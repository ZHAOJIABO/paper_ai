package handler

import (
	"strconv"

	"paper_ai/internal/domain/entity"
	"paper_ai/internal/domain/repository"
	apperrors "paper_ai/pkg/errors"
	"paper_ai/pkg/response"

	"github.com/gin-gonic/gin"
)

// PromptAdminHandler Prompt管理处理器
type PromptAdminHandler struct {
	promptRepo repository.PolishPromptRepository
}

// NewPromptAdminHandler 创建Prompt管理处理器
func NewPromptAdminHandler(promptRepo repository.PolishPromptRepository) *PromptAdminHandler {
	return &PromptAdminHandler{
		promptRepo: promptRepo,
	}
}

// ListPrompts 列出所有Prompts
// @Summary 列出Prompts
// @Description 列出所有Prompt模板（支持过滤和分页）
// @Tags admin
// @Produce json
// @Param version_type query string false "版本类型"
// @Param language query string false "语言"
// @Param style query string false "风格"
// @Param is_active query boolean false "是否激活"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=[]entity.PolishPrompt}
// @Router /api/v1/admin/prompts [get]
func (h *PromptAdminHandler) ListPrompts(c *gin.Context) {
	// 解析查询参数
	filter := repository.PromptFilter{
		VersionType: c.Query("version_type"),
		Language:    c.Query("language"),
		Style:       c.Query("style"),
		Page:        1,
		PageSize:    20,
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			filter.PageSize = pageSize
		}
	}

	prompts, err := h.promptRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, prompts)
}

// GetPrompt 获取Prompt详情
// @Summary 获取Prompt详情
// @Description 根据ID获取Prompt详情
// @Tags admin
// @Produce json
// @Param id path int true "Prompt ID"
// @Success 200 {object} response.Response{data=entity.PolishPrompt}
// @Router /api/v1/admin/prompts/{id} [get]
func (h *PromptAdminHandler) GetPrompt(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.NewInvalidParameterError("无效的用户ID"))
		return
	}

	prompt, err := h.promptRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, prompt)
}

// CreatePrompt 创建Prompt
// @Summary 创建Prompt
// @Description 创建新的Prompt模板
// @Tags admin
// @Accept json
// @Produce json
// @Param request body entity.PolishPrompt true "Prompt信息"
// @Success 200 {object} response.Response{data=entity.PolishPrompt}
// @Router /api/v1/admin/prompts [post]
func (h *PromptAdminHandler) CreatePrompt(c *gin.Context) {
	var prompt entity.PolishPrompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		response.Error(c, err)
		return
	}

	// 设置创建人
	if createdBy, exists := c.Get("username"); exists {
		prompt.CreatedBy = createdBy.(string)
	}

	if err := h.promptRepo.Create(c.Request.Context(), &prompt); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, prompt)
}

// UpdatePrompt 更新Prompt
// @Summary 更新Prompt
// @Description 更新Prompt模板
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "Prompt ID"
// @Param request body entity.PolishPrompt true "Prompt信息"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/prompts/{id} [put]
func (h *PromptAdminHandler) UpdatePrompt(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.NewInvalidParameterError("无效的ID"))
		return
	}

	var prompt entity.PolishPrompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		response.Error(c, err)
		return
	}

	prompt.ID = id

	if err := h.promptRepo.Update(c.Request.Context(), &prompt); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// DeletePrompt 删除Prompt
// @Summary 删除Prompt
// @Description 软删除Prompt（设置为不激活）
// @Tags admin
// @Produce json
// @Param id path int true "Prompt ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/prompts/{id} [delete]
func (h *PromptAdminHandler) DeletePrompt(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.NewInvalidParameterError("无效的ID"))
		return
	}

	if err := h.promptRepo.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// ActivatePrompt 激活Prompt
// @Summary 激活Prompt
// @Description 激活指定的Prompt模板
// @Tags admin
// @Produce json
// @Param id path int true "Prompt ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/prompts/{id}/activate [post]
func (h *PromptAdminHandler) ActivatePrompt(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.NewInvalidParameterError("无效的ID"))
		return
	}

	if err := h.promptRepo.Activate(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// DeactivatePrompt 停用Prompt
// @Summary 停用Prompt
// @Description 停用指定的Prompt模板
// @Tags admin
// @Produce json
// @Param id path int true "Prompt ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/prompts/{id}/deactivate [post]
func (h *PromptAdminHandler) DeactivatePrompt(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.NewInvalidParameterError("无效的ID"))
		return
	}

	if err := h.promptRepo.Deactivate(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GetPromptStats 获取Prompt统计信息
// @Summary 获取Prompt统计
// @Description 获取按版本类型统计的Prompt使用情况
// @Tags admin
// @Produce json
// @Success 200 {object} response.Response{data=map[string]repository.PromptStats}
// @Router /api/v1/admin/prompts/stats [get]
func (h *PromptAdminHandler) GetPromptStats(c *gin.Context) {
	stats, err := h.promptRepo.GetStatsByVersionType(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, stats)
}
