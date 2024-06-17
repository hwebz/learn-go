package com.github.hwebz.service;

import com.github.hwebz.grpc.Laptop;

public interface LaptopStore {
    void Save(Laptop laptop) throws Exception;
    Laptop Find(String id);
}
