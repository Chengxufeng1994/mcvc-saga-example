package application

import "github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/usecase"

type Application struct {
	AuthService     usecase.AuthUseCase
	PurchaseService usecase.PurchaseUseCase
}

func New(
	authService usecase.AuthUseCase,
	purchaseService usecase.PurchaseUseCase) *Application {
	return &Application{
		AuthService:     authService,
		PurchaseService: purchaseService,
	}
}
