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

Plugin mainly has 4 methods: `Enable`, `Disable`, `Status` and `Run`. `Run` is the **only** asynchronous method. disabling a plugin has no effect on previously started tasks but all queued tasks will fail.

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

`DB_PORT` is the port to connect to. Default: `5432`

### tasks

POST request to /plugins/execute/{plugin} will start a goroutine to execute the command.
Later calls to this route will cause the previous results to be lost and it's caller's responsibility to get the latest result before performing subsequent tasks.

---

## Master
