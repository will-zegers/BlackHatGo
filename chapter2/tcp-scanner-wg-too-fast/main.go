package main

import (
    "fmt"
    "net"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    host := "127.0.0.1"

    for i := 1; i < 65535; i++ {
        wg.Add(1)
        go func(j int) {
            defer wg.Done()
            address := fmt.Sprintf("%s:%d", host, j)
            conn, err := net.Dial("tcp", address)
            if err != nil {
                // port is closed or filtered
                return
            }
            conn.Close()
            fmt.Printf("Port %d open\n", j)
        }(i)
    }
    wg.Wait()
}
