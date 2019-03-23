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
To enable kujo on your jobs, make sure you set the annotations
`kujo.sphc.io: true`. Jobs where this is not set up will be ignored.

```bash
cat _examples/jobs.yaml | kujo
```

will output

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi-ignored
spec:
  backoffLimit: 4
  template:
    spec:
      containers:
      - command:
        - perl
        - -Mbignum=bpi
        - -wle
        - print bpi(2000)
        image: perl
        name: pi
      restartPolicy: Never
---
apiVersion: batch/v1
kind: Job
metadata:
  annotations:
    kujo.sphc.io: "true"
  name: pi-unique-54ccg6m6hb
spec:
  backoffLimit: 4
  template:
    spec:
      containers:
      - command:
        - perl
        - -Mbignum=bpi
        - -wle
        - print bpi(2000)
        image: perl
        name: pi
      restartPolicy: Never
```

## Future plans

### Operator

This behaviour would ideally live within an Operator where there is an
abstraction layer on top of Kubernetes' Job Resource. The Operator itself would
then take care of uniquely naming the individual Jobs, much like how the CronJob
Resource works.
