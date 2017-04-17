# SeriousApiarist

SeriousApiarist provides an API exposed to the Continuous Delivery (CD) pipeline
that handels builds, tests, vulnerability scans, releases, and deploys in a
safe and controlled fashion.


Project groups and committers are whitelisted. All API activity has traces and an audit log.


Deploys are protected by 2FA and successful deploys trigger a Slack alert.


## API endpoints
The endpoints and parameters are abstracted by the CLI commands available to you
in the CD pipeline container.
file.

- `/build/:group/:repo`
- `/test/:group/:repo`
- `/scan/:group/:repo`
- `/release/:group/:repo`
- `/deploy/:group/:repo`

### Parameters
- `service` (optional): applicable if you have multiple docker containers in your
repo. Service is expected to be the name of a folder in your repo that contains
a `Dockerfile`.
- `test` (optional): Test parameter is required if you're running tests. Test parameter is provided to your container as the first argument, so your entrypoint
must expect it in that format.
- `imageTag` (optional): imageTag parameter will tag your release image so it looks like <image name>:<tag>-<pipeline id>
- `ref`: Branch/ref is provided from the CD pipeline.
- `commit`: Commit hash is provided from the CD pipeline.
- `pipeline`: Pipeline number is provided from the CD pipeline.
