# Landslide External Server

The `landslide-external-server` project is a backend solution to manage data for a landslide monitoring system. It collects, stores, and serves sensor data from communication units via MQTT and HTTP endpoints, providing a way to persist and configure device data for enhanced landslide detection and management.

## Overview

This project serves as the external backend for a landslide monitoring system. It consists of two main services:

1. **MQTT Broker** - Manages the reception and persistence of sensor data from communication units (e.g., ESP32 devices).
2. **Go-based HTTP Server** - Provides endpoints to manage and configure the communication units, such as retrieving and updating device configurations stored in `config.json`.

---

## Architecture

The `landslide-external-server` project runs two main services within Docker containers:

- **Mosquitto MQTT Broker**: A persistent MQTT broker using the Eclipse Mosquitto Docker image, where all sensor data is relayed and stored in a designated volume.
- **Go HTTP Server**: A REST API built with the Go programming language, which interacts with MQTT and provides configuration endpoints.

Both services are orchestrated with Docker Compose for easier management and deployment.

---

## Setup

### Requirements

- **Docker** and **Docker Compose** installed on your machine. 
  - [Docker Installation Guide](https://docs.docker.com/get-docker/)
  - [Docker Compose Installation Guide](https://docs.docker.com/compose/install/)

### Docker Installation

After installing Docker and Docker Compose, clone this repository:

```bash
git clone https://github.com/rubenszinho/landslide-external-server.git
cd landslide-external-server
```

### Configuration

Configure the project settings by modifying the necessary files in `go_backend/configs/` and `mosquitto/config/`.

- **MQTT Broker Configuration**: `mosquitto/config/mosquitto.conf`
- **Go Server Configuration**: Configure paths for any secrets or external files in `go_backend/configs/config.json`

---

## Running the Project

With Docker Compose configured, start the entire project using:

```bash
docker-compose up -d
```

This command runs both the Mosquitto broker and the Go HTTP server in detached mode.

To stop the containers, use:

```bash
docker-compose down
```

---

## Usage

### MQTT Broker

The Mosquitto MQTT broker runs on port `1883` and stores all incoming messages from sensor units. MQTT clients (e.g., the communication unit) can publish to this broker, which retains data in `./mosquitto/data`.

- **MQTT Port**: `1883`
- **MQTT Storage Path**: `./mosquitto/data`

### HTTP Endpoints

The Go HTTP server provides a RESTful API to manage device configurations.

- **Base URL**: `http://localhost:5000`

#### Endpoints

1. **GET /config**: Fetches the current device configuration.
2. **POST /config**: Updates the device configuration.
   - **Payload**: JSON structure matching `config.json` format.

Example usage of these endpoints can be tested with tools like `curl` or Postman.

---

## Testing the System

To ensure the MQTT broker and HTTP server are functioning correctly:

1. **MQTT Testing**: Publish and subscribe to topics using tools like `mosquitto_pub` and `mosquitto_sub`.
   ```bash
   # Publish test message
   mosquitto_pub -h localhost -p 1883 -t "test/topic" -m "Hello, Mosquitto!"

   # Subscribe to test topic
   mosquitto_sub -h localhost -p 1883 -t "test/topic"
   ```
2. **HTTP Testing**: Use `curl` or Postman to hit the `/config` endpoints.
   ```bash
   # Fetch configuration
   curl -X GET http://localhost:5000/config

   # Update configuration
   curl -X POST http://localhost:5000/config -H "Content-Type: application/json" -d @path_to_config_file
   ```

3. **Persistent Data Check**: Verify that all MQTT messages are stored in `./mosquitto/data`.