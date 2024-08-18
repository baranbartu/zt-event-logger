# ZeroTier Event Logger API

This is a simple web service built with Go and Gin that receives ZeroTier Central event hooks, logs them to a database, and provides an endpoint to search for stored events.

## Features

- Receives event hooks from ZeroTier Central.
- Logs event details (network, device, userID) to a database.
- Exposes an HTTP endpoint to search for stored events by network, device, or user ID.

## Prerequisites

- For local setup: Go 1.21.3 or higher installed on your machine.
- For Docker setup: Docker installed on your machine.

## Setup

### Option 1: Local Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/baranbartu/zt-event-logger.git
    cd zt-event-logger
    ```

2. **Download dependencies:**

    ```bash
    go mod tidy
    ```

3. **Start the server:**

    ```bash
    source .env && go run main.go
    ```

   The service will run on `http://localhost:8080`.

### Option 2: Docker Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/baranbartu/zt-event-logger.git
    cd zt-event-logger
    ```

2. **Build the Docker image and run containers:**

    ```bash
    docker build -t zt-event-logger .
    ```

3. **Run the Docker container:**

    ```bash
    docker run --env-file .env-docker -p 8080:8080 zt-event-logger
    ```

   The service will run on `http://localhost:8080`.

## Running the Service

Regardless of the setup method you choose, the service will be available at `http://localhost:8080`.


## API Endpoints

### 1. Receive Event Hook

- **Endpoint:** `/events/receive`
- **Method:** `POST`
- **Content-Type:** `application/json`
- **Request Body:**

    ```json
    {
      "hook_id": "abc123",
      "org_id": "org456",
      "hook_type": "NETWORK_JOIN",
      "network_id": "net789",
      "member_id": "mem012"
    }
    ```

- **Response:**

    ```json
    {
      "message": "Event received and logged successfully",
      "hook_id": "abc123",
      "hook_type": "NETWORK_JOIN",
      "org_id": "org456"
    }
    ```

### Example cURL Requests For Receiving Events

#### For Known and Handled Events

```bash
curl -X POST http://localhost:8080/events/receive \
     -H "Content-Type: application/json" \
     -d '{
           "hook_id": "abc123",
           "org_id": "org456",
           "hook_type": "NETWORK_JOIN",
           "network_id": "net789",
           "member_id": "mem012"
         }'
```

```bash
curl -X POST http://localhost:8080/events/receive \
     -H "Content-Type: application/json" \
     -d '{
           "hook_id": "abc123",
           "org_id": "org456",
           "hook_type": "NETWORK_CREATED",
           "network_id": "net789",
           "network_config": {
             "config_key": "config_value"
           },
           "user_id": "user123",
           "user_email": "user@example.com",
           "metadata": {
             "meta_key": "meta_value"
           }
         }'
```

```bash
curl -X POST http://localhost:8080/events/receive \
     -H "Content-Type: application/json" \
     -d '{
           "hook_id": "abc123",
           "org_id": "org456",
           "hook_type": "NETWORK_CONFIG_CHANGED",
           "network_id": "net789",
           "user_id": "user123",
           "user_email": "user@example.com",
           "old_config": {
             "old_config_key": "old_config_value"
           },
           "new_config": {
             "new_config_key": "new_config_value"
           },
           "metadata": {
             "meta_key": "meta_value"
           }
         }'
```

#### For Unknown and Unhandled Event Types

```bash
curl -X POST http://localhost:8080/events/receive \
     -H "Content-Type: application/json" \
     -d '{
           "hook_id": "abc123",
           "org_id": "org456",
           "hook_type": "UNKNOWN_EVENT",
           "network_id": "net789",
           "user_id": "user123",
           "user_email": "user@example.com"
         }'
```
