package gossh

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"math"
	"strconv"
	"strings"
)

type command interface {
	Execute(*ssh.Session)
	Name() string
	Results() string
}

type diskcommand struct {
	results string
}

func (d *diskcommand) Name() string {
	return "Disk"
}

func (d *diskcommand) Execute(s *ssh.Session) {
	defer s.Close()
	var stdoutBuf bytes.Buffer
	s.Stdout = &stdoutBuf
	s.Run("df -h | awk 'NR==2' | awk '{print $5}'")
	d.results = stdoutBuf.String()
}

func (d *diskcommand) Results() string {
	return strings.TrimSpace(d.results)
}

type networkcommand struct {
	results string
}

func (n *networkcommand) Name() string {
	return "Network"
}

func (n *networkcommand) Execute(s *ssh.Session) {
	defer s.Close()
	var stdoutBuf bytes.Buffer
	s.Stdout = &stdoutBuf
	s.Run("awk '/eth0/ {i++; rx[i]=$2; tx[i]=$10}; END{print (((rx[2]-rx[1]) + (tx[2]-tx[1])*8)/1000000)}'  <(cat /proc/net/dev; sleep 1; cat /proc/net/dev)")
	n.results = stdoutBuf.String()
}

func (n *networkcommand) Results() string {
	return strings.TrimSpace(n.results) + " Mbps"
}

type uptimecommand struct {
	results string
}

func (u *uptimecommand) Results() string {
	s := u.results
	i, _ := strconv.ParseFloat(s[0:strings.Index(s, " ")], 64)

	days := math.Floor(i / 86400)
	hours := math.Floor(math.Mod(i, 86400) / 3600)
	minutes := math.Floor(math.Mod(math.Mod(i, 86400), 3600) / 60)
	seconds := math.Floor(math.Mod(math.Mod(math.Mod(i, 86400), 3600), 60))
	return fmt.Sprintf("%6.0fd %6.0fh %6.0fm %6.0fs", days, hours, minutes, seconds)
}

func (u *uptimecommand) Name() string {
	return "Uptime"
}

func (u *uptimecommand) Execute(s *ssh.Session) {
	defer s.Close()
	var stdoutBuf bytes.Buffer
	s.Stdout = &stdoutBuf
	s.Run("cat /proc/uptime")
	u.results = stdoutBuf.String()
}

type memorycommand struct {
	results string
}

func (m *memorycommand) Name() string {
	return "Memory"
}

func (m *memorycommand) Results() string {
	return strings.TrimSpace(m.results)
}

func (m *memorycommand) Execute(s *ssh.Session) {
	defer s.Close()
	var stdoutBuf bytes.Buffer
	s.Stdout = &stdoutBuf
	s.Run("free -t | grep \"buffers/cache\" | awk '{print $4}'")
	m.results = stdoutBuf.String()
}
