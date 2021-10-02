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
    Ptype int
    Dest string
    Msg string
}

func client(conn chan net.Conn, username string, scanner *bufio.Scanner) {
    op := -1
    for op != EXIT {
        c, err := net.Dial("tcp", ":9999")
        if err != nil {
            fmt.Println(err)
            return
        }
        enc := gob.NewEncoder(c)

        fmt.Println("\n¡Bienvenido ", username, "! Selecciona una opción:")
        fmt.Println(SEND_MESSAGE, ") Enviar un mensaje")
        fmt.Println(SEND_FILE, ") Enviar un archivo")
        fmt.Println(SHOW_MESSAGES, ") Mostrar tus mensajes")
        fmt.Println(EXIT, ") Salir")
        fmt.Print(">> ")
        fmt.Scan(&op)
        fmt.Println(op, SEND_MESSAGE)
        switch op {
        case SEND_MESSAGE:
            p := &Petition{}
            fmt.Print("Destinatario: ")
            scanner.Scan()
            (*p).Dest = scanner.Text()
            fmt.Print("Mensaje a enviar: ")
            scanner.Scan()
            (*p).Msg = scanner.Text()
            (*p).Ptype = 0
            enc.Encode(p)
            break
        case SEND_FILE:
            break
        case SHOW_MESSAGES:
            break
        case EXIT:
            conn <- c
            return
        default:
            fmt.Println("Opción no válida, vuelva a intentarlo")
        }
        fmt.Println("loop")
    }
}

func main() {
    conn := make(chan net.Conn)
    scanner := bufio.NewScanner(os.Stdin)

    fmt.Print("Ingrese su nombre de usuario: ")
    scanner.Scan()
    username := scanner.Text()

    go client(conn, username, scanner)
    c := <-conn
    c.Close()
    // Se termina la conexión con el servidor
}

