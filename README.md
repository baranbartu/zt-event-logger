# ZeroTier Event Logger API

This is a simple web service built with Go and Gin that receives ZeroTier Central event hooks, logs them to a database, and provides an endpoint to search for stored events.

## Features

- Receives event hooks from ZeroTier Central.
- Logs event details (network, device, userID) to a database.
- Exposes an HTTP endpoint to search for stored events by network, device, or user ID.

## TL;DR

Please check the `Makefile` to make the life easier.

- running tests: `make test`
- generate mocks (optional - make sure a global `mockgen` binary is installed - otherwise skip): `make generate-mocks`
- running docker container: `make docker-run`
- jump to [API Endpoints](#api-endpoints) for testing
- jump to [Potential Improvements](#potential-improvements) for seeing my thought process

## Prerequisites

- For local setup: Go 1.21.3 or higher installed on your machine.
- For Docker setup: Docker installed on your machine.
- To generate more mocks (or update the existing ones), `mockgen` should be globally installed.
  But if no change is required, then the existing mocks can be used without having a global `mockgen`
  installed.

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

### 2. Search Events

- **Endpoint:** `/events/search`
- **Method:** `GET`
- **Content-Type:** `application/json`
- **Query Parameters:** `network_id`, `user_id`, `member_id`

- **Response:**

    ```json
    {
        "events": [
            {
                "id": 1,
                "hook_id": "abc123",
                "org_id": "org456",
                "hook_type": "NETWORK_JOIN",
                "network_id": "net789",
                "member_id": "mem012",
                "created_at": "2024-08-18T18:33:57Z"
            }
        ]
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

### Example cURL Requests For Searching Events

```bash
curl -X GET "http://localhost:8080/events/search?network_id=net789"
```

```bash
curl -X GET "http://localhost:8080/events/search?org_id=org456"
```

```bash
curl -X GET "http://localhost:8080/events/search?org_id=org456&member_id=mem012"
```

## Potential Improvements

### Configuration Management

- I'd use a configuration library like `viper` to manage configurations. This allows for more flexible configuration management, including support for environment variables, configuration files, and more.

### Error Handling

- I'd improve error handling by providing more context and using a structured error handling library like `pkg/errors`.
- I'd ensure that all errors are logged with sufficient context to aid debugging.

### Logging

- I'd integrate a structured logging library like `logrus` or `zap` for better logging.
- I'd ensure all critical operations are logged, including API requests, database operations, and error occurrences.

### Validation

- I'd use validation libraries like `go-playground/validator` to validate incoming request payloads.
- I'd ensure that all user inputs are validated to prevent invalid data from entering the system.

### Security

- I'd use HTTPS for all communications to ensure data encryption in transit.
- I'd validate the `X-ZTC-Signature` header to ensure the integrity and authenticity of incoming requests.
- I'd store sensitive information like the pre-shared key securely and avoid hardcoding it in the codebase.

### Metrics and Monitoring

- I'd integrate metrics collection using libraries like `prometheus/client_golang` to monitor the health and performance of the application.
- I'd set up monitoring and alerting to detect and respond to issues promptly.

### Testing

- I'd increase test coverage to include more edge cases and failure scenarios.
- I'd use table-driven tests to simplify and organize test cases.
- I'd ensure that tests are run in a CI/CD pipeline to catch issues early.

### Documentation

- I'd provide comprehensive documentation for the API endpoints, including request and response formats, error codes, and usage examples.
- I'd document the setup and deployment process, including environment variables and configuration options.

### Containerization

- I'd optimize the Dockerfile to reduce the image size and improve build times.

### CI/CD

- I'd set up a CI/CD pipeline using tools like GitHub Actions, GitLab CI, or Jenkins to automate testing, building, and deployment.
- I'd ensure that code quality checks, tests, and security scans are part of the CI/CD pipeline.

### Storage Layer

- I would use a production-ready database (or cluster), which means its redundancy is considered,
  it is fault-tolerant, and it is backed up constantly.
