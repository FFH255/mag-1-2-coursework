package pb_cache

import (
	"context"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (c *Cache) Get(ctx context.Context, key Key) (models.WPM, bool) {
	wpm, err := c.client.GetFloat64(ctx, key.String())
	if err != nil {
		return 0, false
	}

	return models.WPM(wpm), true
}
