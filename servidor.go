package main

import (
    "fmt"
    "net"
    // "encoding/gob"
)

func server() {
    s, err := net.Listen("tcp", ":9999")
    if err != nil {
        fmt.Println(err)
        return
    }
    for {
        c, err := s.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        go handleClient(c)
    }
}

func handleClient(c net.Conn) {

    c.Close()
}

func main() {
    go server()
    fmt.Scanln()
}
