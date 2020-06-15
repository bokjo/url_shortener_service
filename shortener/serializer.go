package shortener

// ISerializer interface...
type ISerializer interface {
	Decode(input []byte) (*Shortener, error)
	Encode(input *Shortener) ([]byte, error)
}
