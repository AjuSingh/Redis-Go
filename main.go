// Package main is the entry point for our Redis-like server implementation.
// In Go, the main package is special - it defines a standalone executable program, not a library.
package main

// Import necessary standard library packages:
// - fmt: for printing messages and errors
// - net: for network functionality (TCP server)
// - strings: for string manipulation (converting commands to uppercase)
import (
    "fmt"
    "net"
    "strings"
)

// main is the entry point of our program. When you run the program, this function
// gets called first. It sets up our Redis-like server and contains the main server loop.
func main() {
    // Print a message indicating that our server is starting up
    // This will help users know the server is running
    fmt.Println("Listening on port :6379")

    // Create a TCP listener on port 6379 (the default Redis port)
    // net.Listen creates a server that can accept incoming connections
    // "tcp" specifies we want a TCP connection (as opposed to UDP)
    // The second argument ":6379" means listen on all network interfaces on port 6379
    l, err := net.Listen("tcp", ":6379")
    
    // Error handling: if we couldn't create the listener (e.g., port is already in use)
    // print the error and exit the program
    if err != nil {
        fmt.Println(err)
        return
    }

    // Create a new Append-Only File (AOF) for persistence
    // This is how Redis maintains data across server restarts
    // The file will be named "database.aof"
    aof, err := NewAof("database.aof")
    
    // If we couldn't create/open the AOF file, print the error and exit
    if err != nil {
        fmt.Println(err)
        return
    }
    
    // Make sure we close the AOF file when the program exits
    // defer ensures this happens even if we encounter an error
    defer aof.Close()

    // Read existing commands from the AOF file and replay them
    // This restores our database to its state before the last shutdown
    aof.Read(func(value Value) {
        // Extract the command name (like "SET", "GET", etc.) and convert to uppercase
        command := strings.ToUpper(value.array[0].bulk)
        
        // Get the command arguments (everything after the command name)
        args := value.array[1:]

        // Look up the handler function for this command
        handler, ok := Handlers[command]
        
        // If we don't recognize the command, print an error and skip it
        if !ok {
            fmt.Println("Invalid command: ", command)
            return
        }

        // Execute the command with its arguments
        handler(args)
    })

    // Accept a new connection from a client
    // This blocks until a client connects
    conn, err := l.Accept()
    
    // If we couldn't accept the connection, print the error and exit
    if err != nil {
        fmt.Println(err)
        return
    }

    // Ensure we close the connection when we're done with it
    defer conn.Close()

    // Main server loop - this runs forever, processing client commands
    for {
        // Create a new RESP (Redis Serialization Protocol) reader for this connection
        resp := NewResp(conn)
        
        // Read the next command from the client
        value, err := resp.Read()
        
        // If there was an error reading (e.g., client disconnected),
        // print it and exit
        if err != nil {
            fmt.Println(err)
            return
        }

        // Commands should be arrays in RESP format
        // Check that we received an array
        if value.typ != "array" {
            fmt.Println("Invalid request, expected array")
            continue  // Skip this command and wait for the next one
        }

        // Check that the array isn't empty
        // (Every command needs at least a command name)
        if len(value.array) == 0 {
            fmt.Println("Invalid request, expected array length > 0")
            continue
        }

        // Extract the command name and convert to uppercase
        // Commands in Redis are case-insensitive
        command := strings.ToUpper(value.array[0].bulk)
        
        // Get the command arguments
        args := value.array[1:]

        // Create a writer to send responses back to the client
        writer := NewWriter(conn)

        // Look up the handler function for this command
        handler, ok := Handlers[command]
        
        // If we don't recognize the command, send an empty response
        if !ok {
            fmt.Println("Invalid command: ", command)
            writer.Write(Value{typ: "string", str: ""})
            continue
        }

        // If this is a write command (SET or HSET),
        // write it to the AOF file for persistence
        if command == "SET" || command == "HSET" {
            aof.Write(value)
        }

        // Execute the command and send the result back to the client
        result := handler(args)
        writer.Write(result)
    }
}