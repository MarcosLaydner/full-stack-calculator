package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	assertResult(t, Add, 10, 4, 14)
}

func TestSubtract(t *testing.T) {
	assertResult(t, Subtract, 10, 4, 6)
}

func TestMultiply(t *testing.T) {
	assertResult(t, Multiply, 10, 4, 40)
}

func TestDivide(t *testing.T) {
	assertResult(t, Divide, 10, 4, 2.5)
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(1, 0)
	if !errors.Is(err, ErrDivisionByZero) {
		t.Errorf("Divide() error = %v, want %v", err, ErrDivisionByZero)
	}
}

func TestPower(t *testing.T) {
	assertResult(t, Power, 2, 4, 16)
}

func TestPowerOverflow(t *testing.T) {
	_, err := Power(10, 1000)
	if !errors.Is(err, ErrNonFiniteResult) {
		t.Errorf("Power() error = %v, want %v", err, ErrNonFiniteResult)
	}
}

func TestSquareRoot(t *testing.T) {
	got, err := SquareRoot(81)
	if err != nil {
		t.Fatalf("SquareRoot() error = %v", err)
	}
	if got != 9 {
		t.Errorf("SquareRoot() = %v, want 9", got)
	}
}

func TestSquareRootOfNegativeNumber(t *testing.T) {
	_, err := SquareRoot(-1)
	if !errors.Is(err, ErrNegativeSquareRoot) {
		t.Errorf("SquareRoot() error = %v, want %v", err, ErrNegativeSquareRoot)
	}
}

func TestPercentage(t *testing.T) {
	assertResult(t, Percentage, 20, 80, 16)
}

func TestNonFiniteInput(t *testing.T) {
	_, err := Add(math.Inf(1), 1)
	if !errors.Is(err, ErrNonFiniteResult) {
		t.Errorf("Add() error = %v, want %v", err, ErrNonFiniteResult)
	}
}

func assertResult(t *testing.T, operation func(float64, float64) (float64, error), a, b, want float64) {
	t.Helper()
	got, err := operation(a, b)
	if err != nil {
		t.Fatalf("operation error = %v", err)
	}
	if got != want {
		t.Errorf("operation result = %v, want %v", got, want)
	}
}
