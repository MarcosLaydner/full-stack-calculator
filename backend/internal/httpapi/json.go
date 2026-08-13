package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type binaryRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

type unaryRequest struct {
	A *float64 `json:"a"`
}

type resultResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type binaryOperation func(float64, float64) (float64, error)
type unaryOperation func(float64) (float64, error)

func handleBinary(w http.ResponseWriter, r *http.Request, operation binaryOperation) {
	var request binaryRequest
	if err := decodeRequest(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if request.A == nil {
		writeError(w, http.StatusBadRequest, "a is required")
		return
	}
	if request.B == nil {
		writeError(w, http.StatusBadRequest, "b is required")
		return
	}
	result, err := operation(*request.A, *request.B)
	writeOperationResult(w, result, err)
}

func handleUnary(w http.ResponseWriter, r *http.Request, operation unaryOperation) {
	var request unaryRequest
	if err := decodeRequest(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if request.A == nil {
		writeError(w, http.StatusBadRequest, "a is required")
		return
	}
	result, err := operation(*request.A)
	writeOperationResult(w, result, err)
}

func decodeRequest(w http.ResponseWriter, r *http.Request, destination any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return describeDecodeError(err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err != nil {
			return describeDecodeError(err)
		}
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func describeDecodeError(err error) error {
	var maxBytesError *http.MaxBytesError
	var syntaxError *json.SyntaxError
	var typeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &maxBytesError):
		return errors.New("request body must not exceed 1 MB")
	case errors.Is(err, io.EOF):
		return errors.New("request body is required")
	case errors.As(err, &syntaxError), errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("request body contains malformed JSON")
	case errors.As(err, &typeError):
		if typeError.Field != "" {
			return fmt.Errorf("%s must be a number", typeError.Field)
		}
		return errors.New("request body contains an invalid value type")
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		return fmt.Errorf("request body contains %s", strings.TrimPrefix(err.Error(), "json: "))
	default:
		return errors.New("request body must be valid JSON")
	}
}

func writeOperationResult(w http.ResponseWriter, result float64, err error) {
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err := writeJSON(w, http.StatusOK, resultResponse{Result: result}); err != nil {
		log.Printf("write operation response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	if err := writeJSON(w, status, errorResponse{Error: message}); err != nil {
		log.Printf("write error response: %v", err)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode JSON response: %w", err)
	}
	payload = append(payload, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("write JSON response: %w", err)
	}
	return nil
}
