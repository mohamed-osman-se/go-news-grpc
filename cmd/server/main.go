package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	v1 "news-grpc/api/news/v1"
	ingrpc "news-grpc/internal/grpc"
	"news-grpc/internal/memstore"

	"github.com/bufbuild/protovalidate-go"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/proto"
)

func main() {
	// ---------- PROTOVALIDATE ----------
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("validator initialization: %v", err)
	}

	// ---------- UNARY INTERCEPTOR ----------
	unaryInterceptor := func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		// Validate ONLY protobuf messages
		if msg, ok := req.(proto.Message); ok {
			if err := validator.Validate(msg); err != nil {
				return nil, err
			}
		}

		return handler(ctx, req)
	}

	// ---------- STREAM INTERCEPTOR ----------
	streamInterceptor := func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Streaming validation must be handled inside service methods
		return handler(srv, ss)
	}

	// ---------- GRPC SERVER ----------
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptor),
		grpc.ChainStreamInterceptor(streamInterceptor),
	)

	// ---------- SERVICES ----------
	v1.RegisterNewsServiceServer(
		srv,
		ingrpc.NewServer(memstore.New()),
	)

	// ---------- HEALTH ----------
	healthSrv := health.NewServer()
	healthv1.RegisterHealthServer(srv, healthSrv)

	// ---------- GRACEFUL SHUTDOWN ----------
	grp, ctx := errgroup.WithContext(context.Background())

	grp.Go(func() error {
		lis, err := net.Listen("tcp", ":50051") //nolint:gosec
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}

		log.Println("gRPC server started on :50051")
		return srv.Serve(lis)
	})

	grp.Go(func() error {
		interceptSignals(ctx)
		log.Println("shutting down gRPC server...")
		healthSrv.Shutdown()
		srv.GracefulStop()
		return nil
	})

	if err := grp.Wait(); err != nil {
		log.Fatal("server shutdown error:", err)
	}
}

// ---------- SIGNAL HANDLING ----------
func interceptSignals(ctx context.Context) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		return
	case sig := <-sigc:
		log.Println("intercepted signal:", sig.String())
	}
}
