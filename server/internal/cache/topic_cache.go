package cache

import (
	"bbs-go/internal/models/constants"
	"errors"
	"log/slog"
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

var (
	topicRecommendCacheKey = "recommend_topics_cache"
)

var TopicCache = newTopicCache()

type topicCache struct {
	recommendCache       cache.LoadingCache
	needAHandTopicsCache cache.LoadingCache
}

func newTopicCache() *topicCache {
	return &topicCache{
		recommendCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				topics := repositories.TopicRepository.Find(sqls.DB(),
					sqls.NewCnd().Eq("status", constants.StatusOk).Desc("id").Limit(50))
				if topics == nil {
					e = errors.New("数据不存在")
				} else {
					value = topics
				}
				return
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
		needAHandTopicsCache: cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			slog.Info("开始加载呼叫支援话题缓存", slog.String("cache_key", key.(string)))
			cnd := sqls.NewCnd().Eq("status", constants.StatusOk).Eq("need_a_hand", 1).Desc("last_comment_time").Limit(50)
			slog.Info("查询条件", slog.Any("conditions", map[string]interface{}{
				"status": constants.StatusOk,
				"need_a_hand": 1,
				"order": "last_comment_time DESC",
				"limit": 50,
			}))
			topics := repositories.TopicRepository.Find(sqls.DB(), cnd)
			slog.Info("查询结果", slog.Int("count", len(topics)))
			// 确保总是返回有效的数组，避免空指针问题
			if topics == nil {
				topics = []models.Topic{}
			}
			value = topics
			return
		},
		cache.WithMaximumSize(10),
		cache.WithRefreshAfterWrite(10*time.Minute),
	),
	}
}

func (c *topicCache) GetRecommendTopics() []models.Topic {
	val, err := c.recommendCache.Get(topicRecommendCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]models.Topic)
	}
	return nil
}

func (c *topicCache) GetNeedAHandTopics() []models.Topic {
	val, err := c.needAHandTopicsCache.Get("need_a_hand_topics")
	if err != nil {
		slog.Error("获取呼叫支援话题缓存失败", slog.Any("err", err))
		return []models.Topic{}
	}
	if val != nil {
		topics := val.([]models.Topic)
		slog.Info("获取呼叫支援话题缓存成功", slog.Int("count", len(topics)))
		return topics
	}
	slog.Info("呼叫支援话题缓存为空")
	return []models.Topic{}
}

func (c *topicCache) InvalidateRecommend() {
	c.recommendCache.Invalidate(topicRecommendCacheKey)
}

func (c *topicCache) InvalidateNeedAHand() {
	c.needAHandTopicsCache.Invalidate("need_a_hand_topics")
}