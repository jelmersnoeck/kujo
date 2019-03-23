# Kujo

Kujo - Kubernetes Unique Jobs - helps with the creation of unique job names for
your Kubernetes Job Resources.

It's aim is to provide a way to use the same job definition and apply a unique
name based on configuration changes. This can be used when using a Job object to
run migrations for example, where the config - your migrations - will change.

By using Kujo, you keep a history of the jobs that has been run, without the
need to set up a system to remove old jobs so they can be mutated.

## Usage

Usage is based on a piping system where Kujo expects data to be sent on stdin.

```bash
cat examples/job.yaml | kujo | kubectl apply -f -
```

## Future plans

### Operator

This behaviour would ideally live within an Operator where there is an
abstraction layer on top of Kubernetes' Job Resource. The Operator itself would
then take care of uniquely naming the individual Jobs, much like how the CronJob
Resource works.
