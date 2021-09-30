# datadog-parser [![Build Status](https://github.com/paskal/datadog-parser/workflows/CI%20Build/badge.svg)](https://github.com/paskal/datadog-parser/actions?query=workflow%3A%22CI+Build%22)

datadog-parser takes CSV-formatted logs from input or file and produces stats and alerts based on them.

## Run instructions

As a prerequisite, you need to have [Docker](https://www.docker.com/products/docker-desktop) installed on the machine.

### Read log from stdin

To feed datadog-parser with log using stdin, first build it using `docker build . -t paskal/data-parser:latest` and then run it using `docker run -i paskal/data-parser:latest`.

Example of usage in the pipe:

```shell
cat ./sample.csv | docker run -i paskal/data-parser:latest --alert_threshold_per_sec 1000 --alert_window 1h
```

### Read log from file

To run datadog-parser against `sample.csv`, run `docker-compose up`. You can alter the `volumes` block to mount a different file and pass application options using the `environment` block.

### Application parameters

| Command line   | Environment  | Default | Description            |
| ---------------| -------------| --------| -----------------------|
| filepath       | FILEPATH     |         | csv file path, stdin is used if not specified |
| alert_window   | ALERT_WINDOW | `2m`    | alert windows          |
| alert_threshold_per_sec | ALERT_THRESHOLD_PER_SEC] | `10` |  threshold for alert, requests per second |
| help           |              |         | shows the help message |

## Restrictions

Newly appended lines to the log are processed. However, you'll need to restart the application if the log file reduces the size, like when `truncate -s0` was used to clean it.

## Additional notes

I've taken Dockerfile and docker-compose file, GitHub Actions pipeline, and linter setting from [rlb-stats](https://github.com/umputun/rlb-stats) stats collector ([link](https://stats.radio-t.com/rlb/)) which was my first Go project back in 2017. That project was also my first collaboration with [@umputun](https://github.com/umputun), which enormously increased my development skills over the years.
