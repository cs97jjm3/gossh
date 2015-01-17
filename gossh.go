package gossh

import (
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
func (s *Server) AddMonitor(servername, url, username, password string) (err error) {
	s.sites[servername] = &site{Server: servername, Url: url, Username: username, Password: password, Sleep: s.pollInterval}
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
