# GOSSH

> A work in progress but fundementally "Golang + SSH connections = gossh"

My first attempt with Golang; Connect to linux boxes using SSH to obtain some basic statistics on a repeat interval which are then presented on a HTML dashboard.  There's not much at the moment but currently Disk, Memory, Network and Uptime are reported, in varying degrees of accuracy.

## Getting Started
After installing Go and setting up your GOPATH, create your first .go file. We'll call it server.go.
~~~ go
package main

import (
    "github.com/hiddenhippo/gossh"
    "time"
)

func main() {
    //create a new server which polls destination sites every 10 seconds
    g := gossh.NewServer(10 * time.Second)
    g.AddSite("My reliable sftp server", "sftp.myserver.com:22", "username", "password")
    <-g.Start()
}
~~~
Then install the gossh package:
~~~
go get github.com/hiddenhippo/gossh
~~~
Then run your server:
~~~
go run server.go
~~~

You will now have a goosh webserver running on `localhost:4000`
