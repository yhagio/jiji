# JIJI

Web App built with Go and Postgres

<img src="https://cdn-images-1.medium.com/max/600/1*i2skbfmDsHayHhqPfwt6pA.png" height="200px">
<img src="https://raw.githubusercontent.com/docker-library/docs/01c12653951b2fe592c1f93a13b4e289ada0e3a1/postgres/logo.png" height="180px">

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

```bash
go get            # Install dependencies
go run main.go    # Starts the application
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
