package decoders

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"net/http"
)

type MovieDecoder struct {
	Title string `valid:"required"`
}

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
