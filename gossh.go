package main

import (
	"github.com/vaughan0/go-ini"
	"log"
	"os"
	"strconv"
	"time"
)

var applicationResults map[string]result

//master server object
type Server struct {
	pollInterval time.Duration
	sites        map[string]*site
}

//single instance of a result gathered from a server
type result struct {
	Server string
	Values map[string]string
}

func stateMonitor() chan<- result {
	updates := make(chan result)

	go func() {
		for {
			select {
			case r := <-updates:
				applicationResults[r.Server] = r
			}
		}
	}()
	return updates
}

//initialize a new Server
func NewServer(pollInterval time.Duration) *Server {
	return &Server{sites: make(map[string]*site), pollInterval: pollInterval}
}

//Add a site to be monitored
func (s *Server) AddSite(servername, url, username, password string, commands []command) (err error) {
	s.sites[servername] = &site{Server: servername, Url: url, Username: username, Password: password, Sleep: s.pollInterval, Commands: commands}
	return s.sites[servername].Configure()
}

//main poll cycle
func (s *Server) Start() <-chan bool {

	applicationResults = make(map[string]result, 0)
	out := make(chan bool)

	//create object to output completed work
	status := stateMonitor()

	go func() {
		launchBrowser()
	}()

	//launch our monitors concurrently
	for _, v := range s.sites {
		go v.Poll(status)
	}

	return out
}

func main() {
	l := log.New(os.Stdout, "[gossh] ", 0)

	file, err := ini.LoadFile("config.ini")
	if err != nil {
		l.Fatal(err)
	}

	tmpPollInterval, ok := file.Get("settings", "pollInterval")
	if !ok {
		l.Fatal(err)
	}

	pollInterval, err := strconv.Atoi(tmpPollInterval)
	if err != nil {
		l.Fatal(err)
	}

	s := NewServer(time.Duration(pollInterval) * time.Second)

	for name, _ := range file {
		if name == "settings" || name == "commands" {
			continue
		}
		url, _ := file.Get(name, "url")
		username, _ := file.Get(name, "username")
		password, _ := file.Get(name, "password")

		if len(url) == 0 || len(username) == 0 || len(password) == 0 {
			l.Printf("server " + name + " skipped ")
			continue
		}

		//in the future we might want different commands per site being monitored
		//such as a database server being different from a webserver
		c := make([]command, 0)
		for key, value := range file["commands"] {
			c = append(c, command{Name: key, Cmd: value})
		}
		s.AddSite(name, url, username, password, c)
		l.Printf(name + " added")
	}

	<-s.Start()

}
