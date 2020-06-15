package shortener

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

var (
	// ErrorShortenerNotFound custom not found error
	ErrorShortenerNotFound = errors.New("Short url not found")
	// ErrorShortenerInvalidURL custom invalid error
	ErrorShortenerInvalidURL = errors.New("Invalid url, please provide valid url with protocol scheme")
)

type shortenerService struct {
	shortenerRepo ShortenerRepository
}

// NewShortenerService generates new ShortenerService instance
func NewShortenerService(shortenerRepo ShortenerRepository) ShortenerService {
	return &shortenerService{
		shortenerRepo,
	}
}

func (ss *shortenerService) Get(code string) (*Shortener, error) {
	return ss.shortenerRepo.Get(code)
}

func (ss *shortenerService) Store(shortener *Shortener) error {
	// TODO: better url parsing is needed (we must enforce the scheme (http/https) as part of the url)
	parsedURL, err := url.Parse(shortener.URL)
	fmt.Printf("%T, %#v\n", parsedURL.Path, parsedURL.Path)
	fmt.Println(err)
	if err != nil || parsedURL.Scheme == "" {
		return errs.Wrap(ErrorShortenerInvalidURL, "service.Shortener.Store")
	}

	// TODO: we must use some consistent hashing mechanism, so duplicate url will be stored more efficiently
	shortener.Code = shortid.MustGenerate()
	shortener.CreatedAt = time.Now().Unix()

	return ss.shortenerRepo.Store(shortener)

}
