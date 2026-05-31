package pb_cache

import (
	"fmt"
	"time"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/shared/redis"
)

// keyTemplate - pb:user_id:language:mode:submode
// example: pb:12345:english:words:10
const keyTemplate = "pb:%d:%s:%s:%s"

type Key struct {
	UserID   models.ID
	Language models.Language
	Mode     models.Mode
	Submode  models.Submode
}

func (k Key) String() string {
	return fmt.Sprintf(keyTemplate, k.UserID, k.Language, k.Mode, k.Submode)
}

type Cache struct {
	client     *redis.Client
	expiration time.Duration
}

func New(client *redis.Client, expiration time.Duration) *Cache {
	return &Cache{
		client:     client,
		expiration: expiration,
	}
}
