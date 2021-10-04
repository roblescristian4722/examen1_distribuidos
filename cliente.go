package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
    "encoding/gob"
    "strings"
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
    File []byte
}

func sendMsg(c net.Conn, scanner *bufio.Scanner, username string) {
    p := &Petition{}
    fmt.Print("Mensaje a enviar: ")
    scanner.Scan()
    (*p).Msg = scanner.Text()
    (*p).Type = SEND_MESSAGE
    (*p).Sender = username
    err := gob.NewEncoder(c).Encode(p)
    if err == nil {
        fmt.Println("Mensaje enviado con éxito")
    } else { fmt.Println(err) }
}

func listMsg(ps *[]Petition) {
    for _, p := range *ps {
        switch p.Type {
        case SEND_MESSAGE:
            fmt.Printf("\n>%s:\n%s\n\n", p.Sender, p.Msg)
            break
        case SEND_FILE:
            f := strings.Split(p.Msg, ".")
            ext := f[len(f) - 1]
            var t string
            switch ext {
                case "jpg", "jpeg", "png", "raw": t = "Archivo (Imagen)"; break
                case "mp4", "avi", "amv", "webm", "flv": t = "Archivo (Video)"; break
                case "mp3", "3gp", "flac", "m4a": t = "Audio (Audio)"; break
                default: t = "Archivo (" + ext + ")"
            }
            fmt.Printf("\n>%s:\n%s: %s\n\n", p.Sender, t, p.Msg)
            break
        }
    }
}

func readFile(path string) []byte {
    file, err := os.Open(path)
    if err != nil { fmt.Println(err); return []byte{} }
    stat, err := file.Stat()
    if err != nil { fmt.Println(err); return []byte{} }
    bs := make([]byte, stat.Size())
    file.Read(bs)
    return bs
}

func sendFile(c net.Conn, scanner *bufio.Scanner, username string) {
    fmt.Println("Archivo a enviar: ")
    scanner.Scan()
    path := scanner.Text()
    bs := readFile(path)
    if len(bs) == 0 { return }
    pathS := strings.Split(path, "/")
    err := gob.NewEncoder(c).Encode(&Petition{ SEND_FILE, username, pathS[len(pathS) - 1], bs })
    if err == nil {
        fmt.Println("Archivo enviado con éxito")
    } else { fmt.Println(err) }

}

func recieveFile(c net.Conn, username string, p *Petition) {
    if _, err := os.Stat("client_files/"); os.IsNotExist(err) {
        err := os.Mkdir("client_files/", 0777)
        if err != nil { fmt.Println(err); return }
    }
    if _, err := os.Stat("client_files/" + username); os.IsNotExist(err) {
        err := os.Mkdir("client_files/" + username, 0777)
        if err != nil { fmt.Println(err); return }
    }
    path := "client_files/" + username + "/" + p.Msg
    file, err := os.Create(path)
    if err != nil { fmt.Println(err); return }
    file.Write(p.File)
    (*p).File = []byte{}
}

func client(conn chan net.Conn, username string, scanner *bufio.Scanner, ps *[]Petition) {
    op := -1
    c, err := net.Dial("tcp", ":9999")
    if err != nil { fmt.Println(err); return }
    // Primera conexión con el server: se copian todos los msg a ps
    gob.NewEncoder(c).Encode(&Petition{ Type: SHOW_MESSAGES })
    gob.NewDecoder(c).Decode(ps)
    for _, v := range *ps {
        if v.Type == SEND_FILE {
            recieveFile(c, username, &v)
        }
    }
    // Goroutine que obtiene nuevos mensajes de otros clientes
    go func() {
        for {
            p := &Petition{}
            err := gob.NewDecoder(c).Decode(p)
            if err == nil {
                if p.Type == SEND_FILE {
                    recieveFile(c, username, p)
                }
                *ps = append(*ps, *p)
            }
        }
    }()
    for op != EXIT {
        fmt.Println("\n---------------------------------------------")
        fmt.Printf("¡Hola %s! Selecciona una opción:\n", username)
        fmt.Println(SEND_MESSAGE, ") Enviar un mensaje")
        fmt.Println(SEND_FILE, ") Enviar un archivo")
        fmt.Println(SHOW_MESSAGES, ") Mostrar tus mensajes")
        fmt.Print(EXIT, " ) Salir\n>> ")
        fmt.Scan(&op)
        switch op {
        case SEND_MESSAGE:
            sendMsg(c, scanner, username)
            break
        case SEND_FILE:
            sendFile(c, scanner, username)
            break
        case SHOW_MESSAGES:
            listMsg(ps)
            break
        case EXIT:
            gob.NewEncoder(c).Encode(&Petition{})
            conn <- c
            return
        default: fmt.Println("Opción no válida, vuelva a intentarlo")
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
}
