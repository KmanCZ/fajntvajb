# fajntvajb

## Description

This repository contains simple prototype of a web application for searching and advertising free time activities. The application is written in Go.

## Installation

To run the application, you need to have Docker installed on your machine. You can download Docker from the following link: https://www.docker.com/get-started

After installing Go, you can clone the repository and run the application by executing the following commands:

```bash
git clone git@github.com:KmanCZ/fajntvajb.git
cd fajntvajb
docker compose up
```

The application will be available at [http://localhost:8080](http://localhost:8080).

Apllication offers database seeding with test data. To seed the database, uncomment the SEED ENV in the `docker-compose.yml` file.
