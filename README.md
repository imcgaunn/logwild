# logwild

writes json logs in a format similar to RCS at a configurable rate in messages per-second
ensures that each message contains unique contents to make sure logging backends don't sample
events.

## running example

first, create the 'observe' docker network that the containers will use to communicate. There
is a script defined that can do this for you:

```bash
logwild on  main [!?] via 🐳 orbstack via 🐹 v1.22.3 on ☁️  (us-east-1)
✦ ❯ scripts/util/setup_docker_network.sh
a2099b63b53a3fb6f4a7251cc09213ca095fb17ddb2885ef4e0e5c6785bb9f08
```

you can confirm the network was created by issuing `docker network ls`, e.g.:

```bash
logwild on  main [!?] via 🐳 orbstack via 🐹 v1.22.3 on ☁️  (us-east-1)
✦ ❯ docker network ls
NETWORK ID     NAME      DRIVER    SCOPE
1df4780d4551   acs       bridge    local
b67f3623cbd2   bridge    bridge    local
b1e231de1a95   host      host      local
b299faf05206   none      null      local
a2099b63b53a   observe   bridge    local
```

now, you should be able to run the full example in a separate tmux session by
executing `just run-in-panels`, i.e.:

```bash
logwild on  main [!?] via 🐳 orbstack via 🐹 v1.22.3 on ☁️  (us-east-1)
✦ ❯ just run-in-panels
scripts/run_everything_in_panels.sh
killing existing session logwild_demo
program running in other window waiting
```

when the example has started, you should see the message 'program running in other window waiting'
do not press 'return' in this window, or cleanup logic will be executed and tear down the example.

to view the output from the example, switch to the session called `logwild_demo`:

```bash
tmux switch-client -t logwild_demo # or however else you'd like to switch to a different session
```

### accessing openobserve webUI

one of the containers in this example hosts the `openobserve` service, which provides
visualizations for otel/openmetrics, traces, and logs. you can access the frontend
for openobserve once the example is running by navigating to the following URL
a web browser:

```text
http://localhost:5080
```

the default username and password are as follows (found in `run_standalone_observe_backend.sh`):

```text
user email: root@example.com
password: Complexpass#123
```

### triggering log generation

once the `logwild` server starts up, you can issue a `GET` request to its `/loggen` endpoint
and it should start streaming many logs to the configured logger, which without changes
will point to `os.Stdout` (in the logwild container - not in `curl` output). For example:

```bash
logwild on  main [!?] via 🐳 orbstack via 🐹 v1.22.3 on ☁️  (us-east-1)
❯ curl localhost:8888/loggen
{}%
```

## IMPORTANT ACKNOWLEDGMENTS

this is mostly not my code. I started from the venerable https://github.com/stefanprodan/podinfo microservice template
which is an excellent reference, but it has a lot to process for new go developers.
microservice template, replaced its use of uber's 'zap' logger with go's native structured logging library 'slog',
and pruned much of the functionality that I wasn't interested in emphasizing.

### references / links

1. [podinfo microservice template](https://github.com/stefanprodan/podinfo)
2. [openobserve](https://github.com/openobserve/openobserve)
3. [zincsearch](https://zincsearch-docs.zinc.dev/) (predecessor to openobserve - alternative to elasticsearch)
4. [opentelemetry collector](https://opentelemetry.io/docs/collector/)
