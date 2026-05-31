package profiles_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/shared/redis"
)

type Repository struct {
	client     *redis.Client
	expiration time.Duration
}

func (r *Repository) key(username string) string {
	return fmt.Sprintf("profile:%s", username)
}

func (r *Repository) Get(ctx context.Context, username string) (*models.Profile, bool) {
	data, err := r.client.Get(ctx, r.key(username))
	if err != nil {
		return nil, false
	}

	var profile models.Profile

	err = json.Unmarshal(data, &profile)
	if err != nil {
		return nil, false
	}

	return &profile, true
}

func (r *Repository) Save(ctx context.Context, profile *models.Profile) error {
	value, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	return r.client.Save(ctx, r.key(profile.Username), value, r.expiration)
}

func New(client *redis.Client, expiration time.Duration) *Repository {
	return &Repository{
		client:     client,
		expiration: expiration,
	}
}
