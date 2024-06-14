package com.github.hwebz.serializer;

import com.github.hwebz.grpc.Laptop;
import com.github.hwebz.sample.Generator;
import org.junit.Assert;
import org.junit.Test;

import java.io.IOException;

import static org.junit.Assert.*;

public class SerializerTest {

    @Test
    public void writeAndReadBinaryFile() throws IOException {
        String binaryFile = "laptop.bin";
        Laptop laptop1 = new Generator().NewLaptop();

        Serializer serializer = new Serializer();
        serializer.WriteBinaryFile(laptop1, binaryFile);

        Laptop laptop2 = serializer.ReadBinaryFile(binaryFile);
        Assert.assertEquals(laptop1, laptop2);
    }
}