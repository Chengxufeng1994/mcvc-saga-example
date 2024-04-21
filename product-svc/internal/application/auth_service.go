package application

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/common/pb"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/grpc/client"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"
)

type AuthService struct {
	authClient *client.AuthConn
}

func NewAuthService(authConn *client.AuthConn) usecase.AuthUseCase {
	return &AuthService{
		authClient: authConn,
	}
}

// VerifyToken implements usecase.AuthUseCase.
func (a *AuthService) VerifyToken(ctx context.Context, accessToken string) (*dto.VerifyTokenResponse, error) {
	cli := pb.NewAuthServiceClient(a.authClient.Conn())
	res, err := cli.VerifyToken(ctx, &pb.VerifyTokenRequest{AccessToken: accessToken})
	if err != nil {
		return nil, err
	}

	return &dto.VerifyTokenResponse{
		UserId:    res.UserId,
		IsExpired: res.IsExpired,
	}, nil
}
