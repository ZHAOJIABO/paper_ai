package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/repository"
	"paper_ai/pkg/logger"

	"go.uber.org/zap"
)

// PromptService Prompt服务
type PromptService struct {
	promptRepo repository.PolishPromptRepository
	cache      *promptCache
}

// NewPromptService 创建Prompt服务
func NewPromptService(promptRepo repository.PolishPromptRepository) *PromptService {
	return &PromptService{
		promptRepo: promptRepo,
		cache:      newPromptCache(),
	}
}

// GetPrompt 获取Prompt（带缓存）
func (s *PromptService) GetPrompt(ctx context.Context, versionType, language, style string) (*entity.PolishPrompt, error) {
	// 先从缓存获取
	cacheKey := buildPromptCacheKey(versionType, language, style)
	if cached := s.cache.get(cacheKey); cached != nil {
		logger.Debug("prompt cache hit",
			zap.String("version_type", versionType),
			zap.String("language", language),
			zap.String("style", style))
		return cached, nil
	}

	// 缓存未命中，从数据库查询
	prompt, err := s.promptRepo.GetActive(ctx, versionType, language, style)
	if err != nil {
		logger.Error("failed to get prompt from database",
			zap.String("version_type", versionType),
			zap.String("language", language),
			zap.String("style", style),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	// 存入缓存
	s.cache.set(cacheKey, prompt)
	logger.Debug("prompt cached",
		zap.String("version_type", versionType),
		zap.String("language", language),
		zap.String("style", style),
		zap.Int64("prompt_id", prompt.ID))

	return prompt, nil
}

// RenderPrompt 渲染Prompt模板（替换变量）
func (s *PromptService) RenderPrompt(ctx context.Context, versionType, language, style, content string) (*RenderedPrompt, error) {
	// 获取Prompt模板
	prompt, err := s.GetPrompt(ctx, versionType, language, style)
	if err != nil {
		return nil, err
	}

	// 准备变量
	variables := map[string]string{
		"content":  content,
		"language": language,
		"style":    style,
	}

	// 渲染用户提示词
	userPrompt := prompt.RenderUserPrompt(variables)

	return &RenderedPrompt{
		PromptID:     prompt.ID,
		SystemPrompt: prompt.SystemPrompt,
		UserPrompt:   userPrompt,
	}, nil
}

// IncrementUsage 增加Prompt使用次数
func (s *PromptService) IncrementUsage(ctx context.Context, promptID int64) error {
	return s.promptRepo.IncrementUsage(ctx, promptID)
}

// InvalidateCache 清除缓存
func (s *PromptService) InvalidateCache(versionType, language, style string) {
	cacheKey := buildPromptCacheKey(versionType, language, style)
	s.cache.delete(cacheKey)
	logger.Info("prompt cache invalidated",
		zap.String("version_type", versionType),
		zap.String("language", language),
		zap.String("style", style))
}

// InvalidateAllCache 清除所有缓存
func (s *PromptService) InvalidateAllCache() {
	s.cache.clear()
	logger.Info("all prompt cache invalidated")
}

// RenderedPrompt 渲染后的Prompt
type RenderedPrompt struct {
	PromptID     int64
	SystemPrompt string
	UserPrompt   string
}

// buildPromptCacheKey 构建缓存Key
func buildPromptCacheKey(versionType, language, style string) string {
	return fmt.Sprintf("%s:%s:%s", versionType, language, style)
}

// promptCache Prompt缓存（LRU + TTL）
type promptCache struct {
	mu      sync.RWMutex
	data    map[string]*cacheEntry
	maxSize int
	ttl     time.Duration
}

// cacheEntry 缓存条目
type cacheEntry struct {
	prompt    *entity.PolishPrompt
	expiresAt time.Time
}

// newPromptCache 创建缓存
func newPromptCache() *promptCache {
	cache := &promptCache{
		data:    make(map[string]*cacheEntry),
		maxSize: 100,         // 最大缓存100个Prompt
		ttl:     30 * time.Minute, // 30分钟TTL
	}

	// 启动后台清理goroutine
	go cache.cleanupLoop()

	return cache
}

// get 获取缓存
func (c *promptCache) get(key string) *entity.PolishPrompt {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil
	}

	// 检查是否过期
	if time.Now().After(entry.expiresAt) {
		return nil
	}

	return entry.prompt
}

// set 设置缓存
func (c *promptCache) set(key string, prompt *entity.PolishPrompt) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果缓存已满，删除最旧的条目（简单LRU）
	if len(c.data) >= c.maxSize {
		// 找到最早过期的条目并删除
		oldestKey := ""
		oldestTime := time.Now().Add(24 * time.Hour)
		for k, v := range c.data {
			if v.expiresAt.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.expiresAt
			}
		}
		if oldestKey != "" {
			delete(c.data, oldestKey)
		}
	}

	c.data[key] = &cacheEntry{
		prompt:    prompt,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// delete 删除缓存
func (c *promptCache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// clear 清空缓存
func (c *promptCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*cacheEntry)
}

// cleanupLoop 定期清理过期缓存
func (c *promptCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup 清理过期缓存
func (c *promptCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	keysToDelete := make([]string, 0)

	for key, entry := range c.data {
		if now.After(entry.expiresAt) {
			keysToDelete = append(keysToDelete, key)
		}
	}

	for _, key := range keysToDelete {
		delete(c.data, key)
	}

	if len(keysToDelete) > 0 {
		logger.Debug("cleaned up expired prompt cache entries", zap.Int("count", len(keysToDelete)))
	}
}
