package model

// Product 充值产品
type Product struct {
	BaseModel
	Name     string
	Price    float64
	Discount float64
	Days     int
	Power    int
	Enabled  bool
	Sales    int
	SortNum  int
}
