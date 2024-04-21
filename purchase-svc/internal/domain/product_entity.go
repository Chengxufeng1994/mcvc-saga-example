package domain

// ProductDetail value object
type ProductDetail struct {
	ProductName string
	Description string
	BrandName   string
	Inventory   int64
}

// Status enumeration
type Status int

const (
	// ProductOk is ok status
	ProductOk Status = iota
	// ProductNotFound is not found status
	ProductNotFound
)

// ProductStatus value object
type ProductStatus struct {
	ProductID uint64
	Price     int64
	Status    Status
}

func NewProductStatus(productId uint64, price int64, status Status) *ProductStatus {
	return &ProductStatus{
		ProductID: productId,
		Price:     price,
		Status:    status,
	}
}
