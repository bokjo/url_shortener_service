package json

import (
	"encoding/json"

	"github.com/bokjo/url_shortener_service/shortener"
	"github.com/pkg/errors"
)

// Serializer struct...
type Serializer struct{}

// Decode ...
func (s *Serializer) Decode(input []byte) (*shortener.Shortener, error) {
	shortener := &shortener.Shortener{}

	if err := json.Unmarshal(input, shortener); err != nil {
		return nil, errors.Wrap(err, "[serializer] Decode()")
	}

	return shortener, nil

}

// Encode ...
func (s *Serializer) Encode(input *shortener.Shortener) ([]byte, error) {
	msg, err := json.Marshal(input)

	if err != nil {
		return nil, errors.Wrap(err, "[serializer] Encode()")
	}

	return msg, nil
}
