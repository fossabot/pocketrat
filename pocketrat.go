package main

import (
    "fmt"
)

const version = "0.0.0"

func main() {
    fmt.Println("PocketRat version", version)
    config, err := LoadConfiguration("config.json")
    if (err == nil) {
        fmt.Println("listening on", config.ListenAddr, "port", config.ListenPort)
        fmt.Println("server name is", config.ServerName, "and gopher root is", config.GopherRoot)
        RunServer(config.ListenAddr, config.ListenPort, config.ServerName, config.GopherRoot)
    } else {
        panic(err)
    }
}