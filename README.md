# Go-Redis Implementation

<div align="center">
  <img src="https://raw.githubusercontent.com/ashleymcnamara/gophers/master/GO_BUILD.png" width="150">
  <img src="https://raw.githubusercontent.com/redis/redis/unstable/redis.png" width="150">
</div>

## Overview

A high-performance, concurrent Redis-compatible server implementation in Go, featuring thread-safe operations and persistent storage. This project demonstrates the implementation of fundamental Redis commands while leveraging Go's powerful concurrency primitives and efficient I/O handling.

### Key Features

- 🚀 **High-Performance Concurrent Operations**: Utilizes Go's goroutines and mutexes for thread-safe parallel request handling
- 💾 **Persistent Storage**: Implements Append-Only File (AOF) persistence with configurable sync intervals
- 🔄 **RESP Protocol**: Custom implementation of the Redis Serialization Protocol
- 🏗️ **Multiple Data Structures**: Support for Strings and Hash data types
- 🔒 **Thread-Safe**: Robust concurrency handling with Read-Write mutexes
- 🔄 **Background Processing**: Asynchronous AOF syncing for improved performance

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
- Go 1.19 or higher
- Basic understanding of Redis commands

### Installation

1. Clone the repository:
```bash
git clone https://github.com/skylarshi123/RedisFromScratch.git
cd go-redis
```

2. Build the server:
```bash
go build .
```

3. Run the server:
```bash
./go-redis
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

### Performance Considerations

- Utilizes Go's efficient concurrency primitives
- Implements Read-Write locks for optimized concurrent access
- Background AOF syncing to minimize I/O impact
- Buffered I/O operations for improved performance

## Implementation Details

### Persistence (AOF.go)
The Append-Only File (AOF) implementation provides:
- Automatic background syncing every second
- Mutex-protected file operations
- Command replay on server startup

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
Made with ❤️ using Go