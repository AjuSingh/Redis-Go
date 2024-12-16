package main //required for executable programs in Go
import (
	"fmt" //I/O operations
	"io" //I/O interfaces
	"net" //network operations
	"os" //OS functionality
)
func main() {
	fmt.Println("Listening on port :6379")
	// Create a new server
	l, err := net.Listen("tcp", ":6379") //creates TCP network listener
	if err != nil {
		fmt.Println(err)
		return
	}
	// Listen for connections
	conn, err := l.Accept() //The program waits here until a client connects - blocking call
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close() //Schdules conn.Close() to be run after the current function returns
	for { //infinite loop for continous client communication
		buf := make([]byte, 1024) //byte slice
		// read message from client
		_, err = conn.Read(buf) //discards number of bytes read, reads data from connection into buffer
		if err != nil {
			if err == io.EOF { //connection closed by Client
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}
		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n")) //converts string to bytesplice
	}
}