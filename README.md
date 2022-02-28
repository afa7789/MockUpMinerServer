# MockUp Miner Server
A MockUp miner server is a backend code that will receive FTP requests , written in the format of stratum protocols.
Stratum protocols is the agreed upon format to designate jsonRPC calls that miner servers send to miners client and vice-versa
To make use of these code base you gotta have installed not only Golang but Telnet to make the calls for it.

# Running this code

## how to run

```sh
  git clone https://github.com/afa7789/MockUpMinerServer mockup
  cd mockup
  # docker-compose up
  docker-compose up
  # now we gotta seed our database
  # this bellow runs the script as if the database was a remote one, but it's actually in the docker :)
  psql -h localhost -d postgres -U postgres -f etc/postgres/0_init.sql
  # for now we are only having the ids I inserted in the db
  go run . -port=8000 # if you do not set port it will use 8080
```

```python
telnet localhost 8000 # 8080 here is the port change to reflect the above
# use the following lines in the order so you can see the authorization system working
{"id": 1, "method": "mining.subscribe", "params": []}
{"params": ["slush.miner1", "password"], "id": 1, "method": "mining.authorize"}
{"id": 1, "method": "mining.subscribe", "params": []}
# to create a new worker , if you don't have a id ( let's say so) 
{"id": NULL, "method": "mining.new_worker", "params": []} # the id doesn't matter here it will create a new one
# it will return a success with your new id.
```

```bash
docker exec -it postgres sh
psql -U postgres
# entered in the db now
\c
#wanna see the tables ?
\dt
select * from entries;
select * from miners;
select * from subscriptions;
```

Project organization:

```sh
.
├── cmd
│   └── start.go # the starting command
├── docker-compose.yaml
├── domain
│   ├── content.go # shared general structs 
│   ├── postgres.go # shared postgres structs
│   └── stratum.go # shared stratum structs
├── etc
│   └── postgres
│       └── 0_init.sql # script to initialize the db
├── go.mod
├── go.sum
├── Instructions.pdf # challenge description
├── internal
│   ├── middleware
│   │   └── middleware.go # middleware
│   ├── postgres # postgres related code 
│   │   ├── blacklist.go
│   │   ├── entry.go
│   │   ├── miner.go
│   │   ├── pg.go
│   │   └── subscriptions.go
│   ├── server
│   │   └── server.go # ftp server
│   └── stratum
│       └── stratum.go # stractum handling over the connection , have a function for each stratum call
├── main.go
└── README.md
```

# Comments about the code:

This is only a mockup so it does not work to actually mine crypto.
We would still need the client golang code to be running in the computer to comunicate with us, and the other stratum methods finished.

I have choosen to use DB clients without GORM to actually write some SQL , get a habit and see how low level and raw I can do SQL queries on Go.

