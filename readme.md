# Mock SES API

This project is a mock implementation of the AWS Simple Email Service (SES) API using the Gin web framework. It provides endpoints for sending emails, checking the health of the service, and retrieving email statistics.

## Prerequisites

- Go 1.16 or later
- Gin web framework

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/mock-ses-api.git
    cd mock-ses-api
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

## Usage

1. Run the application:

    ```sh
    go run main.go
    ```

2. The API will be available at `http://localhost:8080`.

## Endpoints

### Health Check

- **URL:** `/api/v1/emails/health`
- **Method:** `GET`
- **Description:** Checks the health of the service.

### Send Email

- **URL:** `/api/v1/emails/outbound-emails`
- **Method:** `POST`
- **Description:** Sends an email.
- **Request Body:**
    ```json
    {
        "to": "recipient@example.com",
        "subject": "Email Subject",
        "body": "Email Body"
    }
    ```

### Get Stats

- **URL:** `/api/v1/emails/stats`
- **Method:** `GET`
- **Description:** Retrieves email statistics.

## Project Structure

- `internal/config`: Configuration management.
- `internal/handlers`: HTTP handlers for the API endpoints.
- `pkg/stats`: Email statistics management.
- `routes`: API route setup.

## Testing

Run the tests using the following command:

```sh
go test ./...