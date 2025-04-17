# Change Logs Sidecar

This repo is a small and streamlined implementation of a side car container for
capturing change logs entries in Crossplane providers.

The full design of this feature can be found in the [design doc](https://github.com/crossplane/crossplane/blob/main/design/one-pager-change-logs.md).

## Usage

This repository publishes release images to
`xpkg.crossplane.io/crossplane/changelogs-sidecar`. This image can then be
included as a sidecar container in a provider's pod via a
`DeploymentRuntimeConfig`.

When this container starts up, it starts a gRPC server that listens on a unix
domain socket at the default path of `/var/run/changelogs/changelogs.sock`. The
provider's main pod is expected to also start a gRPC client that connects and
sends requests over this socket. That gRPC client is then given to any managed
reconcilers that the provider uses to reconcile its managed resources. The
managed reconcilers will send change log entries using this client during
typical reconciliation events.

The gRPC server implementation in this repo accepts incoming change log entries
and simply writes them to `stdout` so they will be included in the provider
pod's logs.