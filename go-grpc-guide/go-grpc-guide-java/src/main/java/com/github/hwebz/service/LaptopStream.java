package com.github.hwebz.service;

import com.github.hwebz.grpc.Laptop;

public interface LaptopStream {
    void Send(Laptop laptop);
}
