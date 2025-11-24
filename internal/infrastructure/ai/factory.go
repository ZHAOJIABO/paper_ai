package ai

import (
	"fmt"
	"sync"

	"paper_ai/internal/config"
	"paper_ai/internal/infrastructure/ai/claude"
	"paper_ai/internal/infrastructure/ai/doubao"
	apperrors "paper_ai/pkg/errors"
)

// ProviderFactory AI提供商工厂
type ProviderFactory struct {
	providers map[string]AIProvider
	mu        sync.RWMutex
}

var (
	factoryInstance *ProviderFactory
	factoryOnce     sync.Once
)

// GetFactory 获取工厂单例
func GetFactory() *ProviderFactory {
	factoryOnce.Do(func() {
		factoryInstance = &ProviderFactory{
			providers: make(map[string]AIProvider),
		}
	})
	return factoryInstance
}

// InitProviders 初始化所有配置的提供商
func (f *ProviderFactory) InitProviders(cfg *config.Config) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 初始化所有配置的提供商
	for name, providerCfg := range cfg.AI.Providers {
		switch name {
		case "claude":
			client := claude.NewClient(
				providerCfg.APIKey,
				providerCfg.BaseURL,
				providerCfg.Model,
				providerCfg.Timeout,
			)
			f.providers[name] = client
		case "doubao":
			client := doubao.NewClient(
				providerCfg.APIKey,
				providerCfg.BaseURL,
				providerCfg.Model,
				providerCfg.Timeout,
			)
			f.providers[name] = client
		// 未来可以在这里添加其他提供商
		// case "openai":
		//     client := openai.NewClient(...)
		//     f.providers[name] = client
		default:
			return fmt.Errorf("unsupported AI provider: %s", name)
		}
	}

	return nil
}

// GetProvider 获取指定的AI提供商
func (f *ProviderFactory) GetProvider(providerName string) (AIProvider, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	provider, exists := f.providers[providerName]
	if !exists {
		return nil, apperrors.NewProviderNotFoundError(providerName)
	}

	return provider, nil
}

// GetDefaultProvider 获取默认的AI提供商
func (f *ProviderFactory) GetDefaultProvider() (AIProvider, error) {
	cfg := config.Get()
	return f.GetProvider(cfg.AI.DefaultProvider)
}

// ListProviders 列出所有已注册的提供商
func (f *ProviderFactory) ListProviders() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	names := make([]string, 0, len(f.providers))
	for name := range f.providers {
		names = append(names, name)
	}
	return names
}
