package api

import (
	"io/ioutil"
	"log"
	"net/http"

	js "github.com/bokjo/url_shortener_service/serializers/json"
	"github.com/bokjo/url_shortener_service/shortener"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// API interface...
type API interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type apiHandler struct {
	shortenerService shortener.ShortenerService
}

// NewAPIHandler ...
func NewAPIHandler(shortenerService shortener.ShortenerService) API {
	return &apiHandler{shortenerService: shortenerService}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (ah *apiHandler) serializer(contentType string) shortener.ISerializer {
	return &js.Serializer{}
}

func (ah *apiHandler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	short, err := ah.shortenerService.Get(code)
	if err != nil {
		if errors.Cause(err) == shortener.ErrorShortenerNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, short.URL, http.StatusMovedPermanently)
}

func (ah *apiHandler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	short, err := ah.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ah.shortenerService.Store(short)
	if err != nil {
		if errors.Cause(err) == shortener.ErrorShortenerInvalidURL {
			http.Error(w, shortener.ErrorShortenerInvalidURL.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseBody, err := ah.serializer(contentType).Encode(short)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(w, contentType, responseBody, http.StatusCreated)
}
