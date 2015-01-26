package main

import (
    "4gophers.com/cache/server"
    "log"
    "net"
)

func main() {
    // Запускаем наш сервер
    ln, err := net.Listen("tcp", ":11212")

    if err != nil {
        log.Println(err)
    }
    for {
        // Ждем подключения. При новом подключении ln.Accept()
        // будет зоздавать новое соединение conn net.Conn
        conn, err := ln.Accept()
        if err != nil {
            log.Println(err)
        }
        go server.ConnectionHandler(conn)
    }
}
