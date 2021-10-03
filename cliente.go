package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
    "encoding/gob"
)

const (
    SEND_MESSAGE = iota + 1
    SEND_FILE
    SHOW_MESSAGES
    EXIT = 0
)

type Petition struct {
    Type int
    Sender string
    Msg string
}

func sendMsg(c net.Conn, scanner *bufio.Scanner, username string) {
    p := &Petition{}
    fmt.Print("Mensaje a enviar: ")
    scanner.Scan()
    (*p).Msg = scanner.Text()
    (*p).Type = SEND_MESSAGE
    (*p).Sender = username
    gob.NewEncoder(c).Encode(p)
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

func client(conn chan net.Conn, username string, scanner *bufio.Scanner, ps *[]Petition) {
    op := -1
    c, err := net.Dial("tcp", ":9999")
    if err != nil {
        fmt.Println(err)
        return
    }
    // Primera conexión con el server: se copian todos los msg a ps
    gob.NewEncoder(c).Encode(&Petition{ Type: SHOW_MESSAGES })
    gob.NewDecoder(c).Decode(ps)
    // Goroutine que obtiene nuevos mensajes de otros clientes
    go func() {
        for {
            p := &Petition{}
            err := gob.NewDecoder(c).Decode(p)
            if err == nil { *ps = append(*ps, *p) }
        }
    }()
    for op != EXIT {
        fmt.Println("\n---------------------------------------------")
        fmt.Printf("¡Hola %s! Selecciona una opción:\n", username)
        fmt.Println(SEND_MESSAGE, ") Enviar un mensaje")
        fmt.Println(SEND_FILE, ") Enviar un archivo")
        fmt.Println(SHOW_MESSAGES, ") Mostrar tus mensajes")
        fmt.Println(EXIT, ") Salir")
        fmt.Print(">> ")
        fmt.Scan(&op)
        switch op {
        case SEND_MESSAGE:
            sendMsg(c, scanner, username)
            break
        case SEND_FILE:

            break
        case SHOW_MESSAGES:
            listMsg(ps)
            break
        case EXIT:
            gob.NewEncoder(c).Encode(&Petition{})
            conn <- c
            return
        default:
            fmt.Println("Opción no válida, vuelva a intentarlo")
        }
    }
}

func main() {
    ps := []Petition{}
    conn := make(chan net.Conn)
    scanner := bufio.NewScanner(os.Stdin)

    fmt.Print("Ingrese su nombre de usuario: ")
    scanner.Scan()
    username := scanner.Text()

    go client(conn, username, scanner, &ps)
    c := <-conn

    c.Close()
    // Se termina la conexión con el servidor y la ejecución del cliente termina
}

