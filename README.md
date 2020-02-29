# Looking Glass Project

Derak cloud entrance project.

- Stack
  - iris web framework (Go)
  - Postgres

## Worker

### Plugin

plugins should be located in a directory with the plugin's name and must be compiled as plugins the compiled plugin file should have the same name as the parent directory and end with `.so`.

`./worker/run.sh` compiles all the plugins at `./worker/*/` and starts the worker.

Plugin's initial state(enabled/disabled) is up to the author of the plugin. (**this might change in the future**)

All plugins must expose a variable named `P` which implements the Runner interface

### API

/

/home

/login

/plugins

/tags

### Database

Worker won't start unless it successfully connects to database. environment variables SHOULD be set or default values will be used

`DB_USER` is the username used to connect to database. Default: `test`

`DB_PASSWORD` is the password used to connect to database. Default: `testpassword`

`DB_NAME` is the database name. Default: `test_db`

`DB_PORT` is the port to cennect to. Default: `5432`

---

## Master
