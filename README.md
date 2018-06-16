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
