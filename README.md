# gossh

web dashboard providing CPU, RAM, Network and Load metrics for remote hosts, where the remote hosts are accessible via ssh

go + ssh = gossh


t := gossh.NewServer(10 * time.Second)

err := t.AddMonitor("My SFTP Server", "sftp.mysftpserver.com:22", "username", "password")

if err != nil {

	log.Fatal(err)
	
}

<-t.Start()
