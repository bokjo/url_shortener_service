package shortener

type ShortenerService interface {
	Get(code string) (*Shortener, error)
	Store(shortener *Shortener) error
}
