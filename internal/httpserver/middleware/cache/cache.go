package cache

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	"one-day-server/configs"
	"one-day-server/internal/db/mysql"
	"one-day-server/internal/httpserver/middleware/ratelimit"
	"one-day-server/utils"
)

var (
	Instance                = newManager()
	DefaultCacheExpiredTime = time.Duration(configs.GetEnvDefaultInt("DEFAULT_CACHE_EXPIRED_TIME_IN_MILLISECOND", 1000)) * time.Millisecond
	DefaultCacheCleanTime   = time.Duration(configs.GetEnvDefaultInt("DEFAULT_CACHE_CLEAN_TIME_IN_MINUTE", 1)) * time.Minute
	DefaultRateLimitWeight  = 1
)

type manager struct {
	cache    *cache.Cache
	configMu sync.RWMutex
	config   map[string]*APIConfig
}

type APIConfig struct {
	Path     string        `gorm:"column:path"`
	Duration time.Duration `gorm:"column:duration"`
	Method   string        `gorm:"column:method"`
	Weight   int           `gorm:"column:weight"`
	mu       sync.Mutex
}

func (p *APIConfig) TableName() string {
	return "gateway_api_config"
}

func (p *APIConfig) getAPIConfigKey() string {
	return generateAPIConfigKey(p.Path, p.Method)
}

func generateAPIConfigKey(path, method string) string {
	return path + method
}

func newManager() *manager {
	m := &manager{
		cache:  cache.New(DefaultCacheExpiredTime, DefaultCacheCleanTime),
		config: map[string]*APIConfig{},
	}
	return m
}

func (m *manager) UseCache(c *gin.Context) {
	pathConfig := m.getOrCreatePathConfig(c.Request.URL.Path, c.Request.Method)
	ratelimit.ValidateRateLimit(c, pathConfig.Weight)
	if c.IsAborted() {
		return
	}
	if c.Request.Method != http.MethodGet {
		c.Next()
		return
	}
	key := parseCacheKey(c)
	if cachedValue, found := m.cache.Get(key); found {
		log.Infof("use cache: %s", key)
		cachedValue.(*cachedResponseWriter).FlushTo(c.Writer)
		c.Abort()
		return
	}
	pathConfig.mu.Lock()
	defer pathConfig.mu.Unlock()
	// check cache again in lock to avoid query twice
	if cachedValue, found := m.cache.Get(key); found {
		log.Infof("use cache: %s", key)
		cachedValue.(*cachedResponseWriter).FlushTo(c.Writer)
		c.Abort()
		return
	}
	cachedWriter := newCachedWriter(c.Writer)
	c.Writer = cachedWriter
	c.Next()
	cachedWriter.DoCache()
	m.cache.Set(key, cachedWriter, pathConfig.Duration)
}

func (m *manager) getOrCreatePathConfig(path string, method string) *APIConfig {
	key := generateAPIConfigKey(path, method)
	m.configMu.RLock()
	pathConfig, ok := m.config[key]
	m.configMu.RUnlock()
	if !ok {
		m.configMu.Lock()
		pathConfig, ok = m.config[key]
		if !ok {
			pathConfig = &APIConfig{
				Path:     path,
				Method:   method,
				Duration: DefaultCacheExpiredTime,
				Weight:   DefaultRateLimitWeight,
			}
			m.config[key] = pathConfig
		}
		m.configMu.Unlock()
	}
	return pathConfig
}

func (m *manager) loadConfig() {
	var apiConfigs []APIConfig
	if err := mysql.DB().Model(APIConfig{}).Find(&apiConfigs).Error; err != nil {
		log.Errorf("load config failed, err: %s", err)
		return
	}
	m.configMu.Lock()
	for i := range apiConfigs {
		apiConfig := &apiConfigs[i]
		m.config[apiConfig.getAPIConfigKey()] = &APIConfig{
			Path:     apiConfig.Path,
			Duration: apiConfig.Duration,
			Method:   apiConfig.Method,
			Weight:   apiConfig.Weight,
			mu:       sync.Mutex{},
		}
	}
	m.configMu.Unlock()
}

func parseCacheKey(c *gin.Context) string {
	path := c.Request.URL.Path
	if strings.HasPrefix(path, utils.UserApiPrefix) {
		return fmt.Sprintf("%s##%s##%s", c.Request.URL.Path, utils.SortQueryString(c.Request.URL.RawQuery), c.GetHeader(utils.OneDayApiKey))
	}
	return fmt.Sprintf("%s##%s", c.Request.URL.Path, utils.SortQueryString(c.Request.URL.RawQuery))
}
