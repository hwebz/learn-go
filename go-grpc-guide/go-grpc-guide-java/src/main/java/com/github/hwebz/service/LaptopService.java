package com.github.hwebz.service;

import com.github.hwebz.grpc.*;
import com.google.protobuf.ByteString;
import io.grpc.Context;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.util.UUID;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;

public class LaptopService extends LaptopServiceGrpc.LaptopServiceImplBase {
    private static final Logger logger = Logger.getLogger(LaptopService.class.getName());
    private final LaptopStore laptopStore;
    private final ImageStore imageStore;

    public LaptopService(LaptopStore laptopStore, ImageStore imageStore) {
        this.laptopStore = laptopStore;
        this.imageStore = imageStore;
    }

    @Override
    public void createLaptop(CreateLaptopRequest request, StreamObserver<CreateLaptopResponse> responseObserver) {
        Laptop laptop = request.getLaptop();

        String id = laptop.getId();
        logger.info("got a create-laptop request with ID: " + id);

        UUID uuid;
        if (id.isEmpty()) {
            uuid = UUID.randomUUID();
        } else {
            try {
                uuid = UUID.fromString(id);
            } catch (IllegalArgumentException e) {
                responseObserver.onError(Status.INVALID_ARGUMENT.withDescription(e.getMessage()).asRuntimeException());
                return;
            }
        }

        // FOR TESTING: Simulate heavy processing
//        try {
//            TimeUnit.SECONDS.sleep(6);
//        } catch (InterruptedException e) {
//            e.printStackTrace();
//        }

        // FOR TESTING: Client cancelled the request, we stop the server immediately
        if (Context.current().isCancelled()) {
            logger.info("Request is cancelled");
            responseObserver.onError(Status.CANCELLED.withDescription("Request is cancelled").asRuntimeException());
            return;
        }

        Laptop other = laptop.toBuilder().setId(uuid.toString()).build();
        try {
            laptopStore.Save(other);
        } catch (AlreadyExistsException e) {
            responseObserver.onError(Status.ALREADY_EXISTS.withDescription(e.getMessage()).asRuntimeException());
            return;
        } catch (Exception e) {
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).asRuntimeException());
            return;
        }

        CreateLaptopResponse response = CreateLaptopResponse.newBuilder().setId(other.getId()).build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();

        logger.info("Saved laptop with ID: " + other.getId());
    }

    @Override
    public void searchLaptop(SearchLaptopRequest request, StreamObserver<SearchLaptopResponse> responseObserver) {
        Filter filter = request.getFilter();
        logger.info("Got a search-laptop request with filter: \n" + filter);

        laptopStore.Search(Context.current(), filter, new LaptopStream() {
            @Override
            public void Send(Laptop laptop) {
                logger.info("Found laptop with ID: " + laptop.getId());
                try {
                    TimeUnit.SECONDS.sleep(1);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                SearchLaptopResponse response = SearchLaptopResponse.newBuilder().setLaptop(laptop).build();
                responseObserver.onNext(response);

            }
        });

        responseObserver.onCompleted();
        logger.info("Search laptop completed");
    }

    @Override
    public StreamObserver<UploadImageRequest> uploadImage(StreamObserver<UploadImageResponse> responseObserver) {
        return new StreamObserver<UploadImageRequest>() {
//            private static final int maxImageSize = 1 << 10; // 1Kb FOR TESTING
            private static final int maxImageSize = 1 << 20; // 1Mb
            private String laptopID;
            private String imageType;
            private ByteArrayOutputStream imageData;

            @Override
            public void onNext(UploadImageRequest uploadImageRequest) {
                if (uploadImageRequest.getDataCase() == UploadImageRequest.DataCase.INFO) {
                    ImageInfo info = uploadImageRequest.getInfo();
                    logger.info("Receive image info: " + info);

                    laptopID = info.getLaptopId();
                    imageType = info.getImageType();
                    imageData = new ByteArrayOutputStream();

                    // Check laptop exists
                    Laptop found = laptopStore.Find(laptopID);
                    if (found == null) {
                        responseObserver.onError(
                                Status.NOT_FOUND
                                        .withDescription("Laptop not found")
                                        .asRuntimeException()
                        );
                    }

                    return;
                }

                ByteString chunkData = uploadImageRequest.getChunkData();
                logger.info("Receive image chunk with size: " + chunkData.size());

                if (imageData == null) {
                    logger.info("Image info wasn't sent before");
                    responseObserver.onError(
                            Status.INVALID_ARGUMENT
                                    .withDescription("Image info wasn't sent before")
                                    .asRuntimeException()
                    );
                    return;
                }

                int size = imageData.size() + chunkData.size();
                if (size > maxImageSize) {
                    logger.info("Image is too large: " + size);
                    responseObserver.onError(
                            Status.INVALID_ARGUMENT
                                    .withDescription("Image is too large: " + size)
                                    .asRuntimeException()
                    );
                    return;
                }

                try {
                    chunkData.writeTo(imageData);
                } catch (IOException e) {
                    responseObserver.onError(
                            Status.INTERNAL
                                    .withDescription("Cannot write chunk data: " + e.getMessage())
                                    .asRuntimeException()
                    );
                    return;
                }
            }

            @Override
            public void onError(Throwable throwable) {
                logger.warning(throwable.getMessage());
            }

            @Override
            public void onCompleted() {
                String imageID = "";
                int imageSize = imageData.size();

                try {
                    imageID = imageStore.Save(laptopID, imageType, imageData);
                } catch (IOException e) {
                    responseObserver.onError(
                            Status.INTERNAL
                                    .withDescription("Cannot sav eimage to the store: " + e.getMessage())
                                    .asRuntimeException()
                    );
                }

                UploadImageResponse response = UploadImageResponse.newBuilder()
                        .setId(imageID)
                        .setSize(imageSize)
                        .build();
                responseObserver.onNext(response);
                responseObserver.onCompleted();
            }
        };
    }
}
