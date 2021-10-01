package main

import (
    "fmt"
    "net"
)

const (
    SEND_MESSAGE = iota + 1
    SEND_FILE
    SHOW_MESSAGES
    EXIT
)

func client(conn chan net.Conn) {
    var op int
    c, err := net.Dial("tcp", ":9999")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    for op != EXIT {
    }
}

func main() {
    conn := make(chan net.Conn)
    go client(conn)
    fmt.Scanln()
    // Se termina la conexión con el servidor
}

