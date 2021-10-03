package main

import (
	"encoding/gob"
	"fmt"
	"net"
)

const (
    SEND_MESSAGE = iota + 1
    SEND_FILE
    SHOW_MESSAGES
    LIST_MESSAGES = 1
    BACKUP = 2
    EXIT = 0
)

type Petition struct {
    Type int
    Sender string
    Msg string
}

type Connection struct {
    Id uint
    Type string
    Conn net.Conn
}

func server(ps *[]Petition) {
    id := uint(0)
    cMsg := make(chan Connection)
    go handleConn(cMsg, ps)
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
        go handleClient(id, c, cMsg, ps)
        id++
    }
}

func handleConn(cMsg chan Connection, ps *[]Petition) {
    active := []Connection{}
    for {
        select {
        case msg := <-cMsg:
            switch msg.Type {
            case "kill":
                for i, v := range active {
                    if v.Id == msg.Id {
                        msg.Conn.Close()
                        active = append(active[:i], active[i + 1:]...)
                    }
                }
                break
            case "add": active = append(active, msg)
                break
            case "call":
                for _, v := range active {
                    gob.NewEncoder(v.Conn).Encode((*ps)[len(*ps) - 1])
                }
                break
            }
        }
    }
}

func handleClient(id uint, c net.Conn, cMsg chan Connection, ps *[]Petition) {
    cMsg <- Connection{ id, "add", c }
    for {
        p := &Petition{}
        err := gob.NewDecoder(c).Decode(p)
        if err == nil {
            switch (*p).Type {
            case SEND_MESSAGE:
                *ps = append(*ps, *p)
                cMsg <- Connection{ id, "call", c }
                break
            case SHOW_MESSAGES:
                gob.NewEncoder(c).Encode(ps)
                break
            case EXIT:
                cMsg <- Connection{ id, "kill", c }
                return
            }
        }
    }
}

func listMsg(ps *[]Petition) {
    for _, p := range *ps {
        switch p.Type {
        case SEND_MESSAGE:
            fmt.Printf("\n>%s:\n%s\n\n", p.Sender, p.Msg)
            break
        case SEND_FILE:
            break
        }
    }
}

func main() {
    op := -1
    ps := []Petition{}

    go server(&ps)
    for op != EXIT {
        fmt.Println("\n----------------Group Chat Server-------------------")
        fmt.Println("Seleccione una opción:")
        fmt.Println(LIST_MESSAGES, ") Mostrar mensajes/archivos recibidos")
        fmt.Println(BACKUP, ") Respaldar mensajes/archivos en un archivo")
        fmt.Println(EXIT, ") Salir")
        fmt.Print(">> ")
        fmt.Scanln(&op)
        switch op {
            case LIST_MESSAGES:
                listMsg(&ps)
                break
            case BACKUP:
                break
            case EXIT:
                return
            default:
                fmt.Println("Opción no válida, vuelva a intentarlo")
        }
    }
}
