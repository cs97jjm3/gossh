package main

import (
	"bytes"
	"golang.org/x/crypto/ssh"
)

type command struct {
	Name string
	Cmd  string
}

func (c *command) Execute(s *ssh.Session) string {
	defer s.Close()
	var stdoutBuf bytes.Buffer
	s.Stdout = &stdoutBuf
	s.Run(c.Cmd)
	return stdoutBuf.String()
}
