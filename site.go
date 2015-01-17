//Sites can be considered the SSH end-points e.g.
//the servers you wish to be monitoring

package gossh

import (
	"golang.org/x/crypto/ssh"
	"time"
)

//definition of a site
type site struct {
	Url      string
	Username string
	Server   string
	Password string
	Sleep    time.Duration
	config   *ssh.ClientConfig
	client   *ssh.Client
	commands []command
}

func (s *site) Configure() (err error) {
	s.config = &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
	}

	//associate the commands
	s.commands = make([]command, 0)
	s.commands = append(s.commands, &memorycommand{})
	s.commands = append(s.commands, &uptimecommand{})
	s.commands = append(s.commands, &networkcommand{})
	s.commands = append(s.commands, &diskcommand{})

	//this needs moving
	return s.connect()
}

//establish a connection to the remote site(server)
func (s *site) connect() (err error) {
	s.client, err = ssh.Dial("tcp", s.Url, s.config)
	return err
}

//For the given site iterate through it's commands collecting
//each command output and then pushing the result onto the out channel
func (s *site) Poll(out chan<- result) {

	//the results of a single server
	results := result{Server: s.Server, Values: make(map[string]string, 0)}

	//collect the results
	for i := range s.commands {
		session, _ := s.client.NewSession()
		//execute the command passing the session
		s.commands[i].Execute(session)
		session.Close()
		results.Values[s.commands[i].Name()] = s.commands[i].Results()
	}

	//pump the results out
	out <- results

	//sleep for the next cycle
	s.sleep()

	//and start again *** recursive?
	s.Poll(out)
}

//sleep the cycle for the required time
func (s *site) sleep() {
	time.Sleep(s.Sleep)
}
