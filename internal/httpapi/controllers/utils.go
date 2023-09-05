package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jasonsites/gosk-api/internal/types"
)

func (c *Controller) JSONResponse(w http.ResponseWriter, status int, data []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) ReadJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dest); err != nil {
		return err
	}

	err := dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("request body must contain only a single JSON object")
	}

	return nil
}

func (c *Controller) WriteJSON(w http.ResponseWriter, status int, data types.Map) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.JSONResponse(w, status, b)
}
