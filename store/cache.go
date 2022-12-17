package store

import (
	"context"
	"time"

	gcache "github.com/Code-Hex/go-generics-cache"
	"github.com/rs/zerolog"
	"github.com/stnokott/r6api"
)

type cache struct {
	api       *r6api.R6API
	userCache *gcache.Cache[string, *r6api.Profile] // resolving usernames to profile IDs
	logger    *zerolog.Logger
}

var cacheUserExpiration = gcache.WithExpiration(1 * time.Hour)

func newCache(api *r6api.R6API, logger *zerolog.Logger, ctx context.Context) *cache {
	return &cache{
		api:       api,
		userCache: gcache.NewContext(ctx, gcache.WithJanitorInterval[string, *r6api.Profile](1*time.Hour)),
		logger:    logger,
	}
}

func (c *cache) GetProfile(username string) (p *r6api.Profile, err error) {
	if c.userCache.Contains(username) {
		c.logger.Debug().Str("username", username).Msg("using cached profile")
		p, _ = c.userCache.Get(username)
	} else {
		p, err = c.api.ResolveUser(username)
		if p != nil {
			c.userCache.Set(username, p, cacheUserExpiration)
		}
	}
	return
}
