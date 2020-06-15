package shortener

// ShortenerRepository ...
type ShortenerRepository interface {
	Get(code string) (*Shortener, error)
	Store(shortener *Shortener) error
}
