// Package main implements the main logic of the change logs sidecar.
package main

import (
	"io/fs"
	"net"
	"os"

	"github.com/alecthomas/kong"
	"google.golang.org/grpc"

	"github.com/crossplane/changelogs-sidecar/server"
	changelogs "github.com/crossplane/crossplane-runtime/apis/changelogs/proto/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
)

// CLI is the command line interface for the change logs sidecar server.
type CLI struct {
	SocketPath string `help:"Path to create a Unix domain socket for change logs gRPC." default:"/var/run/changelogs/changelogs.sock"`
}

// Run is the main entry point for the change logs sidecar server.
func (c *CLI) Run() error {
	// clean up any existing socket file from previous runs
	if err := os.Remove(c.SocketPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return errors.Wrapf(err, "failed to remove any existing unix domain socket at %s", c.SocketPath)
	}

	// start listening on the unix domain socket
	lis, err := net.Listen("unix", c.SocketPath)
	if err != nil {
		return errors.Wrapf(err, "failed to listen to unix domain socket at %s", c.SocketPath)
	}

	// initialize the gRPC server
	s := server.Server{}
	grpcServer := grpc.NewServer()
	changelogs.RegisterChangeLogServiceServer(grpcServer, &s)

	// start serving gRPC requests
	if err := grpcServer.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve")
	}

	return nil
}

func main() {
	ctx := kong.Parse(&CLI{}, kong.Description("Change logs sidecar server"))
	ctx.FatalIfErrorf(ctx.Run())
}
