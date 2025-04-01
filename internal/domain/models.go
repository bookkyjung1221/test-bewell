package domain

// InputOrder represents the raw order data from the platform
type InputOrder struct {
	No                int     `json:"no"`
	PlatformProductId string  `json:"platformProductId"`
	Qty               int     `json:"qty"`
	UnitPrice         float64 `json:"unitPrice"`
	TotalPrice        float64 `json:"totalPrice"`
}

// CleanedOrder represents the processed order data
type CleanedOrder struct {
	No         int     `json:"no"`
	ProductId  string  `json:"productId"`
	MaterialId string  `json:"materialId"`
	ModelId    string  `json:"modelId"`
	Qty        int     `json:"qty"`
	UnitPrice  float64 `json:"unitPrice"`
	TotalPrice float64 `json:"totalPrice"`
}

// Product represents the parsed product information
type Product struct {
	ProductId  string
	MaterialId string
	ModelId    string
	Qty        int
}
