package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
    "encoding/gob"
    "strconv"
)

const (
    SEND_MESSAGE = iota + 1
    SEND_FILE
    SHOW_MESSAGES
    EXIT = 0
)

type Petition struct {
    Ptype int
    Sender string
    Msg string
    File []byte
}

func sendMsg(c net.Conn, scanner *bufio.Scanner, username string) {
    p := &Petition{}
    fmt.Print("Mensaje a enviar: ")
    scanner.Scan()
    (*p).Msg = scanner.Text()
    (*p).Ptype = SEND_MESSAGE
    (*p).Sender = username
    gob.NewEncoder(c).Encode(p)
}

func listMsg(ps *[]Petition) {
    for _, p := range *ps {
        switch p.Ptype {
        case SEND_MESSAGE:
            fmt.Printf("\n>%s:\n%s\n\n", p.Sender, p.Msg)
            break
        case SEND_FILE:
            break
        }
    }
}

func showMsg(c net.Conn, ps *[]Petition) {
    p := Petition{ Ptype: SHOW_MESSAGES, Msg: strconv.FormatUint(uint64(len(*ps)), 10) }
    gob.NewEncoder(c).Encode(p)
    psTmp := []Petition{}
    gob.NewDecoder(c).Decode(&psTmp)
    *ps = append(*ps, psTmp...)
    c.Close()
    listMsg(ps)
}

func client(conn chan string, username string, scanner *bufio.Scanner, ps *[]Petition) {
    op := -1
    for op != EXIT {
        c, err := net.Dial("tcp", ":9999")
        if err != nil {
            fmt.Println(err)
            return
        }

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
            showMsg(c, ps)
            break
        case EXIT:
            c.Close()
            conn <- "kill"
            return
        default:
            fmt.Println("Opción no válida, vuelva a intentarlo")
        }
        c.Close()
    }
}

func main() {
    ps := []Petition{}
    conn := make(chan string)
    scanner := bufio.NewScanner(os.Stdin)

    fmt.Print("Ingrese su nombre de usuario: ")
    scanner.Scan()
    username := scanner.Text()

    go client(conn, username, scanner, &ps)
    <-conn
    // Se termina la conexión con el servidor y la ejecución del cliente termina
}

