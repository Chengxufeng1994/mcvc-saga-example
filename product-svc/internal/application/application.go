package application

import "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"

type ProductApplication struct {
	ProductService usecase.ProductUseCase
}

func NewProductApplication(productService usecase.ProductUseCase) *ProductApplication {
	return &ProductApplication{
		ProductService: productService,
	}
}

type OrderApplication struct {
	OrderService usecase.OrderUseCase
}

func NewOrderApplication(orderService usecase.OrderUseCase) *OrderApplication {
	return &OrderApplication{
		OrderService: orderService,
	}
}

type PaymentApplication struct {
	PaymentService usecase.PaymentUseCase
}

func NewPaymentApplication(paymentService usecase.PaymentUseCase) *PaymentApplication {
	return &PaymentApplication{
		PaymentService: paymentService,
	}
}
