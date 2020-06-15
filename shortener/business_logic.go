package shortener

import (
	"errors"
	"net/url"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

var (
	// ErrorShortenerNotFound custom not found error
	ErrorShortenerNotFound = errors.New("Short url not found")
	// ErrorShortenerInvalidURL custom invalid error
	ErrorShortenerInvalidURL = errors.New("Invalid short url")
)

type shortenerService struct {
	shortenerRepo *ShortenerRepository
}

// NewShortenerService generates new ShortenerService instance
func NewShortenerService(shortenerRepo *ShortenerRepository) ShortenerService {
	return &shortenerService{
		shortenerRepo: shortenerRepo,
	}
}

func (ss *shortenerService) Get(code string) (*Shortener, error) {
	return ss.shortenerRepo.Get(code)
}

func (ss *shortenerService) Store(shortener *Shortener) error {
	if _, err := url.Parse(shortener.URL); err != nil {
		return errs.Wrap(ErrorShortenerInvalidURL, "service.Shortener.Store")
	}
	shortener.Code = shortid.MustGenerate()
	shortener.CreatedAt = time.Now().Unix()

	return ss.shortenerRepo.Store(shortener)

}
