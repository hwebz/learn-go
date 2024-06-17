package com.github.hwebz.service;

import com.github.hwebz.grpc.CreateLaptopRequest;
import com.github.hwebz.grpc.CreateLaptopResponse;
import com.github.hwebz.grpc.Laptop;
import com.github.hwebz.grpc.LaptopServiceGrpc;
import com.github.hwebz.sample.Generator;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;

import java.util.concurrent.TimeUnit;
import java.util.logging.Level;
import java.util.logging.Logger;

public class LaptopClient {
    private static final Logger logger = Logger.getLogger(LaptopClient.class.getName());

    private final ManagedChannel channel;
    private final LaptopServiceGrpc.LaptopServiceBlockingStub blockingStub;

    public LaptopClient(String host, int port) {
        channel = ManagedChannelBuilder.forAddress(host, port)
                .usePlaintext()
                .build();
        blockingStub = LaptopServiceGrpc.newBlockingStub(channel);
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

    public static void main(String[] args) throws InterruptedException {
        LaptopClient client = new LaptopClient("0.0.0.0", 8089);

        Generator generator = new Generator();
        Laptop laptop = generator.NewLaptop();
        // Laptop laptop = generator.NewLaptop().toBuilder().setId("5934909a-ac53-4657-9141-dfb895e290bc").build(); // Already exists
        // Laptop laptop = generator.NewLaptop().toBuilder().setId("invalid-uuid").build(); // Invalid ID

        try {
            client.createLaptop(laptop);
        } finally {
            client.shutdown();
        }
    }
}
