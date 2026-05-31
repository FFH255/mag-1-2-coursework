package pb_cache

import (
	"context"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (c *Cache) Save(ctx context.Context, key Key, wpm models.WPM) error {
	return c.client.Save(ctx, key.String(), float64(wpm), c.expiration)
}
