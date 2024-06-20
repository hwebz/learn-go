package com.github.hwebz.service;

public class ImageMetaData {
    private String laptopID;
    private String type;
    private String path;

    public ImageMetaData(String laptopID, String type, String path) {
        this.laptopID = laptopID;
        this.type = type;
        this.path = path;
    }

    public String getLaptopID() {
        return laptopID;
    }

    public void setLaptopID(String laptopID) {
        this.laptopID = laptopID;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getPath() {
        return path;
    }

    public void setPath(String path) {
        this.path = path;
    }
}
