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
