package grpc

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/domain"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/repository"
	infragrpc "github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/server/grpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type GrpcAuthRepository struct {
	AuthConn *infragrpc.AuthConn
}

func NewGrpcAuthRepository(authConn *infragrpc.AuthConn) repository.AuthRepository {
	return &GrpcAuthRepository{
		AuthConn: authConn,
	}
}

// VerifyToken implements repository.AuthRepository.
func (g *GrpcAuthRepository) VerifyToken(ctx context.Context, accessToken string) (*domain.Auth, error) {
	ctx, span := otel.Tracer("purchase").Start(ctx, "Verify Token")
	defer span.End()
	span.SetAttributes(attribute.String("access_token", accessToken))

	cli := pb.NewAuthServiceClient(g.AuthConn.Conn())
	res, err := cli.VerifyToken(ctx, &pb.VerifyTokenRequest{AccessToken: accessToken})
	if err != nil {
		return nil, err
	}

	return &domain.Auth{
		UserId:    res.UserId,
		IsExpired: res.IsExpired,
	}, nil
}
