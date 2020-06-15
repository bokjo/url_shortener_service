package shortener

// Shortener struct
type Shortener struct {
	Code      string `json: "code", bson: "code"`
	URL       string `json: "url", bson: "url"`
	CreatedAt int64  `json: "created_at", bson: "created_at"`
}
