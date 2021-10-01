# datadog-parser [![Build Status](https://github.com/paskal/datadog-parser/workflows/CI%20Build/badge.svg)](https://github.com/paskal/datadog-parser/actions?query=workflow%3A%22CI+Build%22)

datadog-parser takes CSV-formatted logs from input or file and produces stats and alerts based on them.

## Run instructions

As a prerequisite, you need to have [Docker](https://www.docker.com/products/docker-desktop) installed on the machine, and you should run `docker build . -t paskal/data-parser:latest` once.

### Read log from stdin

To feed datadog-parser with logs using stdin, run it using `docker run -i paskal/data-parser:latest`. For example:

```shell
# read from pipe and rewrite standard parameters
cat ./sample.csv | docker run -i paskal/data-parser:latest --alert_threshold_per_sec 1000 --alert_window 1h
# collect data from input and write the results to the file and not stdout to prevent mixing input and output data
docker run -i paskal/data-parser:latest > processed_logs.txt
```

### Read log from file

To run datadog-parser against [sample.csv](sample.csv), run `docker-compose up`. You can alter the `volumes` block to mount a different file and pass application options using the `environment` block.

### Application parameters

| Command line   | Environment  | Default | Description            |
| ---------------| -------------| --------| -----------------------|
| filepath       | FILEPATH     |         | csv file path, stdin is used if not specified |
| alert_window   | ALERT_WINDOW | `2m`    | alert windows          |
| alert_threshold_per_sec | ALERT_THRESHOLD_PER_SEC] | `10` |  threshold for alert, requests per second |
| help           |              |         | shows the help message |

## Restrictions

- When using the file as input, newly appended lines to the log are processed. However, you'll need to restart the application if the log file reduces the size, like when `truncate -s0` was used to clean it.
- As it was specified not to rely on the host machine time, alerts re-evaluation happens only when new log entries are appended. If the last log entry provided to the program is in an alert state, and then there will be no logs, it will be stuck in alerting state.
- HTTP methods are ignored and not counted separately. Response statuses are collected but not shown in the stats.
- Flapping of the alert like the following from the sample is not prevented:
  ```
  2019-02-07 21:14:04 +0000 UTC: Alert GREEN, ~9.97 hits per second which is lower than 10 (1197 total) in the last 2m0s
  2019-02-07 21:14:04 +0000 UTC: Alert RED, ~10.01 hits per second which is higher than 10 (1201 total) in the last 2m0s
  2019-02-07 21:14:05 +0000 UTC: Alert GREEN, ~9.89 hits per second which is lower than 10 (1187 total) in the last 2m0s
  ```

## Additional notes

I've taken Dockerfile and docker-compose file, GitHub Actions pipeline, and linter setting from [rlb-stats](https://github.com/umputun/rlb-stats) stats collector ([link](https://stats.radio-t.com/rlb/)) which was my first Go project back in 2017. That project was also my first collaboration with [@umputun](https://github.com/umputun), which enormously increased my development skills over the years.

Code has 91.4% tests coverage which comes at the cost of hacks here and there to make them work properly.

I've spent roughly six hours writing that because Go has solid guarantees about lack of surprises in the runtime, which comes at the cost of higher writing time. I decided to develop a decent solution rather than hacking something together in Python strictly within the timeframe.

It was an interesting problem to solve.
