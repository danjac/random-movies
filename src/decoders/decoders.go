package decoders

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
)

// MovieDecoder decodes and validates user input
type MovieDecoder struct {
	Title string `valid:"required"`
}

// Decode decodes user input
func (f *MovieDecoder) Decode(r *http.Request) error {
	return decode(r, f)
}

// decodes JSON body of request and runs through validator
func decode(r *http.Request, data interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(data); err != nil {
		return err
	}
	return nil
}
