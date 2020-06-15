package redis

import (
	"fmt"
	"strconv"

	"github.com/bokjo/url_shortener_service/shortener"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "[repository:Redis] newRedisClient()->ParseURL")
	}

	client := redis.NewClient(options)

	_, err = client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "[repository:Redis] newRedisClient()->Ping")
	}

	return client, nil
}

// NewRedisRepository ...
func NewRedisRepository(redisURL string) (shortener.ShortenerRepository, error) {
	repository := &redisRepository{}

	client, err := newRedisClient(redisURL)

	if err != nil {
		return nil, errors.Wrap(err, "[repository:Redis] newRedisClient()->NewRedisRepository")
	}

	repository.client = client

	return repository, nil
}

func (rr *redisRepository) generateRedisKey(code string) string {
	return fmt.Sprintf("shortener:%s", code)
}

func (rr *redisRepository) Get(code string) (*shortener.Shortener, error) {
	short := &shortener.Shortener{}

	key := rr.generateRedisKey(code)

	data, err := rr.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "[repository:Redis] Get()")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrorShortenerNotFound, "[repository:Redis] Get()->No Data!")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "[repository:Redis] Get()->Cannot parse `created_at`")
	}

	short.Code = data["code"]
	short.URL = data["url"]
	short.CreatedAt = createdAt

	return short, nil
}

// Store ...
func (rr *redisRepository) Store(shortener *shortener.Shortener) error {
	key := rr.generateRedisKey(shortener.Code)

	data := map[string]interface{}{
		"code":       shortener.Code,
		"created_at": shortener.CreatedAt,
		"url":        shortener.URL,
	}

	_, err := rr.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "[repository:Redis] Store()")
	}

	return nil
}
