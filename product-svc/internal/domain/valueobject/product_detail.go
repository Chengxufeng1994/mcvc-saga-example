package valueobject

// ProductDetail value object
type ProductDetail struct {
	Name        string
	Description string
	BrandName   string
	Price       int64
}

func NewProductDetail(name string, description string, brandName string, price int64) *ProductDetail {
	return &ProductDetail{
		Name:        name,
		Description: description,
		BrandName:   brandName,
		Price:       price,
	}
}
