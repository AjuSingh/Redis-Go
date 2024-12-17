// Package main implements the RESP (Redis Serialization Protocol) parser and writer
// RESP is the protocol Redis uses for client-server communication
package main

// Import necessary packages for I/O operations and data conversion
import (
    "bufio"     // Provides buffered I/O for efficient reading
    "fmt"       // For formatting and printing error messages
    "io"        // Basic interfaces for I/O operations
    "strconv"   // For converting between strings and numbers
)

// RESP protocol type markers
// Each RESP data type starts with a specific character that identifies its type
const (
    STRING  = '+'  // Simple String: "+OK\r\n"
    ERROR   = '-'  // Error: "-Error message\r\n"
    INTEGER = ':'  // Integer: ":1000\r\n"
    BULK    = '$'  // Bulk String: "$5\r\nHello\r\n"
    ARRAY   = '*'  // Array: "*2\r\n$5\r\nHello\r\n$5\r\nWorld\r\n"
)

// Value represents a RESP data type and its contents
// This is our internal representation of RESP data
type Value struct {
    typ   string    // Type of value ("string", "error", "integer", "bulk", "array")
    str   string    // Holds simple strings and error messages
    num   int       // Holds integer values
    bulk  string    // Holds bulk strings
    array []Value   // Holds arrays (can contain any other RESP values)
}

// Resp represents a RESP protocol parser
// It wraps a buffered reader for efficient reading of RESP data
type Resp struct {
    reader *bufio.Reader
}

// NewResp creates a new RESP parser from any io.Reader
// This could be a network connection, file, or any other input stream
func NewResp(rd io.Reader) *Resp {
    return &Resp{reader: bufio.NewReader(rd)}
}

// readLine reads a RESP line ending with \r\n
// Returns the line without \r\n, the number of bytes read, and any error
func (r *Resp) readLine() (line []byte, n int, err error) {
    // Keep reading bytes until we find \r\n or encounter an error
    for {
        // Read one byte at a time
        b, err := r.reader.ReadByte()
        if err != nil {
            return nil, 0, err
        }
        n += 1  // Track number of bytes read
        line = append(line, b)  // Add byte to our line buffer
        
        // Check if we've found \r\n (CRLF)
        if len(line) >= 2 && line[len(line)-2] == '\r' {
            break
        }
    }
    // Return line without the trailing \r\n
    return line[:len(line)-2], n, nil
}

// readInteger reads a RESP integer
// Used for array lengths and bulk string lengths
func (r *Resp) readInteger() (x int, n int, err error) {
    // Read the line containing the number
    line, n, err := r.readLine()
    if err != nil {
        return 0, 0, err
    }
    
    // Convert string to integer
    i64, err := strconv.ParseInt(string(line), 10, 64)
    if err != nil {
        return 0, n, err
    }
    
    return int(i64), n, nil
}

// Read reads a complete RESP value
// This is the main entry point for parsing RESP data
func (r *Resp) Read() (Value, error) {
    // Read the type marker byte
    _type, err := r.reader.ReadByte()
    if err != nil {
        return Value{}, err
    }

    // Parse different RESP types based on the marker
    switch _type {
    case ARRAY:
        return r.readArray()
    case BULK:
        return r.readBulk()
    default:
        fmt.Printf("Unknown type: %v", string(_type))
        return Value{}, nil
    }
}

// readArray reads a RESP array
// Format: *<length>\r\n<element-1>...<element-n>
func (r *Resp) readArray() (Value, error) {
    v := Value{}
    v.typ = "array"

    // Read array length
    len, _, err := r.readInteger()
    if err != nil {
        return v, err
    }

    // Initialize array to store elements
    v.array = make([]Value, 0)
    
    // Read each array element
    for i := 0; i < len; i++ {
        // Recursively read each value
        val, err := r.Read()
        if err != nil {
            return v, err
        }
        // Add value to array
        v.array = append(v.array, val)
    }

    return v, nil
}

// readBulk reads a RESP bulk string
// Format: $<length>\r\n<data>\r\n
func (r *Resp) readBulk() (Value, error) {
    v := Value{}
    v.typ = "bulk"

    // Read string length
    len, _, err := r.readInteger()
    if err != nil {
        return v, err
    }

    // Allocate buffer for string data
    bulk := make([]byte, len)
    
    // Read the string data
    r.reader.Read(bulk)
    v.bulk = string(bulk)

    // Read the trailing \r\n
    r.readLine()

    return v, nil
}

// Marshal converts a Value into RESP wire format
// Used when sending responses back to clients
func (v Value) Marshal() []byte {
    // Choose appropriate marshaling method based on value type
    switch v.typ {
    case "array":
        return v.marshalArray()
    case "bulk":
        return v.marshalBulk()
    case "string":
        return v.marshalString()
    case "null":
        return v.marshallNull()
    case "error":
        return v.marshallError()
    default:
        return []byte{}
    }
}

// marshalString formats a RESP simple string
// Format: +<string>\r\n
func (v Value) marshalString() []byte {
    var bytes []byte
    bytes = append(bytes, STRING)            // Add type marker
    bytes = append(bytes, v.str...)          // Add string content
    bytes = append(bytes, '\r', '\n')        // Add CRLF
    return bytes
}

// marshalBulk formats a RESP bulk string
// Format: $<length>\r\n<string>\r\n
func (v Value) marshalBulk() []byte {
    var bytes []byte
    bytes = append(bytes, BULK)                          // Add type marker
    bytes = append(bytes, strconv.Itoa(len(v.bulk))...)  // Add length
    bytes = append(bytes, '\r', '\n')                    // Add CRLF
    bytes = append(bytes, v.bulk...)                     // Add string content
    bytes = append(bytes, '\r', '\n')                    // Add CRLF
    return bytes
}

// marshalArray formats a RESP array
// Format: *<length>\r\n<element-1>...<element-n>
func (v Value) marshalArray() []byte {
    len := len(v.array)
    var bytes []byte
    bytes = append(bytes, ARRAY)                     // Add type marker
    bytes = append(bytes, strconv.Itoa(len)...)      // Add array length
    bytes = append(bytes, '\r', '\n')                // Add CRLF
    
    // Marshal each array element
    for i := 0; i < len; i++ {
        bytes = append(bytes, v.array[i].Marshal()...)
    }
    
    return bytes
}

// marshallError formats a RESP error
// Format: -<error>\r\n
func (v Value) marshallError() []byte {
    var bytes []byte
    bytes = append(bytes, ERROR)             // Add type marker
    bytes = append(bytes, v.str...)          // Add error message
    bytes = append(bytes, '\r', '\n')        // Add CRLF
    return bytes
}

// marshallNull formats a RESP null value
// Format: $-1\r\n
func (v Value) marshallNull() []byte {
    return []byte("$-1\r\n")
}

// Writer wraps an io.Writer for writing RESP values
// Used to send responses back to Redis clients
type Writer struct {
    writer io.Writer
}

// NewWriter creates a new RESP writer
func NewWriter(w io.Writer) *Writer {
    return &Writer{writer: w}
}

// Write writes a Value in RESP format to the underlying writer
func (w *Writer) Write(v Value) error {
    // Marshal the value to RESP format
    var bytes = v.Marshal()
    
    // Write to the underlying writer
    _, err := w.writer.Write(bytes)
    if err != nil {
        return err
    }
    
    return nil
}