package main

import (
    "fmt"
    "io/ioutil"
    "mime"
    "net"
    "os"
    "path/filepath"
    "strings"
)

var _serverName string
var _gopherRoot string
var _listenPort uint16

// RunServer runs the server.  Duh.
func RunServer(listenAddr string, listenPort uint16, serverName string, gopherRoot string) {
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenAddr, listenPort))
    if err != nil {
        panic(err)
    }
    
    _serverName = serverName
    _gopherRoot = gopherRoot
    _listenPort = listenPort
    
    for {
        conn, err := listener.Accept()
        if err == nil {
            go handleConnection(conn)
        } else {
            fmt.Println("error accepting connection:", err)
        }
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			conn.Close()
			return
		}

		s := string(buf[0:n])
        if (s == "\r\n") {
            conn.Write(directoryIndex(""))
        } else {
            req := strings.Trim(s, "\r\n")
            
            // TODO: This is vulnerable to directory traversal attacks
            if stat, err := os.Stat(filepath.Join(_gopherRoot, req)); err == nil {
                if !stat.IsDir() {
                    bytes, err := ioutil.ReadFile(filepath.Join(_gopherRoot, req))
                    if err == nil {
                        conn.Write(bytes)
                    } else {
                        conn.Write(gopherError(err))
                    }
                } else {
                    conn.Write(directoryIndex(req))
                }
            } else {
                conn.Write(gopherErrorMessage(req + " does not exist"))
            }
        }
        conn.Write([]byte("."))
        conn.Close()
    }
}

func directoryIndex(relativePath string) []byte {
    files, err := ioutil.ReadDir(filepath.Join(_gopherRoot, relativePath))
    
    if err == nil {
        var index string
        var line string
        
        for _, file := range files {
            typeSigil := (func() string { if file.IsDir() { return "1" }; return getTypeSigil(file.Name()) })()
            line = fmt.Sprintf("%s%s\t%s\t%s\t%d\r\n", typeSigil, file.Name(), filepath.Join(relativePath, file.Name()), _serverName, _listenPort)
            index += line
        }
        return []byte(index)
    }
    
    return gopherError(err)
}

func getTypeSigil(fileName string) string {
    ext := filepath.Ext(fileName)
    
    // check for a few known types by extension before resorting to MIME types
    switch ext {
        case ".gif":
            return "g"
        case ".html", ".htm":
            return "h"
        case ".hqx", ".hcx":
            return "4"
        case ".uue", ".uu":
            return "6"
        case ".txt":
            return "0"
    }

    // if we've made it to this point, it wasn't a known extension.
    // resort to mime types    
    mimeType := mime.TypeByExtension(ext)
    
    // mime "helpfully" provides charset parameters for text/* types,
    // but we don't really need that
    if strings.Index(mimeType, ";") > -1 {
        mimeType = mimeType[0:strings.Index(mimeType, ";")]
    }
    
    if strings.HasPrefix(mimeType, "image/") {
        return "I"
    }

    if strings.HasPrefix(mimeType, "text/") {
        return "0"
    }
    
    if strings.HasPrefix(mimeType, "audio/") {
        return "s"
    }
    
    return "9"
}

func gopherError(err error) []byte {
    return []byte(gopherErrorMessage(err.Error()))
}

func gopherErrorMessage(errMsg string) []byte {
    return []byte(fmt.Sprintf("3%s\t(error)\terror.invalid\t0\r\n", errMsg))
}