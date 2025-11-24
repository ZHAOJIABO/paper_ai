package ai

import (
	"context"

	"paper_ai/internal/infrastructure/ai/types"
)

// AIProvider AI提供商接口
type AIProvider interface {
	// Polish 段落润色
	Polish(ctx context.Context, req *types.PolishRequest) (*types.PolishResponse, error)

	// 预留未来扩展的接口
	// GenerateCode(ctx context.Context, req *CodeGenRequest) (*CodeGenResponse, error)
	// AnalyzeData(ctx context.Context, req *DataAnalysisRequest) (*DataAnalysisResponse, error)
}
