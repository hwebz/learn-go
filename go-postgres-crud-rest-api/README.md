# CRUD Operations on PostgreSQL with a Golang REST API
Link: [codevoweb.com](https://codevoweb.com/crud-operations-on-postgresql-with-a-golang-rest-api/)

## Prerequisites
- Go
- Docker
- PostgreSQL

## Init project
```bash
go mod init github.com/hwebz/go-postgres-crud-rest-api
```

## Add required packages
```bash
# Web framework inspired by Express.js
go get github.com/gofiber/fiber/v2

# A package for generating UUIDs in Go
go get github.com/google/uuid

# Set of tools for validating struct fields
go get github.com/go-playground/validator/v10

# An ORM library for Go
go get -u gorm.io/gorm

# A PostgreSQL driver for GORM
go get gorm.io/driver/postgres

# A configuration management package that loads configuration values from different sources, env files or config files
go get github.com/spf13/viper

# For development purpose ONLY
go install github.com/air-verse/air@latest
```

## Run the application
```bash
# Run PostgreSQL database using Docker
make up

# Run Golang application
make run
```