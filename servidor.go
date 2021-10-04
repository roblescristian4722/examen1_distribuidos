package main

import (
	"encoding/gob"
	"fmt"
	"net"
    "os"
    "strings"
    "strconv"
    "bufio"
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
    File []byte
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
    if err != nil { fmt.Println(err); return }
    for {
        c, err := s.Accept()
        if err != nil { fmt.Println(err); continue}
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
            case "add": active = append(active, msg); break
            case "call":
                for _, v := range active {
                    gob.NewEncoder(v.Conn).Encode((*ps)[len(*ps) - 1])
                }
                (*ps)[len(*ps) - 1].File = []byte{}
                break
            }
        }
    }
}

func createFile(p *Petition) {
    if _, err := os.Stat("server_files/"); os.IsNotExist(err) {
        err := os.Mkdir("server_files/", 0777)
        if err != nil { fmt.Println(err); return }
    }
    path := "server_files/" + (*p).Msg
    file, err := os.Create(path)
    if err != nil { fmt.Println(err); return }

    file.Write((*p).File)
}

func readFile(p *Petition) []byte {
    file, err := os.Open("server_files/" + (*p).Msg)
    if err != nil { fmt.Println(err); return []byte{} }
    stat, err := file.Stat()
    if err != nil { fmt.Println(err); return []byte{} }
    bs := make([]byte, stat.Size())
    file.Read(bs)
    return bs
}

func backup(ps *[]Petition) {
    file, err := os.OpenFile("server.backup", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
    defer file.Close()
    if err != nil { fmt.Println("No se ha hecho un respaldo con anterioridad"); return }
    for _, v := range *ps {
        data := strconv.FormatInt(int64(v.Type), 10) + "|" + v.Sender + "|" + v.Msg + "\n"
        file.WriteString(data)
    }
}

func restore(ps *[]Petition) {
    file, err := os.Open("server.backup")
    if err != nil { fmt.Println(err); return }
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.Split(scanner.Text(), "|")
        t, _ := strconv.ParseInt(line[0], 10, 64)
        *ps = append(*ps, Petition{
            Type: int(t),
            Sender: line[1],
            Msg: line[2],
        })
    }
}

func handleClient(id uint, c net.Conn, cMsg chan Connection, ps *[]Petition) {
    cMsg <- Connection{ id, "add", c }
    for {
        p := &Petition{}
        err := gob.NewDecoder(c).Decode(p)
        if err == nil {
            switch (*p).Type {
            case SEND_MESSAGE, SEND_FILE:
                if (*p).Type == SEND_FILE {
                    createFile(p)
                }
                *ps = append(*ps, *p)
                cMsg <- Connection{ id, "call", c }
                break
            case SHOW_MESSAGES:
                tmpPs := *ps
                for i, v := range tmpPs {
                    if v.Type == SEND_FILE {
                        bs := readFile(&v)
                        tmpPs[i].File = bs
                    }
                }
                gob.NewEncoder(c).Encode(&tmpPs)
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

func main() {
    op := -1
    ps := []Petition{}
    restore(&ps)
    go server(&ps)
    for op != EXIT {
        fmt.Println("\n----------------Group Chat Server-------------------")
        fmt.Println("Seleccione una opción:")
        fmt.Println(LIST_MESSAGES, ") Mostrar mensajes/archivos recibidos")
        fmt.Println(BACKUP, ") Respaldar mensajes/archivos en un archivo")
        fmt.Print(EXIT, " ) Salir\n>> ")
        fmt.Scanln(&op)
        switch op {
            case LIST_MESSAGES:
                listMsg(&ps)
                break
            case BACKUP:
                backup(&ps)
                break
            case EXIT: return
            default: fmt.Println("Opción no válida, vuelva a intentarlo")
        }
    }
}
