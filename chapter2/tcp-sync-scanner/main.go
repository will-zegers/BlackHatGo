package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "sort"
    "strconv"
    "strings"
)

func worker(ports, results chan int, host string) {
    for port := range ports {
        address := fmt.Sprintf("%s:%d", host, port)
        conn, err := net.Dial("tcp", address)
        if err != nil {
            results <- 0
            continue
        }
        conn.Close()
        results <- port
    }
}

func parseCommaSeparatedPorts(portsStr string) []uint16 {
    var ports []uint16
    for _, i := range strings.Split(portsStr, ",") {
        if strings.Contains(i, "-") {
            ports = append(ports, parseDashSeparatedPorts(i)...)
        }
        port, _ := strconv.ParseInt(i, 10, 16)
        ports = append(ports, uint16(port))
    }
    return ports
}

func parseDashSeparatedPorts(portsStr string) []uint16 {
    var ports []uint16
    portRange := strings.Split(portsStr, "-")
    minPort, _ := strconv.ParseInt(portRange[0], 10, 16)
    maxPort, _ := strconv.ParseInt(portRange[1], 10, 16)

    for i := uint16(minPort); i <= uint16(maxPort); i++ {
        ports = append(ports, i)
    }

    return ports
}

func parsePortList(portsStr string) []uint16 {
    var ports []uint16

    if strings.Contains(portsStr, ",") {
        ports = parseCommaSeparatedPorts(portsStr)
    } else if strings.Contains(portsStr, "-") {
        ports = parseDashSeparatedPorts(portsStr)
    } else {
        port, _ := strconv.ParseInt(portsStr, 10, 16)
        ports = []uint16 {uint16(port)}
    }

    return ports
}

func main() {
    hostPtr := flag.String("host", "", "Host IP or FQDN")
    portsPtr := flag.String("ports", "", "Single port number (eg. 80), dash-separated port range (eg. 100-200), comma-separated list of ports (53,80,443), or combination (eg. 20-23,53,80)")
    workerCountPtr := flag.Int("workers", 100, "Number of worker threads to run")
    flag.Parse()

    if *hostPtr == "" || *portsPtr == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    ports := make(chan int, *workerCountPtr)
    results := make(chan int)
    var openports []int

    for i := 0; i < cap(ports); i++ {
        go worker(ports, results, *hostPtr)
    }

    go func() {
        for i := 1; i <= 1024; i++ {
            ports <- i
        }
    }()

    for i := 0; i < 1024; i++ {
        port := <-results
        if port != 0 {
            openports = append(openports, port)
        }
    }

    close(ports)
    close(results)
    sort.Ints(openports)
    for _, port := range openports {
        fmt.Printf("[+] %d is open\n", port)
    }
}
