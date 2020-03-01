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

    requset: GET
    response: 302 Redirect

/home

    request: GET 
    response: application/json 200 OK
    body: {
      message string
    }

/login
  
    request: POST application/json
    body: {
      username string
      password string
    }
    response: application/json 200 OK 
    body: {
      Authorization string
    }

/plugins

- status/{plugin}

      request: GET
      response: application/json 200 OK
      body: {
        pluginStatus string
      }

- enable/{plugin}

      request: POST application/json
      response: application/json 200 OK
      body: {
        success bool
      }

- disable/{plugin}

      request: POST application/json
      response: application/json 200 OK
      body: {
        success bool
      }

/tags

- /

      request: GET
      response: application/json 200 OK
      body: {
        tags []string
      }

- /{tag}

      request: PUT
      response: application/json 200 OK
      body: {
        message string
      }

- /{tag}

      request: DELETE
      response: application/json 204 OK

/task

- /status

      request: GET
      response: application/json 200 OK
      body: {
        cmdStatus string
      }

- /result

      request: GET
      response: application/json 200 OK
      body: {
        result result
      }

### Database

Worker won't start unless it successfully connects to database. environment variables `DATABASE_URL` MUST be set.

### tasks

POST request to /plugins/execute/{plugin} will start a goroutine to execute the command.
Later calls to this route will cause the previous results to be lost and it's caller's responsibility to get the latest result before performing subsequent tasks.

### Benchmark

- Spec:
  - RAM: 16GB
  - CPU: Intel(R) Core(TM) i5-6500 CPU @ 3.20GHz

```bash
$ bombardier --header="Authorization:Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODMwMTg4ODEsImp0aSI6MX0.Zwef_lwan3-En_ZhDUR0yrEmzkqNGfgho0qpCVnnyJvmmqrVSlqxEvB3rSPRUx_8DoA4eH9ZiUFvaDCkrH7rQA" -c 125 -n 1000000 localhost:8080/


Bombarding http://localhost:8080/ with 1000000 request(s) using 125 connection(s)
 1000000 / 1000000 [=======================================] 100.00% 32229/s 31s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec     32404.78    5006.10   43895.57
  Latency        3.86ms     1.28ms    45.29ms
  HTTP codes:
    1xx - 0, 2xx - 0, 3xx - 1000000, 4xx - 0, 5xx - 0
    others - 0
  Throughput:    12.63MB/s
```

```bash
$ bombardier --header="Authorization:Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1ODMwMTg4ODEsImp0aSI6MX0.Zwef_lwan3-En_ZhDUR0yrEmzkqNGfgho0qpCVnnyJvmmqrVSlqxEvB3rSPRUx_8DoA4eH9ZiUFvaDCkrH7rQA" -m PUT -c 125 -n 1000000 localhost:8080/tags/worker1

Bombarding http://localhost:8080/tags/worker1 with 1000000 request(s) using 125 connection(s)
 1000000 / 1000000 [=======================================] 100.00% 30453/s 32s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec     30473.54    4613.93   41012.83
  Latency        4.10ms     1.28ms    44.86ms
  HTTP codes:
    1xx - 0, 2xx - 1000000, 3xx - 0, 4xx - 0, 5xx - 0
    others - 0
  Throughput:    14.29MB/s
```
