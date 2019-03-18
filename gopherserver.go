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

// RunServer runs the server, dispatching a goroutine to handle each connection.
func RunServer(listenAddr string, listenPort uint16, serverName string, gopherRoot string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenAddr, listenPort))
	if err != nil {
		fmt.Printf("couldn't listen on %s: %s", fmt.Sprintf("%s:%d", listenAddr, listenPort), err)
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
		if s == "\r\n" {
			conn.Write(directoryIndex(""))
		} else {
			req := strings.Trim(s, "\r\n")

			reqPath := filepath.Join(_gopherRoot, req)

			if strings.HasPrefix(reqPath, _gopherRoot) {
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
			} else {
				conn.Write(gopherErrorMessage("permission denied"))
			}
		}
		conn.Close()
	}
}

func directoryIndex(relativePath string) []byte {
	reqPath := filepath.Join(_gopherRoot, relativePath)

	var index []byte

	if reqPath != _gopherRoot {
		index = append(index, []byte(fmt.Sprintf("1..\t%s\t%s\t%d\r\n", filepath.Join(relativePath, ".."), _serverName, _listenPort))...)
	}

	if _, err := os.Stat(filepath.Join(reqPath, "gopher.index")); err == nil {
		index = append(index, sendDirectoryIndex(relativePath)...)
	} else if _, err := os.Stat(filepath.Join(reqPath, "Gophermap")); err == nil {
		index = append(index, sendDirectoryIndex(relativePath)...)
	} else {
		index = append(index, generateDirectoryIndex(relativePath)...)
	}

	index = append(index, []byte("i________________________________________________________________________________\r\n")...)
	index = append(index, []byte(fmt.Sprintf("i%80s", fmt.Sprintf("PocketRat/%s\r\n", Version)))...)

	return index
}

func sendDirectoryIndex(relativePath string) []byte {
	reqPath := filepath.Join(_gopherRoot, relativePath)

	if fileBytes, err := ioutil.ReadFile(filepath.Join(reqPath, "gopher.index")); err == nil {
		var index string

		file := string(fileBytes)
		lines := strings.Split(file, "\n")

		for _, l := range lines {
			line := strings.Split(l, "\t")
			if len(line) == 2 {
				selector := line[0]
				text := line[1]

				if strings.TrimSpace(selector) != "" {
					indexPath := filepath.Join(relativePath, selector)
					if file, err := os.Stat(filepath.Join(reqPath, indexPath)); err == nil {
						index += fmt.Sprintf("%s%s\t%s\t%s\t%d\r\n", getTypeSigil(file), text, indexPath, _serverName, _listenPort)
					}
				} else {
					index += fmt.Sprintf("i%s\t(message)\tpocketrat.invalid\t0\r\n", text)
				}
			}
		}

		return []byte(index)
	} else if fileBytes, err := ioutil.ReadFile(filepath.Join(reqPath, "Gophermap")); err == nil {
		var index string

		file := string(fileBytes)
		lines := strings.Split(file, "\n")

		for _, l := range lines {
			if !strings.ContainsAny(l, "\t") {
				index += "i" + l + "\t(message)\tpocketrat.invalid\t0"
			} else {
				index += l
			}
			index += "\r\n"
		}

		return []byte(index)
	}

	return generateDirectoryIndex(relativePath)
}

func generateDirectoryIndex(relativePath string) []byte {
	reqPath := filepath.Join(_gopherRoot, relativePath)
	files, err := ioutil.ReadDir(reqPath)

	if err == nil {
		if strings.HasPrefix(reqPath, _gopherRoot) {
			var index string
			var line string

			for _, file := range files {
				line = fmt.Sprintf("%s%s\t%s\t%s\t%d\r\n", getTypeSigil(file), file.Name(), filepath.Join(relativePath, file.Name()), _serverName, _listenPort)
				index += line
			}
			return []byte(index)
		} else if !strings.HasPrefix(reqPath, _gopherRoot) {
			return gopherErrorMessage("permission denied")
		}
	}

	return gopherError(err)
}

func getTypeSigil(file os.FileInfo) string {
	if !file.IsDir() {
		ext := filepath.Ext(file.Name())

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

	return "1"
}

func gopherError(err error) []byte {
	return []byte(gopherErrorMessage(err.Error()))
}

func gopherErrorMessage(errMsg string) []byte {
	return []byte(fmt.Sprintf("3%s\t(error)\tpocketrat.error.invalid\t0\r\n", errMsg))
}
