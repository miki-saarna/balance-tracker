package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	plaid "github.com/plaid/plaid-go/v21/plaid"
)

func ConvertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	countryCodes := []plaid.CountryCode{}

	for _, countryCodeStr := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}

	return countryCodes
}

func ConvertProducts(productStrs []string) []plaid.Products {
	products := []plaid.Products{}

	for _, productStr := range productStrs {
		products = append(products, plaid.Products(productStr))
	}

	return products
}

func ContainsProduct(products []plaid.Products, product plaid.Products) bool {
	for _, p := range products {
		if p == product {
			return true
		}
	}
	return false
}

func RenderError(c *gin.Context, originalErr error) {
	if plaidError, err := plaid.ToPlaidError(originalErr); err == nil {
		// Return 200 and allow the front end to render the error.
		c.JSON(http.StatusOK, gin.H{"error": plaidError})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": originalErr.Error()})
}
