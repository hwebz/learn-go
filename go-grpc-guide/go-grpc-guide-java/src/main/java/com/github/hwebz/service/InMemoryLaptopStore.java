package com.github.hwebz.service;

import com.github.hwebz.grpc.Laptop;

import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

public class InMemoryLaptopStore implements LaptopStore{
    private ConcurrentMap<String, Laptop> data;

    public InMemoryLaptopStore() {
        data = new ConcurrentHashMap<>(0);
    }
    @Override
    public void Save(Laptop laptop) throws Exception {
        if (data.containsKey(laptop.getId())) {
            throw new AlreadyExistsException("Laptop ID already exists");
        }

        // deep copy
        Laptop other = laptop.toBuilder().build();
        data.put(laptop.getId(), other);
    }

    @Override
    public Laptop Find(String id) {
        if (!data.containsKey(id)) {
            return null;
        }

        // deep copy
        return data.get(id).toBuilder().build();
    }
}
