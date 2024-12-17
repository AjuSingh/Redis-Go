# Go-Redis Implementation

<div align="center">
  <img src="https://raw.githubusercontent.com/ashleymcnamara/gophers/master/GO_BUILD.png" width="150">
  <img src="https://redis.com/wp-content/uploads/2021/08/redis-logo.png" width="200">
  <img src="https://www.docker.com/wp-content/uploads/2022/03/vertical-logo-monochromatic.png" width="100">
</div>

## Overview

A high-performance, containerized Redis-compatible server implementation in Go, featuring thread-safe operations and persistent storage. This project demonstrates the implementation of fundamental Redis commands while leveraging Go's powerful concurrency primitives and efficient I/O handling, all packaged in a lightweight Docker container.

### Key Features

- 🚀 **High-Performance Concurrent Operations**: Utilizes Go's goroutines and mutexes for thread-safe parallel request handling
- 💾 **Persistent Storage**: Implements Append-Only File (AOF) persistence with configurable sync intervals
- 🔄 **RESP Protocol**: Custom implementation of the Redis Serialization Protocol
- 🏗️ **Multiple Data Structures**: Support for Strings and Hash data types
- 🔒 **Thread-Safe**: Robust concurrency handling with Read-Write mutexes
- 🔄 **Background Processing**: Asynchronous AOF syncing for improved performance
- 🐳 **Docker Support**: Optimized multi-stage builds reducing image size by 90%

## Supported Commands

### String Operations
- `SET`: Set key to hold a string value
- `GET`: Get the value of a key
- `DEL`: Delete a key

### Hash Operations
- `HSET`: Set field in a hash stored at key to value
- `HGET`: Get the value of a field in a hash
- `HGETALL`: Get all fields and values in a hash

### Connection Management
- `PING`: Test connection to server

## Quick Start

### Prerequisites
- Go 1.21 or higher (for local development)
- Docker (for containerized deployment)
- Basic understanding of Redis commands

### Docker Installation (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/skylarshi123/RedisFromScratch.git
cd RedisFromScratch
```

2. Build and run with Docker:
```bash
# Build the image
docker build -t redis-from-scratch .

# Run the container
docker run -p 6379:6379 -d redis-from-scratch
```

3. For persistence, use Docker volumes:
```bash
docker run -p 6379:6379 -v redis-data:/app/data -d redis-from-scratch
```

### Local Installation (Alternative)

1. Build the server:
```bash
go build .
```

2. Run the server:
```bash
./redis-from-scratch
```

The server will start listening on port 6379 (default Redis port).

### Usage Example

Using `redis-cli`:
```bash
$ redis-cli
127.0.0.1:6379> SET mykey "Hello"
OK
127.0.0.1:6379> GET mykey
"Hello"
127.0.0.1:6379> HSET user:1 name "John"
OK
127.0.0.1:6379> HGET user:1 name
"John"
```

## Technical Details

### Architecture

The implementation follows a modular design with several key components:
- **RESP Protocol Handler**: Manages Redis protocol serialization/deserialization
- **Command Handlers**: Individual command implementations with thread safety
- **AOF Persistence**: Ensures data durability through append-only logging
- **Concurrent Access Management**: RWMutex-based access control
- **Docker Container**: Optimized multi-stage build with Alpine Linux base

### Performance Considerations

- Utilizes Go's efficient concurrency primitives
- Implements Read-Write locks for optimized concurrent access
- Background AOF syncing to minimize I/O impact
- Buffered I/O operations for improved performance
- Minimal Docker image footprint for faster deployments

### Docker Optimization

- Multi-stage builds for minimal image size
- Alpine-based images for reduced footprint
- Volume mounting for persistent storage
- Container health checks
- Automated container orchestration ready

## Implementation Details

### Persistence (AOF.go)
The Append-Only File (AOF) implementation provides:
- Automatic background syncing every second
- Mutex-protected file operations
- Command replay on server startup
- Docker volume support for data persistence

### Protocol (resp.go)
RESP protocol implementation supports:
- Bulk Strings
- Arrays
- Simple Strings
- Errors
- Null values

### Command Handling (handler.go)
Thread-safe command implementations with:
- Concurrent access protection
- Error handling
- Type-safe operations

## Contributing

Contributions are welcome! Feel free to:
- Report bugs
- Suggest new features
- Submit pull requests

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by Redis
- Built with Go's excellent standard library
- Special thanks to the Go and Redis communities

---
Made with ❤️ using Go and Docker
