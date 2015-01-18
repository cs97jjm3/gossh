# GOSSH
> A work in progress but fundementally "Golang + SSH connections = gossh"

My first attempt with Golang; A standalone binary that will periodically execute commands on remote linux boxes via SSH and present the output within a browser.

## Getting Started
Grab the source code
~~~
go get github.com/hiddenhippo/gossh
~~~
Navigate to the source code and type
~~~
go build
~~~
Before running the binary make sure you configure the config.ini file.  It's pretty self explanatory
~~~
[commands]
Disk=df -h | awk 'NR==2' | awk '{print $5}'
Free Memory=free -t | grep "buffers/cache" | awk '{print ($4 / 1000) "mb"}'
Uptime=cat /proc/uptime | awk '{print int($1 / 86400) "d " int(($1 % 86400)/3600) "h " int((($1 % 86400) % 3600 / 60)) "m "}'
Network=awk '/eth0/ {i++; rx[i]=$2; tx[i]=$10}; END{printf("%.2fMbps\n",((rx[2]-rx[1]) + (tx[2]-tx[1]))/131072) }'  <(cat /proc/net/dev; sleep 1; cat /proc/net/dev)
Load=uptime | awk -F'[a-z]:' '{ print $2}'
~~~

The commands section represents the commands that you want to execute on the remote server.  You may need to change "eth0" to the name of your network adaptor.
~~~
[yourserver]
url=sftp.myserver.com:22
username=username
password=password
~~~
You can have as many [yourserver] sections as you require, although each should have a unique name, requiring url, username and password.

Launch the executable and you'll have a gossh webserver running on `localhost:4000`
