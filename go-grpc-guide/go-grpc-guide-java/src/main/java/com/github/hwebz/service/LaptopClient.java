package com.github.hwebz.service;

import com.github.hwebz.grpc.*;
import com.github.hwebz.sample.Generator;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;

import java.sql.Time;
import java.util.Iterator;
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

    public static void main(String[] args) throws InterruptedException {
        LaptopClient client = new LaptopClient("0.0.0.0", 8089);

        Generator generator = new Generator();
//        Laptop laptop = generator.NewLaptop();
        // Laptop laptop = generator.NewLaptop().toBuilder().setId("5934909a-ac53-4657-9141-dfb895e290bc").build(); // Already exists
        // Laptop laptop = generator.NewLaptop().toBuilder().setId("invalid-uuid").build(); // Invalid ID

        try {
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
        } finally {
            client.shutdown();
        }
    }
}
