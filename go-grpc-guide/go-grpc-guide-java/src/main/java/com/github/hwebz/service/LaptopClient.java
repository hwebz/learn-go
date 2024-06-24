package com.github.hwebz.service;

import com.github.hwebz.grpc.*;
import com.github.hwebz.sample.Generator;
import com.google.protobuf.ByteString;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.NettyChannelBuilder;
import io.grpc.stub.StreamObserver;
import io.netty.handler.ssl.SslContext;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import javax.net.ssl.SSLException;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.sql.Time;
import java.util.Iterator;
import java.util.Scanner;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.logging.Level;
import java.util.logging.Logger;

public class LaptopClient {
    private static final Logger logger = Logger.getLogger(LaptopClient.class.getName());
    private static final Log log = LogFactory.getLog(LaptopClient.class);

    private final ManagedChannel channel;
    private final LaptopServiceGrpc.LaptopServiceBlockingStub blockingStub;
    private final LaptopServiceGrpc.LaptopServiceStub asyncStub;

    public LaptopClient(String host, int port, SslContext sslContext) {
        channel = NettyChannelBuilder.forAddress(host, port)
                .sslContext(sslContext)
                .build();
        blockingStub = LaptopServiceGrpc.newBlockingStub(channel);
        asyncStub = LaptopServiceGrpc.newStub(channel);
    }

