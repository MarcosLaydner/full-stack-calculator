package httpapi

import (
	"net/http"

	"full-stack-calculator/backend/internal/calculator"
)

func add(w http.ResponseWriter, r *http.Request) {
	handleBinary(w, r, calculator.Add)
}

func subtract(w http.ResponseWriter, r *http.Request) {
	handleBinary(w, r, calculator.Subtract)
}

func multiply(w http.ResponseWriter, r *http.Request) {
	handleBinary(w, r, calculator.Multiply)
}

func divide(w http.ResponseWriter, r *http.Request) {
	handleBinary(w, r, calculator.Divide)
}

func power(w http.ResponseWriter, r *http.Request) {
	handleBinary(w, r, calculator.Power)
}

func squareRoot(w http.ResponseWriter, r *http.Request) {
	handleUnary(w, r, calculator.SquareRoot)
}

func percentage(w http.ResponseWriter, r *http.Request) {
	handleBinary(w, r, calculator.Percentage)
}
