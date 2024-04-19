package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Transaction represents a payment transaction
type Transaction struct {
	CardNumber string `json:"cardNumber"`
	ExpiryDate string `json:"expiryDate"`
	CVV        string `json:"cvv"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	MerchantID string `json:"merchantId"`
	Status     string `json:"status"`
}

// Transactions is the in-memory store for transactions
var Transactions = make(map[string]Transaction)

// TransactionNumber is incremented each time a new transaction is added
var TransactionIdCounter = 0

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	// Endpoint to process a payment transaction
	r.POST("/", processPayment)

	r.Run()
}

// processPayment handles a payment request
func processPayment(c *gin.Context) {
	var transaction Transaction

	// Parse and validate the incoming transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the card details using Luhn's algorithm
	isValid, validationErrors := validateCardDetails(transaction)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card details", "details": validationErrors})
		return
	}

	// Simulate acquirer communication and update transaction status
	status := simulateAcquirer(transaction.CardNumber)
	transaction.Status = status

	// Store the transaction and increment the transaction number
	TransactionId := strconv.Itoa(TransactionIdCounter)
	addTransaction(TransactionId, transaction)
	TransactionIdCounter++

	// Respond to the merchant with the transaction status
	c.JSON(http.StatusOK, gin.H{
		"transactionId": TransactionId,
		"status":        status,
	})
}

// addTransaction adds a transaction to the in-memory store
func addTransaction(transactionID string, transaction Transaction) {
	Transactions[transactionID] = transaction
}

// validateCardDetails checks if the card number is valid using Luhn's algorithm
func validateCardDetails(details Transaction) (bool, map[string]string) {
	isValid := true
	var validationErrors = make(map[string]string)

	if !checkLuhn(details.CardNumber) {
		validationErrors["cardNumber"] = "Invalid card number"
		isValid = false
	}

	if len(details.ExpiryDate) != 4 {
		validationErrors["expiryDate"] = "Invalid expire date"
		isValid = false
	}

	// Check if CVV is valid
	if len(details.CVV) != 3 {
		validationErrors["cvv"] = "Invalid cvv"
		isValid = false
	}

	// Check if the amount is a valid float
	if _, err := strconv.ParseFloat(details.Amount, 64); err != nil {
		validationErrors["amount"] = "Invalid amount"
		isValid = false
	}

	// Check if currency is valid
	if details.Currency == "" {
		validationErrors["currency"] = "Invalid currency"
		isValid = false
	}

	// Check if merchant ID is valid
	if details.MerchantID == "" {
		validationErrors["merchantId"] = "Invalid merchant Id"
		isValid = false
	}

	return isValid, validationErrors
}

// checkLuhn implements Luhn's algorithm to validate a card number
func checkLuhn(cardNumber string) bool {
	var sum int
	numDigits := len(cardNumber)
	oddeven := numDigits & 1

	for count := 0; count < numDigits; count++ {
		digit, _ := strconv.Atoi(string(cardNumber[count]))

		if ((count & 1) ^ oddeven) == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return (sum % 10) == 0
}

// simulateAcquirer simulates the role of an acquirer which approves or denies a transaction
func simulateAcquirer(cardNumber string) string {
	lastDigit := int(cardNumber[len(cardNumber)-1] - '0')
	if lastDigit%2 == 0 {
		return "Approved"
	} else {
		return "Denied"
	}
}
