package com.github.hwebz.sample;

import com.github.hwebz.grpc.*;
import com.google.protobuf.Timestamp;

import java.time.Instant;
import java.util.Arrays;
import java.util.Objects;
import java.util.Random;

public class Generator {
    private Random rand;

    public Generator() {
        rand = new Random();
    }
    public Keyboard NewKeyboard() {
        return Keyboard.newBuilder()
                .setLayout(randomKeyboardLayout())
                .setBacklit(rand.nextBoolean())
                .build();
    }

    private Keyboard.Layout randomKeyboardLayout() {
        switch (rand.nextInt(2)) {
            case 1:
                return Keyboard.Layout.QWERTY;
            case 2:
                return Keyboard.Layout.QWERTZ;
            default:
                return Keyboard.Layout.AZERTY;
        }
    }

    public CPU NewCPU() {
        String brand = randomCPUBrand();
        String name = randomCPUName(brand);

        int numberCores = randomInt(2, 8);
        int numberThreads = randomInt(numberCores, 12);

        double minGhz = randomDouble(2.0, 3.5);
        double maxGhz = randomDouble(minGhz, 5.0);

        return CPU.newBuilder()
                .setBrand(brand)
                .setName(name)
                .setNumberCores(numberCores)
                .setNumberThreads(numberThreads)
                .setMinGhz(minGhz)
                .setMaxGhz(maxGhz)
                .build();
    }

    public GPU NewGPU() {
        String brand = randomGPUBrand();
        String name = randomGPUName(brand);

        double minGhz = randomDouble(1.0, 1.5);
        double maxGhz = randomDouble(minGhz, 2.0);

        Memory memory = Memory.newBuilder()
                .setValue(randomInt(2, 6))
                .setUnit(Memory.Unit.GIGABYTE)
                .build();

        return GPU.newBuilder()
                .setBrand(brand)
                .setName(name)
                .setMinGhz(minGhz)
                .setMaxGhz(maxGhz)
                .setMemory(memory)
                .build();
    }

    public Memory NewRAM() {
        return Memory.newBuilder()
                .setValue(randomInt(4, 64))
                .setUnit(Memory.Unit.GIGABYTE)
                .build();
    }

    public Storage NewSSD() {
        Memory memory = Memory.newBuilder()
                .setValue(randomInt(128, 1024))
                .setUnit(Memory.Unit.GIGABYTE)
                .build();

        return Storage.newBuilder()
                .setDriver(Storage.Driver.SSD)
                .setMemory(memory)
                .build();
    }

    public Storage NewHDD() {
        Memory memory = Memory.newBuilder()
                .setValue(randomInt(1, 6))
                .setUnit(Memory.Unit.TERABYTE)
                .build();

        return Storage.newBuilder()
                .setDriver(Storage.Driver.HDD)
                .setMemory(memory)
                .build();
    }

    private Screen NewScreen() {
        int height = randomInt(1080, 4320);
        int width = height * 16 / 9;

        Screen.Resolution resolution = Screen.Resolution.newBuilder()
                .setHeight(height)
                .setWidth(width)
                .build();

        return Screen.newBuilder()
                .setResolution(resolution)
                .setMultitouch(rand.nextBoolean())
                .setPanel(randomScreenPanel())
                .setSizeInch(randomInt(13, 17))
                .build();
    }

    public Laptop NewLaptop() {
        String brand = randomLaptopBrand();
        String name = randomLaptopName(brand);

        double weightKg = randomDouble(1.0, 3.0);
        double priceUsd = randomDouble(1500.0, 3500.0);

        int releaseYear = randomInt(2015, 2019);

        return Laptop.newBuilder()
                .setBrand(brand)
                .setName(name)
                .setCpu(NewCPU())
                .setRam(NewRAM())
                .addGpus(NewGPU())
                .addStorages(NewSSD())
                .addStorages(NewHDD())
                .setScreen(NewScreen())
                .setKeyboard(NewKeyboard())
                .setWeightKg(weightKg)
                .setPriceUsd(priceUsd)
                .setReleaseYear(releaseYear)
                .setUpdatedAt(timestampNow())
                .build();
    }

    private Timestamp timestampNow() {
        Instant now = Instant.now();
        return Timestamp.newBuilder()
                .setSeconds(now.getEpochSecond())
                .setNanos(now.getNano())
                .build();
    }

    private String randomCPUBrand() {
        return randomStringFromSet("Intel", "AMD");
    }

    private String randomStringFromSet(String... a) {
        int n = a.length;
        if (n == 0) {
            return "";
        }
        return a[rand.nextInt(n)];
    }

    private String randomCPUName(String brand) {
        if (Objects.equals(brand, "Intel")) {
            return randomStringFromSet("Xeon E-1234v", "Core i9-9900K", "Core i7-7700K", "Core i5-5500K", "Core i3 3300K");
        }

        return randomStringFromSet("Ryzen 7 Pro 2700U", "Ryzen 5 Pro 5100U", "Ryzen 3 Pro 3200U");
    }

    private int randomInt(int min, int max) {
        return min + rand.nextInt(max-min+1);
    }

    private String randomGPUBrand() {
        return randomStringFromSet("NVIDIA", "AMD");
    }

    private Float randomFloat32(Float min, Float max) {
        return min + rand.nextFloat() * (max-min);
    }

    private Double randomDouble(Double min, Double max) {
        return min + rand.nextDouble() * (max - min);
    }

    private String randomGPUName(String brand) {
        if (Objects.equals(brand, "NVIDIA")) {
            return randomStringFromSet("RTX 3090", "RTX 3080", "RTX 3070", "RTX 3060", "RTX 3050");
        }

        return randomStringFromSet("RX 590", "RX 580", "RX 5700XT", "RX Vega-56");
    }

    private Screen.Resolution randomScreenResolution() {
        int height = randomInt(1080, 4320);
        int width = height * 16 / 9;

        return Screen.Resolution.newBuilder()
                .setHeight(height)
                .setWidth(width)
                .build();
    }

    private Screen.Panel randomScreenPanel() {
        if (rand.nextInt(2) == 1) {
            return Screen.Panel.OLED;
        }

        return Screen.Panel.IPS;
    }

    private String randomUUID() {
        return "";
    }

    private String randomLaptopBrand() {
        return randomStringFromSet("Asus", "Dell", "Acer", "HP", "Lenovo", "Apple");
    }

    private String randomLaptopName(String brand) {
        switch (brand) {
            case "Apple":
                return randomStringFromSet("MacBook Air", "MacBook Pro");
            case "Asus":
                return randomStringFromSet("ZenBook", "VivoBook", "ROG Zephyrus", "TUF Gaming");
            case "Dell":
                return randomStringFromSet("XPS", "Inspiron", "Alienware", "Latitude", "Vostro");
            case "Acer":
                return randomStringFromSet("Aspire", "Swift", "Predator", "Nitro", "Spin");
            case "HP":
                return randomStringFromSet("Spectre", "Pavilion", "Envy", "Omen", "EliteBook");
            case "Lenovo":
                return randomStringFromSet("ThinkPad", "IdeaPad", "Legion", "Yoga");
            default:
                return "Unknown Model";
        }
    }

    public static void main(String[] args) {
        Generator g = new Generator();
        Laptop laptop = g.NewLaptop();
        System.out.println(laptop);
    }
}
