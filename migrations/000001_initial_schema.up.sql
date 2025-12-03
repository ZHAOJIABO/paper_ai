-- ============================================
-- 多版本润色功能数据库迁移脚本（PostgreSQL）
-- 版本: 001
-- 作者: Claude Code
-- 日期: 2025-12-03
-- ============================================

-- ============================================
-- 1. 修改主表 polish_records，添加 mode 字段
-- ============================================
ALTER TABLE polish_records
ADD COLUMN IF NOT EXISTS mode VARCHAR(20) NOT NULL DEFAULT 'single';

-- 为 mode 字段添加索引
CREATE INDEX IF NOT EXISTS idx_mode ON polish_records(mode);

-- 添加列注释
COMMENT ON COLUMN polish_records.mode IS '润色模式: single(单版本) / multi(多版本)';

-- ============================================
-- 2. 创建从表 polish_versions (润色版本详情表)
-- ============================================
CREATE TABLE IF NOT EXISTS polish_versions (
    id BIGSERIAL PRIMARY KEY,
    record_id BIGINT NOT NULL,

    -- 版本信息
    version_type VARCHAR(32) NOT NULL,

    -- 输出内容
    polished_content TEXT NOT NULL,
    polished_length INT NOT NULL,
    suggestions JSONB,

    -- AI信息
    model_used VARCHAR(64) NOT NULL,
    prompt_id BIGINT,

    -- 性能指标
    process_time_ms INT NOT NULL DEFAULT 0,

    -- 状态
    status VARCHAR(20) NOT NULL DEFAULT 'success',
    error_message TEXT,

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- 外键约束
    CONSTRAINT fk_record_id FOREIGN KEY (record_id) REFERENCES polish_records(id) ON DELETE CASCADE,

    -- 唯一约束：同一个record_id下，version_type唯一
    CONSTRAINT uk_record_version UNIQUE (record_id, version_type)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_record_id ON polish_versions(record_id);
CREATE INDEX IF NOT EXISTS idx_version_type ON polish_versions(version_type);
CREATE INDEX IF NOT EXISTS idx_status_version ON polish_versions(status);

-- 为表和字段添加注释
COMMENT ON TABLE polish_versions IS '润色版本详情表 - 存储多版本润色的各个版本结果';
COMMENT ON COLUMN polish_versions.record_id IS '关联主表 polish_records.id';
COMMENT ON COLUMN polish_versions.version_type IS '版本类型: conservative(保守) / balanced(平衡) / aggressive(激进)';
COMMENT ON COLUMN polish_versions.polished_content IS '润色后的内容';
COMMENT ON COLUMN polish_versions.polished_length IS '润色后的长度';
COMMENT ON COLUMN polish_versions.suggestions IS '改进建议JSON数组，如: ["Changed discuss to discusses", "Added academic tone"]';
COMMENT ON COLUMN polish_versions.model_used IS '使用的AI模型';
COMMENT ON COLUMN polish_versions.prompt_id IS '关联使用的prompt模板ID';
COMMENT ON COLUMN polish_versions.process_time_ms IS '处理耗时(毫秒)';
COMMENT ON COLUMN polish_versions.status IS '状态: success / failed';
COMMENT ON COLUMN polish_versions.error_message IS '错误信息';
COMMENT ON COLUMN polish_versions.created_at IS '创建时间';

-- ============================================
-- 3. 创建 Prompt 管理表 polish_prompts
-- ============================================
CREATE TABLE IF NOT EXISTS polish_prompts (
    id BIGSERIAL PRIMARY KEY,

    -- 基本信息
    name VARCHAR(128) NOT NULL,
    version_type VARCHAR(32) NOT NULL,
    language VARCHAR(16) NOT NULL,
    style VARCHAR(32) NOT NULL,

    -- Prompt内容
    system_prompt TEXT NOT NULL,
    user_prompt_template TEXT NOT NULL,

    -- 版本管理
    version INT NOT NULL DEFAULT 1,
    is_active BOOLEAN DEFAULT true,

    -- 元数据
    description TEXT,
    tags JSONB,

    -- A/B测试
    ab_test_group VARCHAR(32),
    weight INT DEFAULT 100,

    -- 统计信息
    usage_count INT DEFAULT 0,
    success_rate DECIMAL(5,2),
    avg_satisfaction DECIMAL(3,2),

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(128),

    -- 唯一约束：version_type + language + style + version 在 is_active=true 时唯一
    CONSTRAINT uk_unique_active_prompt UNIQUE (version_type, language, style, version, is_active)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_version_type_prompt ON polish_prompts(version_type);
CREATE INDEX IF NOT EXISTS idx_language_prompt ON polish_prompts(language);
CREATE INDEX IF NOT EXISTS idx_style_prompt ON polish_prompts(style);
CREATE INDEX IF NOT EXISTS idx_active_prompt ON polish_prompts(is_active);
CREATE INDEX IF NOT EXISTS idx_created_at_prompt ON polish_prompts(created_at);

-- 添加表注释
COMMENT ON TABLE polish_prompts IS 'Prompt模板管理表 - 数据库化管理Prompt，支持版本管理和A/B测试';
COMMENT ON COLUMN polish_prompts.name IS 'Prompt名称';
COMMENT ON COLUMN polish_prompts.version_type IS '版本类型，决定润色强度: conservative / balanced / aggressive';
COMMENT ON COLUMN polish_prompts.language IS '语言: en / zh / all(通用)';
COMMENT ON COLUMN polish_prompts.style IS '风格: academic / formal / concise / all(通用)';
COMMENT ON COLUMN polish_prompts.system_prompt IS '系统提示词';
COMMENT ON COLUMN polish_prompts.user_prompt_template IS '用户提示词模板(支持变量替换)，如 {{content}}, {{style}}, {{language}}';
COMMENT ON COLUMN polish_prompts.version IS 'Prompt版本号';
COMMENT ON COLUMN polish_prompts.is_active IS '是否启用';
COMMENT ON COLUMN polish_prompts.description IS 'Prompt描述';
COMMENT ON COLUMN polish_prompts.tags IS '标签(用于分类和搜索)';
COMMENT ON COLUMN polish_prompts.ab_test_group IS 'A/B测试分组';
COMMENT ON COLUMN polish_prompts.weight IS '权重值，用于灰度发布，100=全量';
COMMENT ON COLUMN polish_prompts.usage_count IS '使用次数';
COMMENT ON COLUMN polish_prompts.success_rate IS '成功率';
COMMENT ON COLUMN polish_prompts.avg_satisfaction IS '平均满意度';
COMMENT ON COLUMN polish_prompts.created_at IS '创建时间';
COMMENT ON COLUMN polish_prompts.updated_at IS '更新时间';
COMMENT ON COLUMN polish_prompts.created_by IS '创建人';

-- ============================================
-- 4. 扩展 users 表，添加多版本功能权限字段
-- ============================================
ALTER TABLE users
ADD COLUMN IF NOT EXISTS enable_multi_version BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS multi_version_quota INT DEFAULT 0;

-- 添加索引
CREATE INDEX IF NOT EXISTS idx_enable_multi_version ON users(enable_multi_version);

-- 添加列注释
COMMENT ON COLUMN users.enable_multi_version IS '是否允许使用多版本润色功能';
COMMENT ON COLUMN users.multi_version_quota IS '多版本润色配额，0表示无限制';

-- ============================================
-- 5. 插入初始 Prompt 数据
-- ============================================

-- Conservative (保守) 版本 - English Academic
INSERT INTO polish_prompts (name, version_type, language, style, system_prompt, user_prompt_template, version, description, tags)
VALUES (
    'Conservative English Academic',
    'conservative',
    'en',
    'academic',
    'You are an academic writing assistant. Your task is to polish the text while maintaining the original meaning and structure. Make only necessary corrections for grammar, punctuation, and clarity. Keep changes minimal and conservative.',
    'Please polish the following academic text in a conservative manner. Only fix grammatical errors and improve clarity without changing the original structure or meaning:

{{content}}',
    1,
    '保守版本 - 英文学术风格，仅修正语法错误和提升清晰度',
    '["conservative", "english", "academic", "minimal-changes"]'::jsonb
) ON CONFLICT (version_type, language, style, version, is_active) DO NOTHING;

-- Balanced (平衡) 版本 - English Academic
INSERT INTO polish_prompts (name, version_type, language, style, system_prompt, user_prompt_template, version, description, tags)
VALUES (
    'Balanced English Academic',
    'balanced',
    'en',
    'academic',
    'You are an academic writing assistant. Your task is to polish the text with moderate improvements. Fix grammar, enhance clarity, and improve sentence structure while maintaining the core message and most of the original phrasing.',
    'Please polish the following academic text in a balanced manner. Fix grammar, improve clarity and sentence structure, and enhance the academic tone:

{{content}}',
    1,
    '平衡版本 - 英文学术风格，适度优化语法、结构和学术性',
    '["balanced", "english", "academic", "moderate-changes"]'::jsonb
) ON CONFLICT (version_type, language, style, version, is_active) DO NOTHING;

-- Aggressive (激进) 版本 - English Academic
INSERT INTO polish_prompts (name, version_type, language, style, system_prompt, user_prompt_template, version, description, tags)
VALUES (
    'Aggressive English Academic',
    'aggressive',
    'en',
    'academic',
    'You are an academic writing assistant. Your task is to significantly enhance the text. Rewrite sentences for better flow, use more sophisticated academic vocabulary, improve logical structure, and elevate the overall academic quality. Feel free to make substantial changes while preserving the core arguments and data.',
    'Please polish the following academic text in an aggressive manner. Significantly improve the writing quality, use sophisticated academic language, enhance logical flow, and elevate the overall academic standard:

{{content}}',
    1,
    '激进版本 - 英文学术风格，大幅提升写作质量和学术水平',
    '["aggressive", "english", "academic", "substantial-changes"]'::jsonb
) ON CONFLICT (version_type, language, style, version, is_active) DO NOTHING;

-- Conservative (保守) 版本 - Chinese Academic
INSERT INTO polish_prompts (name, version_type, language, style, system_prompt, user_prompt_template, version, description, tags)
VALUES (
    'Conservative Chinese Academic',
    'conservative',
    'zh',
    'academic',
    '你是一位学术写作助手。你的任务是在保持原文含义和结构的前提下对文本进行润色。仅进行必要的语法、标点和清晰度修正。保持改动最小化和保守性。',
    '请以保守的方式润色以下学术文本。仅修正语法错误和提升清晰度，不要改变原文结构或含义：

{{content}}',
    1,
    '保守版本 - 中文学术风格，仅修正语法错误和提升清晰度',
    '["conservative", "chinese", "academic", "minimal-changes"]'::jsonb
) ON CONFLICT (version_type, language, style, version, is_active) DO NOTHING;

-- Balanced (平衡) 版本 - Chinese Academic
INSERT INTO polish_prompts (name, version_type, language, style, system_prompt, user_prompt_template, version, description, tags)
VALUES (
    'Balanced Chinese Academic',
    'balanced',
    'zh',
    'academic',
    '你是一位学术写作助手。你的任务是对文本进行适度的改进。修正语法、提升清晰度、改善句子结构，同时保持核心信息和大部分原始表述。',
    '请以平衡的方式润色以下学术文本。修正语法、改善清晰度和句子结构，并提升学术语气：

{{content}}',
    1,
    '平衡版本 - 中文学术风格，适度优化语法、结构和学术性',
    '["balanced", "chinese", "academic", "moderate-changes"]'::jsonb
) ON CONFLICT (version_type, language, style, version, is_active) DO NOTHING;

-- Aggressive (激进) 版本 - Chinese Academic
INSERT INTO polish_prompts (name, version_type, language, style, system_prompt, user_prompt_template, version, description, tags)
VALUES (
    'Aggressive Chinese Academic',
    'aggressive',
    'zh',
    'academic',
    '你是一位学术写作助手。你的任务是显著提升文本质量。重写句子以改善流畅度，使用更精准的学术词汇，优化逻辑结构，全面提升学术水平。可以进行大幅改动，同时保留核心论点和数据。',
    '请以激进的方式润色以下学术文本。显著提升写作质量，使用精准的学术语言，增强逻辑流畅度，全面提升学术水准：

{{content}}',
    1,
    '激进版本 - 中文学术风格，大幅提升写作质量和学术水平',
    '["aggressive", "chinese", "academic", "substantial-changes"]'::jsonb
) ON CONFLICT (version_type, language, style, version, is_active) DO NOTHING;

-- ============================================
-- 6. 创建更新时间戳触发器（for polish_prompts）
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_polish_prompts_updated_at
BEFORE UPDATE ON polish_prompts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 迁移完成
-- ============================================
-- 说明:
-- 1. 主表 polish_records 添加了 mode 字段，用于区分单版本和多版本
-- 2. 新建从表 polish_versions，存储多版本润色的每个版本详情
-- 3. 新建 polish_prompts 表，实现 Prompt 数据库管理
-- 4. 扩展 users 表，添加多版本功能权限控制
-- 5. 插入了 6 个初始 Prompt 模板（3种版本 × 2种语言）
-- 6. 所有表都添加了合适的索引和约束
-- 7. 使用 ON DELETE CASCADE 确保数据一致性
-- ============================================
