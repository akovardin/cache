package main

import (
    // "encoding/gob"
    "fmt"
    "log"
    "net"
)

func main() {
    // mycache()
    memcached()
}

func memcached() {
    fmt.Println("start client")
    conn, err := net.Dial("tcp", "localhost:11211")
    if err != nil {
        log.Fatal("Connection error", err)
    }

    conn.Write([]byte("set gokey 0 0 9\r\nsomevalue\r\n"))
    conn.Close()
    fmt.Println("done")
}

func mycache() {
    fmt.Println("start client")
    conn, err := net.Dial("tcp", "localhost:11212")
    if err != nil {
        log.Fatal("Connection error", err)
    }

    conn.Write([]byte("hi"))
    conn.Close()
    fmt.Println("done")
}
