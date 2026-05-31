package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const Nil = redis.Nil

type Config struct {
	Address  string
	Password string
}

type Client struct {
	config *Config
	client *redis.Client
}

func (c *Client) MustConnect(ctx context.Context) {
	pingTimeout := 5 * time.Second

	client := redis.NewClient(&redis.Options{
		Addr:     c.config.Address,
		Password: c.config.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("redis connection test failed: %s", err.Error()))
	}

	c.client = client
}

func (c *Client) MustClose() {
	err := c.client.Close()
	if err != nil {
		panic(fmt.Sprintf("redis connection test failed: %s", err.Error()))
	}
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

func (c *Client) GetFloat64(ctx context.Context, key string) (float64, error) {
	return c.client.Get(ctx, key).Float64()
}

func (c *Client) Save(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) UnderlyingClient() *redis.Client {
	return c.client
}

func New(config *Config) *Client {
	return &Client{
		config: config,
	}
}
