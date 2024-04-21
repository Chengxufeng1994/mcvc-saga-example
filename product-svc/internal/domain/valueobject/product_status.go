package valueobject

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

func NewProductStatus(productID uint64, price int64, existed bool) *ProductStatus {
	status := ProductOk
	if !existed {
		status = ProductNotFound
	}

	return &ProductStatus{
		ProductID: productID,
		Price:     price,
		Status:    status,
	}
}
