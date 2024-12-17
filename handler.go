// Package main implements a Redis-like server with basic command handling functionality
// This file specifically handles the implementation of Redis commands like SET, GET, HSET, etc.
package main

// Import the sync package which provides basic synchronization primitives
// We need this for mutual exclusion (mutex) to handle concurrent access to our data stores
import (
    "sync"
	"strconv"
)

// Handlers maps Redis command names to their corresponding handler functions
// Each handler function takes a slice of Values (command arguments) and returns a Value (the response)
// This is our command registry - it tells the server which function to call for each Redis command
var Handlers = map[string]func([]Value) Value{
    "PING":    ping,     // Simple server health check command
    "SET":     set,      // Set a key-value pair
    "GET":     get,      // Retrieve a value by key
    "HSET":    hset,     // Set a field in a hash structure
    "HGET":    hget,     // Get a field from a hash structure
    "HGETALL": hgetall,  // Get all fields and values from a hash structure
	"DEL":     del,  // Add our new DEL command
}

// ping implements the PING command from Redis protocol
// If called without arguments, returns "PONG"
// If called with an argument, echoes back that argument
// This is commonly used to test if the server is alive and responding
func ping(args []Value) Value {
    // If no arguments provided, return the standard "PONG" response
    if len(args) == 0 {
        return Value{typ: "string", str: "PONG"}
    }

    // If an argument was provided, echo it back to the client
    // args[0].bulk contains the first argument's value
    return Value{typ: "string", str: args[0].bulk}
}

// SETs is our key-value store for string values
// This is a simple map that stores key-value pairs for the SET and GET commands
// It's the equivalent of Redis's string data type storage
var SETs = map[string]string{}

// SETsMu is a read-write mutex that protects access to the SETs map
// RWMutex allows multiple simultaneous readers but only one writer
// This ensures thread-safety when multiple clients are accessing the data
var SETsMu = sync.RWMutex{}

// set implements the Redis SET command
// It stores a key-value pair in the SETs map
// The command format is: SET key value
func set(args []Value) Value {
    // SET command requires exactly 2 arguments: key and value
    if len(args) != 2 {
        return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
    }

    // Extract key and value from the arguments
    key := args[0].bulk    // First argument is the key
    value := args[1].bulk  // Second argument is the value

    // Lock the mutex before modifying the map
    // This ensures no other goroutine can access the map while we're writing
    SETsMu.Lock()
    SETs[key] = value  // Store the key-value pair
    SETsMu.Unlock()    // Release the lock immediately after writing

    // Return OK to indicate successful operation
    return Value{typ: "string", str: "OK"}
}

// get implements the Redis GET command
// It retrieves a value from the SETs map by its key
// The command format is: GET key
func get(args []Value) Value {
    // GET command requires exactly 1 argument: the key
    if len(args) != 1 {
        return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
    }

    // Extract the key from the arguments
    key := args[0].bulk

    // Get a read lock - multiple goroutines can read simultaneously
    SETsMu.RLock()
    value, ok := SETs[key]  // Attempt to get the value and whether it exists
    SETsMu.RUnlock()       // Release the read lock

    // If the key doesn't exist, return null
    // This matches Redis behavior for non-existent keys
    if !ok {
        return Value{typ: "null"}
    }

    // Return the value as a bulk string
    return Value{typ: "bulk", bulk: value}
}

// HSETs is our hash table store
// It's a nested map: the outer map keys are hash names, and each value is another map
// The inner maps represent hash fields and their values
// This implements Redis's hash data structure
var HSETs = map[string]map[string]string{}

// HSETsMu protects access to the HSETs map
// Like SETsMu, this ensures thread-safe access to our hash structures
var HSETsMu = sync.RWMutex{}

// hset implements the Redis HSET command
// It sets a field value within a hash structure
// The command format is: HSET hash field value
func hset(args []Value) Value {
    // HSET requires exactly 3 arguments: hash name, field, and value
    if len(args) != 3 {
        return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
    }

    // Extract arguments
    hash := args[0].bulk   // Name of the hash
    key := args[1].bulk    // Field name within the hash
    value := args[2].bulk  // Value to store

    // Lock for writing since we're modifying the structure
    HSETsMu.Lock()
    // If this hash doesn't exist yet, create a new empty hash map
    if _, ok := HSETs[hash]; !ok {
        HSETs[hash] = map[string]string{}
    }
    // Set the field value in the hash
    HSETs[hash][key] = value
    HSETsMu.Unlock()

    // Return OK to indicate successful operation
    return Value{typ: "string", str: "OK"}
}

// hget implements the Redis HGET command
// It retrieves the value of a field from a hash structure
// The command format is: HGET hash field
func hget(args []Value) Value {
    // HGET requires exactly 2 arguments: hash name and field
    if len(args) != 2 {
        return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
    }

    // Extract arguments
    hash := args[0].bulk  // Name of the hash
    key := args[1].bulk   // Field name to retrieve

    // Get a read lock
    HSETsMu.RLock()
    value, ok := HSETs[hash][key]  // Attempt to get the field value
    HSETsMu.RUnlock()

    // If either the hash doesn't exist or the field doesn't exist, return null
    if !ok {
        return Value{typ: "null"}
    }

    // Return the field value
    return Value{typ: "bulk", bulk: value}
}

// hgetall implements the Redis HGETALL command
// It returns all fields and values of a hash structure
// The command format is: HGETALL hash
func hgetall(args []Value) Value {
    // HGETALL requires exactly 1 argument: the hash name
    if len(args) != 1 {
        return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
    }

    // Extract the hash name
    hash := args[0].bulk

    // Get a read lock
    HSETsMu.RLock()
    value, ok := HSETs[hash]  // Get the entire hash structure
    HSETsMu.RUnlock()

    // If the hash doesn't exist, return null
    if !ok {
        return Value{typ: "null"}
    }

    // Create an array to hold all field-value pairs
    // In Redis protocol, HGETALL returns an array where elements alternate between
    // field names and their values
    values := []Value{}
    for k, v := range value {
        // Add field name to array
        values = append(values, Value{typ: "bulk", bulk: k})
        // Add field value to array
        values = append(values, Value{typ: "bulk", bulk: v})
    }

    // Return the array of field-value pairs
    return Value{typ: "array", array: values}
}

func del(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'del' command"}
	}
	deletedCount := 0
	SETsMu.Lock()
	HSETsMu.Lock()
	defer SETsMu.Unlock()
	defer HSETsMu.Unlock()
	for _, arg := range args {
		key := arg.bulk
		
		// Check SETs
		if _, exists := SETs[key]; exists {
			delete(SETs, key)
			deletedCount++
			continue
		}
	
		// Check HSETs
		if _, exists := HSETs[key]; exists {
			delete(HSETs, key)
			deletedCount++
		}
	}
	return Value{
		typ: "string",
		str: strconv.Itoa(deletedCount),
	}
}