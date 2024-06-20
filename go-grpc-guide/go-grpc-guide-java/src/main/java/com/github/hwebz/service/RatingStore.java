package com.github.hwebz.service;

public interface RatingStore {
    Rating Add(String laptopID, double score);
}
