package main

import (
	"flag"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	changelogs "github.com/crossplane/crossplane-runtime/apis/changelogs/proto/v1alpha1"

	"github.com/crossplane/changelogs-sidecar/server"
)

func main() {
	// get socket file path, allowing user to specify via command line flag
	socketPath := flag.String("socket-path", "/var/run/changelogs/changelogs.sock", "Path to create a Unix domain socket for change logs gRPC")
	flag.Parse()

	// clean up any existing socket file from previous runs
	if err := os.Remove(*socketPath); err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to remove any existing unix domain socket at %s: %+v", *socketPath, err)
	}

	// start listening on the unix domain socket
	lis, err := net.Listen("unix", *socketPath)
	if err != nil {
		log.Fatalf("failed to listen to unix domain socket at %s: %+v", *socketPath, err)
	}

	// initialize the gRPC server
	s := server.Server{}
	grpcServer := grpc.NewServer()
	changelogs.RegisterChangeLogServiceServer(grpcServer, &s)

	// start serving gRPC requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}
