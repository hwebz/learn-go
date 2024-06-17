package com.github.hwebz.service;

import com.github.hwebz.grpc.Filter;
import com.github.hwebz.grpc.Laptop;
import io.grpc.Context;

public interface LaptopStore {
    void Save(Laptop laptop) throws Exception;
    Laptop Find(String id);
    void Search(Context context, Filter filter, LaptopStream stream);
}

