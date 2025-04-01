// test/service/order_service_test.go
package service_test

import (
	"encoding/json"
	"reflect"

	"test_bewell/internal/domain"
	"test_bewell/internal/service"
	"testing"
)

func TestProcessOrders(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		input    []domain.InputOrder
		expected []domain.CleanedOrder
	}{
		{
			name: "Case 1: Only one product",
			input: []domain.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         50,
					TotalPrice:        100,
				},
			},
			expected: []domain.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "IPHONE16PROMAX",
					Qty:        2,
					UnitPrice:  50.00,
					TotalPrice: 100.00,
				},
				{
					No:         2,
					ProductId:  "WIPING-CLOTH",
					Qty:        2,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
				{
					No:         3,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
			},
		},
		{
			name: "Case 2: One product with wrong prefix",
			input: []domain.InputOrder{
				{
					No:                1,
					PlatformProductId: "x2-3&FG0A-CLEAR-IPHONE16PROMAX",
					Qty:               2,
					UnitPrice:         50,
					TotalPrice:        100,
				},
			},
			expected: []domain.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-IPHONE16PROMAX",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "IPHONE16PROMAX",
					Qty:        2,
					UnitPrice:  50.00,
					TotalPrice: 100.00,
				},
				{
					No:         2,
					ProductId:  "WIPING-CLOTH",
					Qty:        2,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
				{
					No:         3,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
			},
		},
		{
			name: "Case 3: One product with wrong prefix and has * symbol that indicates the quantity",
			input: []domain.InputOrder{
				{
					No:                1,
					PlatformProductId: "x2-3&FG0A-MATTE-IPHONE16PROMAX*3",
					Qty:               1,
					UnitPrice:         90,
					TotalPrice:        90,
				},
			},
			expected: []domain.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-MATTE-IPHONE16PROMAX",
					MaterialId: "FG0A-MATTE",
					ModelId:    "IPHONE16PROMAX",
					Qty:        3,
					UnitPrice:  30.00,
					TotalPrice: 90.00,
				},
				{
					No:         2,
					ProductId:  "WIPING-CLOTH",
					Qty:        3,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
				{
					No:         3,
					ProductId:  "MATTE-CLEANNER",
					Qty:        3,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
			},
		},
		{
			name: "Case 4: One bundle product with wrong prefix and split by / symbol into two product",
			input: []domain.InputOrder{
				{
					No:                1,
					PlatformProductId: "FG0A-CLEAR-OPPOA3/%20xFG0A-CLEAR-OPPOA3-B",
					Qty:               1,
					UnitPrice:         80,
					TotalPrice:        80,
				},
			},
			expected: []domain.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Qty:        1,
					UnitPrice:  40.00,
					TotalPrice: 40.00,
				},
				{
					No:         2,
					ProductId:  "FG0A-CLEAR-OPPOA3-B",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3-B",
					Qty:        1,
					UnitPrice:  40.00,
					TotalPrice: 40.00,
				},
				{
					No:         3,
					ProductId:  "WIPING-CLOTH",
					Qty:        2,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
				{
					No:         4,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
			},
		},
		{
			name: "Case 6: One bundle product with wrong prefix and have / symbol and * symbol",
			input: []domain.InputOrder{
				{
					No:                1,
					PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3",
					Qty:               1,
					UnitPrice:         120,
					TotalPrice:        120,
				},
			},
			expected: []domain.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Qty:        2,
					UnitPrice:  40.00,
					TotalPrice: 80.00,
				},
				{
					No:         2,
					ProductId:  "FG0A-MATTE-OPPOA3",
					MaterialId: "FG0A-MATTE",
					ModelId:    "OPPOA3",
					Qty:        1,
					UnitPrice:  40.00,
					TotalPrice: 40.00,
				},
				{
					No:         3,
					ProductId:  "WIPING-CLOTH",
					Qty:        3,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
				{
					No:         4,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
				{
					No:         5,
					ProductId:  "MATTE-CLEANNER",
					Qty:        1,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
			},
		},
		{
			name: "Case 7: one product and one bundle product with wrong prefix and have / symbol and * symbol",
			input: []domain.InputOrder{
				{
					No:                1,
					PlatformProductId: "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3*2",
					Qty:               1,
					UnitPrice:         160,
					TotalPrice:        160,
				},
				{
					No:                2,
					PlatformProductId: "FG0A-PRIVACY-IPHONE16PROMAX",
					Qty:               1,
					UnitPrice:         50,
					TotalPrice:        50,
				},
			},
			expected: []domain.CleanedOrder{
				{
					No:         1,
					ProductId:  "FG0A-CLEAR-OPPOA3",
					MaterialId: "FG0A-CLEAR",
					ModelId:    "OPPOA3",
					Qty:        2,
					UnitPrice:  40.00,
					TotalPrice: 80.00,
				},
				{
					No:         2,
					ProductId:  "FG0A-MATTE-OPPOA3",
					MaterialId: "FG0A-MATTE",
					ModelId:    "OPPOA3",
					Qty:        2,
					UnitPrice:  40.00,
					TotalPrice: 80.00,
				},
				{
					No:         3,
					ProductId:  "FG0A-PRIVACY-IPHONE16PROMAX",
					MaterialId: "FG0A-PRIVACY",
					ModelId:    "IPHONE16PROMAX",
					Qty:        1,
					UnitPrice:  50.00,
					TotalPrice: 50.00,
				},
				{
					No:         4,
					ProductId:  "WIPING-CLOTH",
					Qty:        5,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
				{
					No:         5,
					ProductId:  "CLEAR-CLEANNER",
					Qty:        2,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
				{
					No:         6,
					ProductId:  "MATTE-CLEANNER",
					Qty:        2,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
				{
					No:         7,
					ProductId:  "PRIVACY-CLEANNER",
					Qty:        1,
					UnitPrice:  0.00,
					TotalPrice: 0.00,
				},
			},
		},
	}

	// Create the service
	orderService := service.NewOrderService()

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Process the orders
			actual, err := orderService.ProcessOrders(tc.input)
			if err != nil {
				t.Fatalf("Error processing orders: %v", err)
			}

			// Debug output
			expectedJSON, _ := json.MarshalIndent(tc.expected, "", "  ")
			actualJSON, _ := json.MarshalIndent(actual, "", "  ")

			// Compare results
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", string(expectedJSON), string(actualJSON))
			}
		})
	}
}
