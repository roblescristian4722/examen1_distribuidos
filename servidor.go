package main

import (
    "fmt"
    "net"
    "encoding/gob"
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
    Ptype int
    Sender string
    Msg string
    File []byte
}


func server(ps *[]Petition) {
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
        go handleClient(c, ps)
    }
}

func handleClient(c net.Conn, ps *[]Petition) {
    p := &Petition{}
    dec := gob.NewDecoder(c)
    dec.Decode(p)
    if (*p).Ptype != EXIT && (*p).Ptype != SHOW_MESSAGES {
        *ps = append(*ps, *p)
    }
    c.Close()
}

func listMsg(ps *[]Petition) {
    fmt.Println()
    for _, p := range *ps {
        switch p.Ptype {
        case SEND_MESSAGE:
            fmt.Printf(">%s:\n", p.Sender)
            fmt.Printf("%s\n\n", p.Msg)
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
