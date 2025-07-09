package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Wuli-Giao-Giao/tools/server"
	grpc "google.golang.org/grpc"
)

type GRPCServer struct {
	srv  *grpc.Server
	addr string
}

func NewGRPCServer(addr string, grpcServer *grpc.Server) server.Server {
	return &GRPCServer{
		srv:  grpcServer,
		addr: addr,
	}
}

func (g *GRPCServer) Start() error {
	listener, err := net.Listen("tcp", g.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", g.addr, err)
	}
	return g.srv.Serve(listener)
}

func (g *GRPCServer) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		g.srv.GracefulStop()
		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		g.srv.Stop()
		return ctx.Err()
	case <-done:
		return nil
	}
}
