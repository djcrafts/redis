# Redis-like Server in Go

This project implements a Redis-like key-value store server in Go, with support for commands like `SET`, `GET`, `DEL`, and `HELLO`. It demonstrates how to build a multi-peer server that manages connections, handles key-value operations, and processes commands from clients using the RESP (Redis Serialization Protocol).

## Features
- **Key-Value Store:** Provides a thread-safe in-memory key-value store with `SET` and `GET` commands.
- **Command Support:** Implements basic Redis commands like `SET`, `GET`, `DEL`, and `HELLO`.
- **RESP Protocol:** Communicates with clients using the RESP format, ensuring compatibility with Redis-like clients.
- **Peer Management:** Supports dynamic management of peer connections.
- **Concurrency:** Handles multiple concurrent peer connections efficiently using Go's goroutines.

## Supported Commands
The following commands are implemented in this server:

### `SET <key> <value>`
Sets the value of a key in the key-value store.

- **Example:**  
  `SET mykey myvalue`
- **Response:**  
  `OK`

### `GET <key>`
Retrieves the value associated with the given key.

- **Example:**  
  `GET mykey`
- **Response:**  
  `myvalue` (if the key exists)  
  `(nil)` (if the key does not exist)

### `DEL <key>`
Deletes the specified key from the key-value store.

- **Example:**  
  `DEL mykey`
- **Response:**  
  `(integer) 1` (if the key was successfully deleted)  
  `(integer) 0` (if the key does not exist)

### `HELLO`
Sends a greeting message from the server to the client, providing basic server information.

- **Example:**  
  `HELLO`
- **Response:**  
  `{"server": "GoRedis", "version": "1.0"}`

## Architecture

The server is built with modular components to maintain clean code and separation of concerns. The architecture includes the following key parts:

- **Key-Value Store (`KV`):** Handles storing and retrieving data in memory. Operations are synchronized using `sync.RWMutex` for thread safety.
- **Server:** The entry point for the server that listens for incoming peer connections, handles commands, and routes responses.
- **Peer:** Represents an active connection to the server. Each peer sends commands which are parsed and processed individually.
- **Commands:** Implements core Redis-like commands (`SET`, `GET`, `DEL`, `HELLO`), each of which interacts with the in-memory key-value store.

### Flow Overview
1. **Peer Connections:** The server listens for incoming TCP connections and accepts new peers, adding them to the list of active connections.
2. **Command Processing:** Each peer sends commands in RESP format. The server parses these commands and invokes the corresponding methods to interact with the key-value store.
3. **Key-Value Operations:** Operations like `SET`, `GET`, and `DEL` modify the in-memory store, and results are returned to the client.
4. **Multi-Peer Support:** The server can handle multiple peer connections simultaneously, ensuring independent operation for each peer.

## Future Enhancements

This project serves as a Redis-like server and is continually being improved. Some potential future enhancements include:

- **Support for Additional Data Structures:**  
  Implement Redis-like data types such as Lists, Sets, Hashes, and Sorted Sets, which will extend the functionality of the server.

- **Persistence:**  
  Introduce persistence mechanisms such as snapshots (RDB) or append-only file (AOF) to ensure that data is saved to disk, even if the server is restarted.

- **Replication and Clustering:**  
  Add support for server replication, enabling fault tolerance and horizontal scaling through multiple instances of the server.

- **Authentication:**  
  Implement user authentication and permission management to secure the server from unauthorized access.

- **Advanced Command Handling:**  
  Add support for additional advanced Redis commands such as `INCR`, `EXPIRE`, and `PUBLISH` to make the server more feature-complete.

- **Performance Improvements:**  
  Explore ways to optimize the server for handling high loads and large datasets, ensuring improved throughput and latency.
