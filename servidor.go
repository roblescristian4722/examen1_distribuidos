package main

import (
    "fmt"
    "net"
    "encoding/gob"
)

type Petition struct {
    Ptype int
    Dest string
    Msg string
}

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
    p := &Petition{}
    dec := gob.NewDecoder(c)
    dec.Decode(p)
    fmt.Println(p)
    c.Close()
}

func main() {
    go server()
    fmt.Scanln()
}
