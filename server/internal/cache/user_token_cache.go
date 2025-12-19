package cache

import (
	"errors"
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

var UserTokenCache = newUserTokenCache()

type userTokenCache struct {
	cache cache.LoadingCache
}

func newUserTokenCache() *userTokenCache {
	return &userTokenCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				token := repositories.UserTokenRepository.GetByToken(sqls.DB(), key.(string))
				if token == nil {
					e = errors.New("数据不存在")
				} else {
					value = token
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(60*time.Minute),
		),
	}
}

func (c *userTokenCache) Get(token string) *models.UserToken {
	if len(token) == 0 {
		return nil
	}
	val, err := c.cache.Get(token)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.(*models.UserToken)
	}
	return nil
}

func (c *userTokenCache) Invalidate(token string) {
	c.cache.Invalidate(token)
}
