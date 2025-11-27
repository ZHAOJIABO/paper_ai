package service

import (
	"context"
	"errors"
	"time"

	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/model"
	"paper_ai/internal/domain/repository"
	"paper_ai/internal/infrastructure/security"
	apperrors "paper_ai/pkg/errors"
)

// AuthService 认证服务
type AuthService struct {
	userRepo         repository.UserRepository
	tokenRepo        repository.RefreshTokenRepository
	jwtManager       *security.JWTManager
}

// NewAuthService 创建认证服务实例
func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.RefreshTokenRepository,
	jwtManager *security.JWTManager,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.UserInfo, error) {
	// 1. 验证参数
	if req.Password != req.ConfirmPassword {
		return nil, apperrors.NewBadRequestError("两次输入的密码不一致")
	}

	// 2. 验证用户名格式
	if err := security.ValidateUsername(req.Username); err != nil {
		return nil, apperrors.NewBadRequestError(err.Error())
	}

	// 3. 验证邮箱格式
	if err := security.ValidateEmail(req.Email); err != nil {
		return nil, apperrors.NewBadRequestError(err.Error())
	}

	// 4. 验证密码强度
	if err := security.ValidatePasswordStrength(req.Password); err != nil {
		return nil, apperrors.NewBadRequestError(err.Error())
	}

	// 5. 检查用户名是否已存在
	exists, err := s.userRepo.ExistsUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.NewInternalError("检查用户名失败")
	}
	if exists {
		return nil, apperrors.NewBadRequestError("用户名已存在")
	}

	// 6. 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("检查邮箱失败")
	}
	if exists {
		return nil, apperrors.NewBadRequestError("邮箱已被注册")
	}

	// 7. 加密密码
	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, apperrors.NewInternalError("密码加密失败")
	}

	// 8. 创建用户实体
	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Nickname:     req.Nickname,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if user.Nickname == "" {
		user.Nickname = req.Username
	}

	// 9. 保存用户
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperrors.NewInternalError("创建用户失败")
	}

	// 10. 返回用户信息
	return &model.UserInfo{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Nickname:      user.Nickname,
		AvatarURL:     user.AvatarURL,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest, ip, userAgent string) (*model.LoginResponse, error) {
	// 1. 查找用户（支持用户名或邮箱登录）
	var user *entity.User
	var err error

	// 尝试用户名登录
	user, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.NewInternalError("查询用户失败")
	}

	// 如果用户名不存在，尝试邮箱登录
	if user == nil {
		user, err = s.userRepo.GetByEmail(ctx, req.Username)
		if err != nil {
			return nil, apperrors.NewInternalError("查询用户失败")
		}
	}

	// 2. 用户不存在
	if user == nil {
		return nil, apperrors.NewUnauthorizedError("用户名或密码错误")
	}

	// 3. 检查用户状态
	if user.IsBanned() {
		return nil, apperrors.NewForbiddenError("账号已被封禁")
	}

	// 4. 验证密码
	if err := security.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		// 增加失败登录次数
		_ = s.userRepo.IncrementFailedLoginCount(ctx, user.ID)
		return nil, apperrors.NewUnauthorizedError("用户名或密码错误")
	}

	// 5. 生成访问令牌
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, apperrors.NewInternalError("生成访问令牌失败")
	}

	// 6. 生成刷新令牌
	refreshTokenStr, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("生成刷新令牌失败")
	}

	// 7. 保存刷新令牌到数据库
	refreshToken := &entity.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(time.Duration(s.jwtManager.GetRefreshExpiry()) * time.Second),
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.tokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, apperrors.NewInternalError("保存刷新令牌失败")
	}

	// 8. 更新登录信息
	if err := s.userRepo.UpdateLoginInfo(ctx, user.ID, ip); err != nil {
		// 记录日志但不影响登录
	}

	// 9. 返回登录响应
	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    s.jwtManager.GetAccessExpiry(),
		TokenType:    "Bearer",
		User: &model.UserInfo{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			Nickname:      user.Nickname,
			AvatarURL:     user.AvatarURL,
			Status:        user.Status,
			EmailVerified: user.EmailVerified,
			LastLoginAt:   user.LastLoginAt,
			CreatedAt:     user.CreatedAt,
		},
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*model.RefreshTokenResponse, error) {
	// 1. 验证刷新令牌格式
	claims, err := s.jwtManager.ValidateToken(refreshTokenStr)
	if err != nil {
		if errors.Is(err, security.ErrExpiredToken) {
			return nil, apperrors.NewUnauthorizedError("刷新令牌已过期")
		}
		return nil, apperrors.NewUnauthorizedError("无效的刷新令牌")
	}

	// 2. 从数据库查询令牌
	token, err := s.tokenRepo.GetByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, apperrors.NewInternalError("查询令牌失败")
	}

	// 3. 令牌不存在
	if token == nil {
		return nil, apperrors.NewUnauthorizedError("刷新令牌不存在")
	}

	// 4. 检查令牌是否有效
	if !token.IsValid() {
		return nil, apperrors.NewUnauthorizedError("刷新令牌已失效")
	}

	// 5. 获取用户信息
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.NewInternalError("查询用户失败")
	}

	if user == nil {
		return nil, apperrors.NewUnauthorizedError("用户不存在")
	}

	// 6. 检查用户状态
	if user.IsBanned() {
		return nil, apperrors.NewForbiddenError("账号已被封禁")
	}

	// 7. 生成新的访问令牌
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, apperrors.NewInternalError("生成访问令牌失败")
	}

	// 8. 返回响应
	return &model.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   s.jwtManager.GetAccessExpiry(),
		TokenType:   "Bearer",
	}, nil
}

// Logout 登出
func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string) error {
	// 1. 撤销令牌
	if err := s.tokenRepo.Revoke(ctx, refreshTokenStr); err != nil {
		return apperrors.NewInternalError("登出失败")
	}

	return nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(ctx context.Context, userID int64) (*model.UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("查询用户失败")
	}

	if user == nil {
		return nil, apperrors.NewNotFoundError("用户不存在")
	}

	return &model.UserInfo{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Nickname:      user.Nickname,
		AvatarURL:     user.AvatarURL,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
	}, nil
}
