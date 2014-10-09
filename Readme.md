# Сталин
Сталин сука простой, он просто берет и херачит данные и в graphite и в redis и в другой riemann-server, либо в stdout.

```
Usage of stalin:
  -Bind="0.0.0.0:5555": Listen riemann events
  -GraphiteAdr="": Send metrics to graphite
  -RedisAdr="": Connect to redis, like '127.0.0.1:6379'
  -RedisList="riemann.events": Redis event list
  -ResendAdr="": Resend to next riemann
  -alsologtostderr=false: log to standard error as well as files
  -log_backtrace_at=:0: when logging hits line file:N, emit a stack trace
  -log_dir="": If non-empty, write log files in this directory
  -logtostderr=false: log to standard error instead of files
  -stderrthreshold=0: logs at or above this threshold go to stderr
  -v=0: log level for V logs
  -vmodule=: comma-separated list of pattern=N settings for file-filtered logging
``` 

# Todo:
* дропать старые сообщения по переполнению счетчика
* наследование
