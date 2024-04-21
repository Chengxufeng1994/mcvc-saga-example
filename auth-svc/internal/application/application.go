package application

import "github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/usecase"

type Application struct {
	AuthService usecase.AuthUseCase
	UserService usecase.UserUseCase
}

func New(authService usecase.AuthUseCase, userService usecase.UserUseCase) *Application {
	return &Application{
		AuthService: authService,
		UserService: userService,
	}
}