    public void shutdown() throws InterruptedException {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS);
    }

    public void createLaptop(Laptop laptop) {
        CreateLaptopRequest request = CreateLaptopRequest.newBuilder().setLaptop(laptop).build();
        CreateLaptopResponse response = CreateLaptopResponse.getDefaultInstance();

        try {
            response = blockingStub.withDeadlineAfter(5, TimeUnit.SECONDS).createLaptop(request);
        } catch (StatusRuntimeException e) {
            if (e.getStatus().getCode() == Status.Code.ALREADY_EXISTS) {
                logger.info("Laptop is already exists.");
                return;
            }
            logger.log(Level.SEVERE, "request failed: " + e.getMessage());
            return;
        } catch (Exception e) {
            logger.log(Level.SEVERE, "request failed: " + e.getMessage());
            return;
        }

        logger.info("Laptop created with ID: " + response.getId());
    }

    public void searchLaptop(Filter filter) {
        logger.info("Search started");

        SearchLaptopRequest request = SearchLaptopRequest.newBuilder().setFilter(filter).build();
        try {
            Iterator<SearchLaptopResponse> responseIterator = blockingStub.withDeadlineAfter(5, TimeUnit.SECONDS).searchLaptop(request);

            while (responseIterator.hasNext()) {
                SearchLaptopResponse response = responseIterator.next();
                Laptop laptop = response.getLaptop();
                logger.info("- found" + laptop.getId());
            }
        } catch (Exception e) {
            logger.log(Level.SEVERE, "search failed: " + e.getMessage());
            return;
        }

        logger.info("Search completed");
    }

    public void uploadImage(String laptopID, String imagePath) throws InterruptedException {
        CountDownLatch finishLatch = new CountDownLatch(1);
        StreamObserver<UploadImageRequest> requestStreamObserver = asyncStub.withDeadlineAfter(5, TimeUnit.SECONDS).uploadImage(new StreamObserver<UploadImageResponse>() {
            @Override
            public void onNext(UploadImageResponse uploadImageResponse) {
                logger.info("Receive response: " + uploadImageResponse);
            }

            @Override
            public void onError(Throwable throwable) {
                logger.log(Level.SEVERE, "Upload failed: " + throwable);
                finishLatch.countDown();
            }

            @Override
            public void onCompleted() {
                logger.info("Image uploaded");
                finishLatch.countDown();
            }
        });

        FileInputStream fileInputStream;
        try {
            fileInputStream = new FileInputStream(imagePath);
        } catch (FileNotFoundException e) {
            logger.log(Level.SEVERE, "Cannot read image file: " + e.getMessage());
            return;
        }

        String imageType = imagePath.substring(imagePath.lastIndexOf("."));
        ImageInfo info = ImageInfo.newBuilder().setLaptopId(laptopID).setImageType(imageType).build();
        UploadImageRequest request = UploadImageRequest.newBuilder().setInfo(info).build();

        try {
            requestStreamObserver.onNext(request);
            logger.info("Image info sent: " + info);

            byte[] buffer = new byte[1024];
            while (true) {
                int n = fileInputStream.read(buffer);
                if (n <= 0) {
                    break;
                }

                if (finishLatch.getCount() == 0) {
                    return;
                }

                request = UploadImageRequest.newBuilder()
                        .setChunkData(ByteString.copyFrom(buffer, 0, n))
                        .build();
                requestStreamObserver.onNext(request);
                logger.info("Image chunk sent with size: " + n);
            }
        } catch (Exception e) {
            logger.log(Level.SEVERE, "Unexpected error: " + e.getMessage());
            requestStreamObserver.onError(e);
        }

        requestStreamObserver.onCompleted();

        if (!finishLatch.await(1, TimeUnit.MINUTES)) {
            logger.warning("Request cannot finish within 1 minute");
        }
    }

    public void rateLaptop(String[] laptopIDs, double[] scores) throws InterruptedException {
        CountDownLatch finishLatch = new CountDownLatch(1);
        StreamObserver<RateLaptopRequest> requestStreamObserver = asyncStub.withDeadlineAfter(5, TimeUnit.SECONDS)
                .rateLaptop(new StreamObserver<RateLaptopResponse>() {
                    @Override
                    public void onNext(RateLaptopResponse rateLaptopResponse) {
                        logger.info("laptop rate id: " + rateLaptopResponse.getLaptopId() +
                                ", count = " + rateLaptopResponse.getRatedCount() +
                                ", average = " + rateLaptopResponse.getAverageScore());
                    }

                    @Override
                    public void onError(Throwable throwable) {
                        logger.log(Level.SEVERE, "Rate laptop failed: " + throwable.getMessage());
                        finishLatch.countDown();
                    }

                    @Override
                    public void onCompleted() {
                        logger.info("Rate laptop completed");
                        finishLatch.countDown();
                    }
                });
        int n = laptopIDs.length;
        try {
            for (int i = 0; i < n; i++) {
                RateLaptopRequest request = RateLaptopRequest.newBuilder()
                        .setLaptopId(laptopIDs[i])
                        .setScore(scores[i])
                        .build();
                requestStreamObserver.onNext(request);
                logger.info("Sent rate laptop request with id = " + request.getLaptopId() + ", score = " + request.getScore());
            }
        } catch (Exception e) {
            logger.log(Level.SEVERE, "Unexpected error: " + e.getMessage());
            requestStreamObserver.onError(e);
            return;
        }

        requestStreamObserver.onCompleted();
        if (!finishLatch.await(1, TimeUnit.MINUTES)) {
            logger.warning("Request cannot finish within 1 minute");
        }
    }

    public static void testCreateLaptop(LaptopClient client, Generator generator) {
        Laptop laptop = generator.NewLaptop();
        client.createLaptop(laptop);
    }

    public static void testSearchLaptop(LaptopClient client, Generator generator) {
        for (int i = 0; i < 10; i++) {
            Laptop laptop = generator.NewLaptop();
            client.createLaptop(laptop);
        }

        Memory minRam = Memory.newBuilder().setValue(4).setUnit(Memory.Unit.GIGABYTE).build();
        Filter filter = Filter.newBuilder()
                .setMaxPriceUsd(3000)
                .setMinCpuCores(4)
                .setMinCpuGhz(2.0)
                .setMinRam(minRam)
                .build();

        client.searchLaptop(filter);
    }

    public static void testUploadImage(LaptopClient client, Generator generator) throws InterruptedException {
        Laptop laptop = generator.NewLaptop();
        client.createLaptop(laptop); // commented out FOR TESTING "Laptop not found"
        client.uploadImage(laptop.getId(), "tmp/laptop.png");
    }

    public static void testRateLaptop(LaptopClient client, Generator generator) throws InterruptedException {
        int n = 3;
        String[] laptopIDs = new String[n];

        for (int i = 0; i < n; i++) {
            Laptop laptop = generator.NewLaptop();
            laptopIDs[i] = laptop.getId();
            client.createLaptop(laptop);
        }

        Scanner scanner = new Scanner(System.in);
        while (true) {
            logger.info("Rate laptop (y/n)?");
            String answer = scanner.nextLine();
            if (!answer.toLowerCase().trim().equals("y")) {
                break;
            }

            double[] scores = new double[n];
            for (int i = 0; i < n; i++) {
                scores[i] = generator.randomLaptopScore();
            }

            client.rateLaptop(laptopIDs, scores);
        }
    }

    public static SslContext loadTLSCredentials() throws SSLException {
        File serverCACertFile = new File("cert/ca-cert.pem");
        File clientCertFile = new File("cert/client-cert.pem");
        File clientKeyFile = new File("cert/client-key.pem");

        return GrpcSslContexts.forClient()
                .keyManager(clientCertFile, clientKeyFile)
                .trustManager(serverCACertFile)
                .build();
    }

    public static void main(String[] args) throws InterruptedException, SSLException {
        SslContext sslContext = LaptopClient.loadTLSCredentials();
        LaptopClient client = new LaptopClient("0.0.0.0", 8089, sslContext);

        Generator generator = new Generator();
//        Laptop laptop = generator.NewLaptop();
        // Laptop laptop = generator.NewLaptop().toBuilder().setId("5934909a-ac53-4657-9141-dfb895e290bc").build(); // Already exists
        // Laptop laptop = generator.NewLaptop().toBuilder().setId("invalid-uuid").build(); // Invalid ID/home/vinai/learning/learn-go/go-grpc-guide/tmp/laptop.png

        try {
//            testCreateLaptop(client, generator);
//            testSearchLaptop(client, generator);
//            testUploadImage(client, generator);
            testRateLaptop(client, generator);
        } finally {
            client.shutdown();
        }
    }
}
