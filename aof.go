// Package main implements AOF (Append Only File) persistence
// AOF is Redis's primary persistence mechanism, logging all write commands
package main

// Import required packages
import (
    "bufio"    // For buffered I/O operations
    "io"       // For basic I/O interfaces
    "os"       // For file operations
    "sync"     // For mutex synchronization
    "time"     // For sleep operations
)

// Aof represents an Append Only File
// It handles persistence by logging all write operations to disk
type Aof struct {
    file *os.File         // The actual file on disk
    rd   *bufio.Reader    // Buffered reader for reading the file
    mu   sync.Mutex       // Mutex to protect concurrent access
}

// NewAof creates a new AOF handler
// path: the filesystem path where the AOF file will be stored
func NewAof(path string) (*Aof, error) {
    // Open the file with create, read, and write permissions
    // O_CREATE: create file if it doesn't exist
    // O_RDWR: open for reading and writing
    // 0666: read/write permissions for user, group, and others
    f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
    if err != nil {
        return nil, err
    }

    // Create new AOF instance
    aof := &Aof{
        file: f,
        rd:   bufio.NewReader(f),
    }

    // Start background goroutine for periodic disk sync
    // This ensures durability while maintaining performance
    go func() {
        for {
            aof.mu.Lock()           // Acquire lock
            aof.file.Sync()         // Force write to disk
            aof.mu.Unlock()         // Release lock
            time.Sleep(time.Second) // Wait 1 second before next sync
        }
    }()

    return aof, nil
}

// Close safely closes the AOF file
// This should be called when shutting down the server
func (aof *Aof) Close() error {
    aof.mu.Lock()
    defer aof.mu.Unlock()  // Ensure lock is released even if Close fails

    return aof.file.Close()
}

// Write appends a new command to the AOF file
// This is called for every write operation (SET, HSET, etc.)
func (aof *Aof) Write(value Value) error {
    aof.mu.Lock()
    defer aof.mu.Unlock()  // Ensure lock is released after write

    // Marshal the command to RESP format and write to file
    _, err := aof.file.Write(value.Marshal())
    if err != nil {
        return err
    }

    return nil
}

// Read processes all commands in the AOF file
// This is called during server startup to rebuild the database state
// fn is a callback function that processes each command
func (aof *Aof) Read(fn func(value Value)) error {
    aof.mu.Lock()
    defer aof.mu.Unlock()

    // Seek to start of file
    aof.file.Seek(0, io.SeekStart)

    // Create a RESP reader for the file
    reader := NewResp(aof.file)

    // Read and process each command
    for {
        // Read next command
        value, err := reader.Read()
        if err != nil {
            // If we've reached end of file, break
            if err == io.EOF {
                break
            }
            return err
        }

        // Process the command using callback function
        fn(value)
    }

    return nil
}