# Landslide External Server

The `landslide-external-server` project is a backend solution for a landslide monitoring system. It collects, stores, and manages sensor data from communication units via MQTT, providing a way to persist and configure device data for enhanced landslide detection and management.

## Overview

This project serves as the external backend for a landslide monitoring system. It consists of two main components:

1. **MQTT Broker** - Manages the reception, persistence, and distribution of sensor data and configuration updates to communication units (e.g., ESP32 devices).
2. **Go MQTT Backend** - A Go application that publishes configuration updates to devices via MQTT topics and listens for status updates, ensuring configurations are applied successfully.

---

## Architecture

The `landslide-external-server` project runs two main services within Docker containers:

- **Mosquitto MQTT Broker**: A persistent MQTT broker using the Eclipse Mosquitto Docker image, where all sensor data and configuration messages are relayed and stored in a designated volume.
- **Go MQTT Backend**: A service built with Go that handles publishing configurations to communication units and listens for their status responses via MQTT.

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
git clone https://github.com/yourusername/landslide-external-server.git
cd landslide-external-server
```

### Configuration

Configure the project settings by modifying the necessary files:

- **MQTT Broker Configuration**: `mosquitto/config/mosquitto.conf`
- **Go MQTT Backend Configuration**: Adjust MQTT broker address, ports, and other settings in the Go code (`backend/main.go`) if necessary.

Ensure that any credentials or sensitive information are securely managed, possibly through environment variables or Docker secrets.

---

## Running the Project

With Docker Compose configured, start the entire project using:

```bash
docker-compose up -d
```

This command runs both the Mosquitto broker and the Go MQTT backend in detached mode.

To stop the containers, use:

```bash
docker-compose down
```

---

## Usage

### MQTT Broker

The Mosquitto MQTT broker runs on port `1883` and manages all incoming and outgoing messages between the backend and communication units.

- **MQTT Port**: `1883`
- **MQTT Storage Path**: `./mosquitto/data`

### Communication Units

Communication units (e.g., ESP32 devices) connect to the MQTT broker to:

- **Publish Sensor Data**: To topics determined by the system design (e.g., `sensors/{sensor_id}/data`).
- **Subscribe to Configuration Updates**: On their specific topic `comm_unit/{CommUnitID}/config`.
- **Publish Configuration Status**: To `comm_unit/{CommUnitID}/config/status`.

### Backend Configuration Updates

The Go MQTT backend handles:

- **Publishing Configuration Updates**:

  - To update a communication unit's configuration, the backend publishes a message to `comm_unit/{CommUnitID}/config`.

- **Listening for Status Messages**:

  - The backend subscribes to `comm_unit/{CommUnitID}/config/status` to receive feedback from communication units.