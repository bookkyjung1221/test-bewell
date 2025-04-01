package main

import (
	"encoding/json"
	"fmt"
	"log"

	"test_bewell/internal/domain"
	"test_bewell/internal/service"
)

func main() {
	// Initialize the order service
	orderService := service.NewOrderService()

	// Sample input (for demonstration)
	inputData := []domain.InputOrder{
		{
			No:                1,
			PlatformProductId: "FG0A-CLEAR-IPHONE16PROMAX",
			Qty:               2,
			UnitPrice:         50,
			TotalPrice:        100,
		},
	}

	// Process the orders
	cleanedOrders, err := orderService.ProcessOrders(inputData)
	if err != nil {
		log.Fatalf("Error processing orders: %v", err)
	}

	// Output the result
	output, _ := json.MarshalIndent(cleanedOrders, "", "  ")
	fmt.Println(string(output))
}
