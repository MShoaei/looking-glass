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

---

## Master
