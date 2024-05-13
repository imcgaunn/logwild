# logwild

writes json logs in a format similar to RCS at a configurable rate in messages per-second
ensures that each message contains unique contents to make sure logging backends don't sample
events.

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
