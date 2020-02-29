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

---

## Master
