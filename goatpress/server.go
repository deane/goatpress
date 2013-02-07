package goatpress

import (
  "net"
  "os"
  "fmt"
  "time"
  "net/http"
)

const serverAddress = "localhost:4123"

var newPlayers = make(chan Player)
var removePlayers = make(chan string)

type Server struct {
  Tournament *Tournament
}

func newServer() *Server {
  gameType := newGameType(5, DefaultWordSet)
  tourney := newTournament(*gameType)
  randomPlayer := newInternalPlayer("Random", newRandomFinder(DefaultWordSet))
  tourney.RegisterPlayer(randomPlayer)
  return &Server{tourney}
}

func (c *Server) Run() {
  go c.RunWeb()
  c.RunTournament()
}

func homePage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (c *Server) RunWeb() {
  http.HandleFunc("/", homePage)
  http.ListenAndServe(":5123", nil)
}

func (c *Server) RunTournament() {
  listener, err := net.Listen("tcp", serverAddress)
  if err != nil {
    fmt.Printf("error listening:", err.Error())
    os.Exit(1)
  }
  go AcceptPlayers(listener)
  for {
    select {
    case newPlayer := <-newPlayers:
      if newPlayer.Name() != "" {
        fmt.Printf("Player Online: %s\n", newPlayer.Name())
        c.Tournament.RegisterPlayer(newPlayer)
      }
    case removePlayerName := <-removePlayers:
      if removePlayerName != "" {
        c.Tournament.DeregisterPlayer(removePlayerName)
      }
    default:
      if c.Tournament.Size() > 1 {
        c.Tournament.PlayMatch()
      } else {
        for _, player := range c.Tournament.Players {
          player.Ping()
        }
        time.Sleep(0.2*1e9)
      }
    }
    fmt.Printf("Players: %s\n", c.Tournament.PlayerList())
    time.Sleep(1)
  }
}

const serverSig = "goatpress<VERSION=1> ; \n"

func AcceptPlayers(listener net.Listener) {
  for {
    conn, err := listener.Accept()
    if err != nil {
      println("Error accept:", err.Error())
      return
    }
    conn.Write([]byte(serverSig))
    player := newClientPlayer(conn, removePlayers)
    if player != nil {
      newPlayers <- player
    }
  }
}


