package service

import (
	"math"
	"strings"

	"test_bewell/internal/domain"
	"test_bewell/internal/parser"
)

// OrderService handles the business logic for processing orders
type OrderService struct {
	parser *parser.ProductParser
}

// NewOrderService creates a new order service
func NewOrderService() *OrderService {
	return &OrderService{
		parser: parser.NewProductParser(),
	}
}

// ProcessOrders processes a batch of input orders and returns cleaned orders
func (s *OrderService) ProcessOrders(inputOrders []domain.InputOrder) ([]domain.CleanedOrder, error) {
	cleanedOrders := make([]domain.CleanedOrder, 0)

	// Process each input order
	orderNo := 1
	for _, input := range inputOrders {
		products, err := s.parser.Parse(input.PlatformProductId)
		if err != nil {
			return nil, err
		}

		// Calculate the unit price per product
		pricePerProduct := input.TotalPrice / float64(len(products))
		totalQty := 0

		// Add each product to the cleaned orders
		for _, product := range products {
			unitPrice := pricePerProduct / float64(product.Qty)
			// Round to 2 decimal places
			unitPrice = math.Round(unitPrice*100) / 100

			cleanedOrder := domain.CleanedOrder{
				No:         orderNo,
				ProductId:  product.ProductId,
				MaterialId: product.MaterialId,
				ModelId:    product.ModelId,
				Qty:        product.Qty,
				UnitPrice:  unitPrice,
				TotalPrice: unitPrice * float64(product.Qty),
			}

			cleanedOrders = append(cleanedOrders, cleanedOrder)
			orderNo++
			totalQty += product.Qty
		}

		// Track the textures used to add appropriate cleaners
		textureQty := make(map[string]int)
		for _, product := range products {
			textureParts := strings.Split(product.MaterialId, "-")
			if len(textureParts) == 2 {
				texture := textureParts[1]
				textureQty[texture] += product.Qty
			}
		}

		// Add complementary items

		// Add wiping cloth per ordered quantity
		cleanedOrders = append(cleanedOrders, domain.CleanedOrder{
			No:         orderNo,
			ProductId:  "WIPING-CLOTH",
			Qty:        totalQty,
			UnitPrice:  0.00,
			TotalPrice: 0.00,
		})
		orderNo++

		// Add cleaners based on texture
		for texture, qty := range textureQty {
			cleanedOrders = append(cleanedOrders, domain.CleanedOrder{
				No:         orderNo,
				ProductId:  texture + "-CLEANNER",
				Qty:        qty,
				UnitPrice:  0.00,
				TotalPrice: 0.00,
			})
			orderNo++
		}
	}

	return cleanedOrders, nil
}
