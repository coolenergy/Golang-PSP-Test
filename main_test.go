package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test for validateCardDetails function
func TestValidateCardDetails(t *testing.T) {
	transaction := Transaction{
		CardNumber: "4532015112830366",
		ExpiryDate: "0424",
		CVV:        "123",
		Amount:     "100.50",
		Currency:   "USD",
		MerchantID: "123456",
	}

	isValid, _ := validateCardDetails(transaction)
	assert.True(t, isValid)

	// Testing with invalid card number
	transaction.CardNumber = "1234"
	isValid, _ = validateCardDetails(transaction)
	assert.False(t, isValid)
}

// Test for checkLuhn function
func TestCheckLuhn(t *testing.T) {
	// Test with valid card number
	assert.True(t, checkLuhn("4532015112830366"))

	// Test with invalid card number
	assert.False(t, checkLuhn("1234"))
}

// Test for simulateAcquirer function
func TestSimulateAcquirer(t *testing.T) {
	// Test with a card number ending in an even digit
	assert.Equal(t, "Approved", simulateAcquirer("4532015112830366"))

	// Test with a card number ending in an odd digit
	assert.Equal(t, "Denied", simulateAcquirer("4532015112830365"))
}

// Test for processPayment endpoint
func TestProcessPayment(t *testing.T) {
	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// Setup your router, just like you did in your main function
	r := gin.Default()

	// Define the route just like you did in main
	r.POST("/", processPayment)

	// Now, create a request to send to above route
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{
		"cardNumber": "4532015112830366",
		"expiryDate": "0424",
		"cvv": "123",
		"amount": "100.50",
		"currency": "USD",
		"merchantId": "123456"
	}`))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	resp := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(resp, req)

	// Assert we got the expected response code
	assert.Equal(t, 200, resp.Code)

	// Now you can check if the response body is what you expected
	expectedBody := `{"status":"Approved","transactionId":"0"}`
	assert.JSONEq(t, expectedBody, resp.Body.String())
}
