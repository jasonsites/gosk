package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Envelope
type Envelope map[string]any

// RequestBody
type RequestBody struct {
	Data *RequestResource `json:"data" validate:"required"`
}

// RequestResource
type RequestResource struct {
	Type       string `json:"type" validate:"required"`
	ID         string `json:"id" validate:"omitempty,uuid4"`
	Attributes any    `json:"attributes" validate:"required"`
}

// JSONDecode
func (c *Controller) JSONDecode(w http.ResponseWriter, r *http.Request, dest any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(1048576))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dest); err != nil {
		return err
	}

	err := dec.Decode(&struct{}{})
	if err != io.EOF {
		// err = fmt.Errorf("request body must contain only one json object %w", err)
		return errors.New("request body must contain only one json object")
	}

	return nil
}

// JSONEncode
func (c *Controller) JSONEncode(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Add("Content-Type", "application/json") // TODO: .Set("Content-Type", ...) ?
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		c.Error(w, r, err)
	}
}
