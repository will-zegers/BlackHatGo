package main

import (
    "io"
    "log"
    "net"
)

// echo is a handler funtion that simply echoes received data
func echo(conn net.Conn) {
    defer conn.Close()

    // Create a buffer to store received data
    b := make([]byte, 512)
    for {
        // Receive data via conn.Read into a buffer.
        size, err := conn.Read(b[0:])
        if err == io.EOF {
            log.Println("Client disconnected")
            break
        } else if err != nil {
            log.Println("Unexpected error")
            break
        }
        log.Printf("Recieved %d bytes: %s\n", size, string(b))

        // Send data via conn.Write
        log.Println("Writing data")
        if _, err := conn.Write(b[0:size]); err != nil {
            log.Fatalln("Unable to write data")
        }
    }
}

func main() {
    // Bind to TCP port 20080 on all interfaces
    listener, err := net.Listen("tcp", ":20080")
    if err != nil {
        log.Fatalln("Unable to bind to port")
    }
    log.Println("Listening on 0.0.0.0:20080")
    for {
        // Wait for connection. Create net.Conn on connection established.
        conn, err := listener.Accept()
        if log.Println("Received connection"); err != nil {
            log.Fatalln("Unable to accept connection")
        }
        // Handle the connection. Using goroutine for concurrency.
        go echo(conn)
    }
}
