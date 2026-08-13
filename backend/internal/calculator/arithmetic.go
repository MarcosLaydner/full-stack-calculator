package calculator

import (
	"errors"
	"math"
)

var (
	ErrDivisionByZero     = errors.New("division by zero")
	ErrNegativeSquareRoot = errors.New("cannot calculate the square root of a negative number")
	ErrNonFiniteResult    = errors.New("calculation result is not finite")
)

func Add(a, b float64) (float64, error) {
	return finite(a + b)
}

func Subtract(a, b float64) (float64, error) {
	return finite(a - b)
}

func Multiply(a, b float64) (float64, error) {
	return finite(a * b)
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return finite(a / b)
}

func Power(a, b float64) (float64, error) {
	return finite(math.Pow(a, b))
}

func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSquareRoot
	}
	return finite(math.Sqrt(a))
}

func Percentage(a, b float64) (float64, error) {
	return finite(a * b / 100)
}

func finite(result float64) (float64, error) {
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, ErrNonFiniteResult
	}
	return result, nil
}
