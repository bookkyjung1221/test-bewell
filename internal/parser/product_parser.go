package parser

import (
	"regexp"
	"strconv"
	"strings"

	"test_bewell/internal/domain"
)

// ProductParser handles the extraction of product details from various formats
type ProductParser struct {
	// Regular expression to match the internal product code pattern
	productRegex *regexp.Regexp
}

// NewProductParser creates a new product parser
func NewProductParser() *ProductParser {
	// This regex captures the entire pattern of film type ID, texture ID, and phone model ID
	// The pattern is: FG0A-CLEAR-IPHONE16PROMAX or similar patterns
	regex := regexp.MustCompile(`(FG0[A-Z0-9]{1,2})-(CLEAR|MATTE|PRIVACY)-([A-Z0-9-]+)`)

	return &ProductParser{
		productRegex: regex,
	}
}

// Parse extracts product information from a platform-specific product ID
func (p *ProductParser) Parse(platformProductId string) ([]domain.Product, error) {
	// Split by "/" to handle bundled products
	bundledProducts := strings.Split(platformProductId, "/")
	products := make([]domain.Product, 0, len(bundledProducts))

	for _, productStr := range bundledProducts {
		// Clean up the product string by removing any prefixes or special characters
		cleanProductStr := p.cleanProductString(productStr)

		// Check for quantity indicator "*"
		quantity := 1
		if strings.Contains(cleanProductStr, "*") {
			parts := strings.Split(cleanProductStr, "*")
			if len(parts) == 2 {
				cleanProductStr = parts[0]
				qty, err := strconv.Atoi(parts[1])
				if err == nil && qty > 0 {
					quantity = qty
				}
			}
		}

		// Extract product details using regex
		matches := p.productRegex.FindStringSubmatch(cleanProductStr)
		if len(matches) != 4 {
			// If we can't match the pattern, try to look for the pattern within the string
			fullStr := p.productRegex.FindString(cleanProductStr)
			if fullStr != "" {
				matches = p.productRegex.FindStringSubmatch(fullStr)
			}

			if len(matches) != 4 {
				return nil, domain.ErrInvalidProductFormat
			}
		}

		// Create product with extracted details
		product := domain.Product{
			ProductId:  matches[0],
			MaterialId: matches[1] + "-" + matches[2],
			ModelId:    matches[3],
			Qty:        quantity,
		}

		products = append(products, product)
	}

	return products, nil
}

// cleanProductString removes any non-essential characters from the product string
func (p *ProductParser) cleanProductString(productStr string) string {
	// Trim spaces
	productStr = strings.TrimSpace(productStr)

	// Remove URL encoding like %20
	productStr = strings.ReplaceAll(productStr, "%20", "")

	// Find the product pattern within the string
	match := p.productRegex.FindString(productStr)
	if match != "" {
		// Extract quantity if present
		if idx := strings.Index(productStr, "*"); idx > 0 && strings.HasPrefix(productStr[idx-len(match):idx], match) {
			return match + productStr[idx:]
		}
		return match
	}

	return productStr
}
