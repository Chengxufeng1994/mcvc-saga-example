package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/usecase"
	"github.com/Chengxufeng1994/go-saga-example/common/bootstrap"
	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	application     string
	bootstrapConfig *bootstrap.BootstrapConfig
	authService     usecase.AuthUseCase
	Srv             *grpc.Server
	pb.UnimplementedAuthServiceServer
}

func New(bootstrapConfig *bootstrap.BootstrapConfig, authService usecase.AuthUseCase) *GrpcServer {
	grpcSrv := &GrpcServer{
		application:     bootstrapConfig.Application,
		bootstrapConfig: bootstrapConfig,
		authService:     authService,
	}

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 8), // increase to 8 MB (default: 4 MB)
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // terminate the connection if a client pings more than once every 5 seconds
			PermitWithoutStream: true,            // allow pings even when there are no active streams
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second,  // if a client is idle for 15 seconds, send a GOAWAY
			MaxConnectionAge:      600 * time.Second, // if any connection is alive for more than maxConnectionAge, send a GOAWAY
			MaxConnectionAgeGrace: 5 * time.Second,   // allow 5 seconds for pending RPCs to complete before forcibly closing connections
			Time:                  5 * time.Second,   // ping the client if it is idle for 5 seconds to ensure the connection is still active
			Timeout:               1 * time.Second,   // wait 1 second for the ping ack before assuming the connection is dead
		}),
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	}

	grpcSrv.Srv = grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(grpcSrv.Srv, grpcSrv)
	reflection.Register(grpcSrv.Srv)
	return grpcSrv
}

func (s *GrpcServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	resp, err := s.authService.VerifyToken(ctx, &dto.VerifyTokenRequest{AccessToken: req.AccessToken})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error: %v", err),
		)
	}

	return &pb.VerifyTokenResponse{
		UserId: resp.UserId,
	}, nil
}

func (s *GrpcServer) Run() error {
	addr := fmt.Sprintf("%s:%d", s.bootstrapConfig.Grpc.Host, s.bootstrapConfig.Grpc.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	config.ContextLogger.Infoln("grpc.Run listening on", s.bootstrapConfig.Grpc.Port)
	if err := s.Srv.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *GrpcServer) GracefulShutdown(ctx context.Context) {
	config.ContextLogger.Infoln("grpc.GracefulShutdown")
	s.Srv.GracefulStop()
}
