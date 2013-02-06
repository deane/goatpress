package main

import (
  "fmt"
  "goatpress"
  "flag"
)

func usage() {
  fmt.Printf("usage: goatpress-cli OPTIONS (demo|server|client|web)\n")
}

func main() {
  playerName := flag.String("name", "", "name")
  flag.Parse()
  fmt.Printf("Connecting as player %s\n", *playerName)
  if len(flag.Args()) < 1 {
    usage()
  } else {
    command := flag.Args()[0]
    if command == "demo" {
      goatpress.Demo()
    } else if command == "server" {
      goatpress.ServerStart()
    } else if command == "client" {
      goatpress.ClientStart(*playerName)
    } else if command == "web" {
    } else {
      usage()
    }
  }
}

