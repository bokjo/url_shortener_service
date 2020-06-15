package shortener

type Serializer interface {
	Encode(input []byte) (*Shortener, error)
	Decode(input *Shortener) ([]byte, error)
}
