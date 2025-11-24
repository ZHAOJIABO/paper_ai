package service

import (
	"context"

	"paper_ai/internal/config"
	"paper_ai/internal/domain/model"
	"paper_ai/internal/infrastructure/ai"
	"paper_ai/internal/infrastructure/ai/types"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/logger"
	"go.uber.org/zap"
)

// PolishService 润色服务
type PolishService struct {
	providerFactory *ai.ProviderFactory
}

// NewPolishService 创建润色服务
func NewPolishService(factory *ai.ProviderFactory) *PolishService {
	return &PolishService{
		providerFactory: factory,
	}
}

// Polish 执行段落润色
func (s *PolishService) Polish(ctx context.Context, req *model.PolishRequest) (*types.PolishResponse, error) {
	// 参数验证
	if err := req.Validate(); err != nil {
		logger.Warn("invalid polish request", zap.Error(err))
		return nil, apperrors.NewInvalidParameterError(err.Error())
	}

	// 设置默认值
	req.SetDefaults()

	// 获取AI提供商
	var provider ai.AIProvider
	var err error

	if req.Provider == "" {
		// 使用默认提供商
		provider, err = s.providerFactory.GetDefaultProvider()
		if err != nil {
			logger.Error("failed to get default provider", zap.Error(err))
			return nil, err
		}
		req.Provider = config.Get().AI.DefaultProvider
	} else {
		// 使用指定的提供商
		provider, err = s.providerFactory.GetProvider(req.Provider)
		if err != nil {
			logger.Error("failed to get provider", zap.String("provider", req.Provider), zap.Error(err))
			return nil, err
		}
	}

	// 构建AI请求
	aiReq := &types.PolishRequest{
		Content:  req.Content,
		Style:    req.Style,
		Language: req.Language,
	}

	// 调用AI服务
	logger.Info("calling ai provider for polish",
		zap.String("provider", req.Provider),
		zap.Int("content_length", len(req.Content)),
	)

	resp, err := provider.Polish(ctx, aiReq)
	if err != nil {
		logger.Error("ai provider polish failed",
			zap.String("provider", req.Provider),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Info("polish completed successfully",
		zap.String("provider", req.Provider),
		zap.Int("original_length", resp.OriginalLength),
		zap.Int("polished_length", resp.PolishedLength),
	)

	return resp, nil
}
