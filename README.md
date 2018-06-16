# JIJI

Web App built with Go and Postgres

<img src="https://cdn-images-1.medium.com/max/600/1*i2skbfmDsHayHhqPfwt6pA.png" height="200px" style="display: inline-block">
<img src="https://raw.githubusercontent.com/docker-library/docs/01c12653951b2fe592c1f93a13b4e289ada0e3a1/postgres/logo.png" height="180px" style="display: inline-block">

### DB setup - Postgres (on Ubuntu)

```
sudo -i -u postgres

psql
> CREATE DATABASE jiji_dev;
> CREATE USER jiji_dev_user WITH PASSWORD '123test';
> GRANT ALL PRIVILEGES ON DATABASE jiji_dev to jiji_dev_user;
> ALTER USER jiji_dev_user WITH SUPERUSER;
> ALTER ROLE jiji_dev_user CREATEROLE CREATEDB;
```

## Dev setup

Create a file `.config`
```json
{
  "port": 3000,
  "env": "dev",
  "pepper": "secret-random-string",
  "hmac_key": "secret-hmac-key",
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "jiji_dev_user",
    "password": "your-password",
    "name": "jiji_dev"
  }
}
```

```bash
go get            # Install dependencies
go run *.go       # Starts the application
godoc -http=:6060 # Documentation http://localhost:6060/pkg/jiji/
```

### (Optional) Hot reloading

Go hot-reloading, so no need to stop and restart server

```bash
go get github.com/pilu/fresh
```

In this app repo, run `fresh`

```bash
fresh
```

### Other commands
```bash 
go run *.go --help   # Check if flag is provided

go build -o app *.go # Build a binary named app
./app --help

go run *.go -prod    # Production, ensures to use .config
```

### Digital Ocean setup

After ssh into the droplet server you created

Setup DB
```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib

sudo -u postgres psql
> ALTER USER postgres WITH ENCRYPTED PASSWORD 'YOUR_PASSWORD';
# ctrl + d to quit

vi /etc/postgresql/10/main/pg_hba.conf
# Replace 
# from: local   all   postgres   peer
# to:   local   all   postgres   md5

sudo service postgresql restart

psql -U postgres
> CREATE DATABASE jiji_demo;
```

Setup Go
```bash
curl -O https://storage.googleapis.com/golang/go1.10.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
mkdir -p go/src

# Check Go version
go version
/usr/local/go/bin/go version
```

Setup Caddy
```bash
cd ~/go/src
go get -u github.com/mholt/caddy
go get -u github.com/caddyserver/builds

# Build Caddy
cd ~/go/src/github.com/mholt/caddy/caddy
go run build.go -goos=linux -goarch=amd64

./caddy # Test Caddy

cp ./caddy /usr/local/bin/ # Copy Caddy binary

sudo vi /etc/systemd/system/caddy.service
```
and fill with
```
[Unit]
Description=caddy server for serving jiji-demo

[Service]
WorkingDirectory=/root/app
ExecStart=/usr/local/bin/caddy -email your@email.com
Restart=always
RestartSec=120
LimitNOFILE=8192

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable caddy.service
```

```bash
mkdir /root/app
sudo service caddy restart
journalctl -r # view logs from our services, and the -r flag says to display them in reverse order so we see the newest logs first
```

```bash
sudo service caddy stop
```