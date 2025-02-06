# WebSocket Server

## Overview

This is a WebSocket server implemented in Go for Optionx. It supports multiple client connections, ping/pong messages, and message routing and broadcasting.

### Assignment Overview

In this assignment, I developed a WebSocket server in Go, designed to handle multiple client connections efficiently. The server is capable of maintaining active connections, sending regular ping messages, and gracefully disconnecting inactive clients. Each client is assigned a unique UUID for internal identification and a randomly generated username for user-friendly interactions.

Upon a new client's connection, the server broadcasts a welcome message, including the usernames of all connected clients, and notifies others of the new client's arrival. The server supports both broadcasting messages to all clients and sending direct messages to specific clients, with messages formatted to include the sender's username for clarity.

The server is structured to handle graceful shutdowns, ensuring all resources are properly released. The codebase is organized into distinct components, following industry standards, and includes detailed comments for maintainability and clarity.

## Running the Server

1. Ensure you have Go installed on your system.
2. Clone the repository.
3. Navigate to the project root directory.
4. Run the server using the following command:
   ```bash
   bash ./script/run.sh
   ```

## Dependencies

- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [Google UUID](https://github.com/google/uuid)
- [Gofakeit](https://github.com/brianvoe/gofakeit/v6)

## Features

- Accepts and maintains multiple connections.
- Sends ping messages and handles pong responses.
- Disconnects inactive clients.
- Routes messages to specific clients based on ID.
- Broadcasts new client connections to all clients.
- Assigns random usernames to clients for display.

## Testing with Postman

### Prerequisites

- [Postman](https://www.postman.com/downloads/) installed on your machine.

### Steps

1. Start the WebSocket server by running the provided script or command.
2. Open Postman and create a new WebSocket request.
3. Enter the WebSocket server URL in the request URL field (e.g., `ws://localhost:8081/ws`).
4. Send a WebSocket message to the server and observe the response.
5. Test broadcasting messages to all clients and sending direct messages to specific clients.
6. Verify the server's handling of ping/pong messages and disconnections.

### Example Request

- **URL**: `ws://localhost:8081/ws`
- **Request Body**: `{ "id": "broadcast", "message": "Hello, everyone!" }`
- **Request Body for direct message**: `{ "id": "client123", "message": "Hello, client123! This is a direct message for you." }`

### Expected Responses

- Upon connecting to the server, you should receive a welcome message.
- Broadcasting a message should result in all connected clients receiving the message.
- Sending a direct message to a specific client should only be received by that client.
- The server's handling of ping/pong messages and disconnections should be observed.

By following these steps, you can use Postman to test the functionality of the WebSocket server and verify its behavior in different scenarios.
