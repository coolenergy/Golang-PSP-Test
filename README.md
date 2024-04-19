<h1>GoLang Payment Processing API<h1>

# Overview

This is a simple payment processing API written in Go. It accepts payment details, validates them, simulates transaction processing with a mock acquirer, and stores the transaction status in an in-memory data storage.

# Design Choices

1. Gin Web Framework: Gin was chosen for its simplicity and performance. It’s one of the fastest web frameworks for Go and it’s very easy to set up and use for building APIs.

2. In-memory Data Storage: For simplicity, we use an in-memory map to store transactions. This is not suitable for production use, but it’s enough for this simple demo.

3. Luhn’s Algorithm for Card Validation: The card number is validated using Luhn’s algorithm, which is a simple checksum formula used to validate a variety of identification numbers.

4. Mock Acquirer Communication: A simple function simulates sending the transaction to an acquirer. The acquirer approves if the last digit of the card number is even and denies if it’s odd.

# Setup

## Prerequisites

1. Go version 1.14 or higher installed on your machine.

2. Make sure you have set up your GOPATH.

3. Install Gin package for Go using the following command:

   `go get github.com/gin-gonic/gin`

## Steps to Run the Application

1. Clone the repository or download the source code.
2. Navigate to the directory in which main.go is located.
3. Run the command go run main.go in your terminal.
4. The application will start and listen on localhost:8080.

# Usage

Send a POST request to http://localhost:8080/ with the following JSON payload:

```
{
"cardNumber": "4532015112830366",
"expiryDate": "0424",
"cvv": "123",
"amount": "100.50",
"currency": "USD",
"merchantId": "123456"
}
```

The server will respond with a JSON containing the transaction ID and the status:

```
{
"transactionId": "0",
"status": "Approved"
}
```

# Running Tests

To run tests, navigate to the directory in which main_test.go is located and run the following command in your terminal:

`go test -v`

# Note

This is a simple demo application and is not suitable for production use as it lacks proper error handling, data persistence, security measures like input sanitization, encryption, and many other features expected in a production-grade payment processing system.
