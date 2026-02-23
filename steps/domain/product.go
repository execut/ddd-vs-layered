package domain

type Product struct {
	Name string
}

func NewProduct(name string) (Product, error) {
	return Product{Name: name}, nil
}
