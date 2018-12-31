[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/yhagio/jiji/blob/master/LICENSE)

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
  },
  "mailgun": {
    "api_key": "key-b9cc7d8sdf897sdfkjsdf87dsfjlk8sdf",
    "public_api_key": "pubkey-2343249809823423jhjlksdfjhhf",
    "domain": "yourdomain.mailgun.org"
  },
  "dropbox": {
    "id": "sdf809sdn89dsf",
    "secret": "dsf789ds8f89sff2da",
    "auth_url": "https://www.dropbox.com/oauth2/authorize",
    "token_url": "https://api.dropboxapi.com/oauth2/token"
  }
}
```

```bash
chmod +x scripts/release.sh

go get            # Install dependencies
go run *.go       # Starts the application
godoc -http=:6060 # Documentation http://localhost:6060/pkg/jiji/
```

### OAuth2
https://godoc.org/golang.org/x/oauth2

### Dropbox SDK
Unofficial Dropbox SDK: https://github.com/dropbox/dropbox-sdk-go-unofficial

Chooser: https://www.dropbox.com/developers/chooser

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

### Digital Ocean setup - Deployment

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
sudo vi /etc/systemd/system/jiji_demo.service
```
fill with
```
[Unit]
Description=jiji_demo app

[Service]
WorkingDirectory=/root/app
ExecStart=/root/app/server -prod
Restart=always
RestartSec=30

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable jiji_demo.service
```

Create production config
```bash
touch /root/app/.config
vi /root/app/.config
```
and content is
```json
{
  "port": 3000,
  "env": "prod",
  "pepper": "SECRET_STUFF",
  "hmac_key": "SECRET_STUFF",
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "SECRET_STUFF",
    "name": "jiji_demo"
  },
  "mailgun": {
    "api_key": "your-api-key",
    "public_api_key": "your-public-key",
    "domain": "your-domain-setup-with-mailgun"
  }
}
```

### Create deployment script

```bash
mkdir scripts
touch scripts/releash.sh
```

**/scripts/release.sh** example (replace your own IP)
```bash
#!/bin/bash
cd "$GOPATH/src/jiji"

echo "==== Releasing jiji ===="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm jiji
echo "  Done!"


echo "  Deleting existing code..."
ssh root@123.78.123.156 "rm -rf /root/go/src/jiji"
echo "  Code deleted successfully!"


echo "  Uploading code..."
# The \ at the end of the line tells bash that our
# command isn't done and wraps to the next line.
rsync -avr --exclude '.git/*' --exclude 'tmp/*' --exclude 'images/*' ./ \
  root@123.78.123.156:/root/go/src/jiji/
echo "  Code uploaded successfully!"


echo "  Go getting deps..."
ssh root@123.78.123.156 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
ssh root@123.78.123.156 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/mux"
ssh root@123.78.123.156 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/schema"
ssh root@123.78.123.156 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/lib/pq"
ssh root@123.78.123.156 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh root@123.78.123.156 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/csrf"


echo "  Building the code on remote server..."
ssh root@123.78.123.156 'export GOPATH=/root/go; \
  cd /root/app; \
  /usr/local/go/bin/go build -o ./server \
    $GOPATH/src/jiji/*.go'
echo "  Code built successfully!"


echo "  Moving assets..."
ssh root@123.78.123.156 "cd /root/app; \
  cp -R /root/go/src/jiji/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@123.78.123.156 "cd /root/app; \
  cp -R /root/go/src/jiji/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@123.78.123.156 "cd /root/app; \
  cp /root/go/src/jiji/Caddyfile ."
echo "  Views moved successfully!"


echo "  Restarting the server..."
ssh root@123.78.123.156 "sudo service jiji_demo restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@123.78.123.156 "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing jiji ===="
```


Deploy command
```bash
./scripts/release.sh
```

Also update `resetBaseURL` in `./email.mailgun.go` 
