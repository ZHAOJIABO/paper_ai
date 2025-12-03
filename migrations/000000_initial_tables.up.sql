-- ============================================
-- 初始表结构创建脚本（PostgreSQL）
-- 版本: 000000
-- 说明: 创建基础表 users, polish_records, refresh_tokens
-- ============================================

-- ============================================
-- 1. 创建 users 表
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,  -- 使用 Snowflake ID，不自增
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    avatar_url VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    email_verified BOOLEAN DEFAULT false,
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(50),
    login_count INT DEFAULT 0,
    failed_login_count INT DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_username ON users(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_email ON users(email) WHERE deleted_at IS NULL;

-- 创建普通索引
CREATE INDEX IF NOT EXISTS idx_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_deleted_at ON users(deleted_at);

-- 添加表和列注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.id IS '用户ID (Snowflake ID)';
COMMENT ON COLUMN users.status IS '用户状态: active / suspended / deleted';

-- ============================================
-- 2. 创建 polish_records 表
-- ============================================
CREATE TABLE IF NOT EXISTS polish_records (
    id BIGSERIAL PRIMARY KEY,
    trace_id VARCHAR(20) NOT NULL,
    user_id BIGINT NOT NULL,

    -- 原始内容
    original_content TEXT NOT NULL,
    style VARCHAR(20) NOT NULL,
    language VARCHAR(10) NOT NULL,

    -- 润色结果
    polished_content TEXT NOT NULL,
    original_length INT NOT NULL,
    polished_length INT NOT NULL,

    -- AI 信息
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,

    -- 性能指标
    process_time_ms INT DEFAULT 0,

    -- 状态
    status VARCHAR(20) NOT NULL DEFAULT 'success',
    error_message TEXT,

    -- 对比数据
    comparison_data JSONB,
    changes_count INT DEFAULT 0,
    accepted_changes JSONB,
    rejected_changes JSONB,
    final_content TEXT,

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_trace_id ON polish_records(trace_id);

-- 创建普通索引
CREATE INDEX IF NOT EXISTS idx_user_id ON polish_records(user_id);
CREATE INDEX IF NOT EXISTS idx_style ON polish_records(style);
CREATE INDEX IF NOT EXISTS idx_language ON polish_records(language);
CREATE INDEX IF NOT EXISTS idx_provider ON polish_records(provider);
CREATE INDEX IF NOT EXISTS idx_process_time ON polish_records(process_time_ms);
CREATE INDEX IF NOT EXISTS idx_status_record ON polish_records(status);
CREATE INDEX IF NOT EXISTS idx_created_at_record ON polish_records(created_at);
CREATE INDEX IF NOT EXISTS idx_deleted_at_record ON polish_records(deleted_at);

-- 添加表和列注释
COMMENT ON TABLE polish_records IS '润色记录表';
COMMENT ON COLUMN polish_records.trace_id IS '润色追踪ID(纯数字)';
COMMENT ON COLUMN polish_records.comparison_data IS '对比数据JSON';
COMMENT ON COLUMN polish_records.changes_count IS '修改总数';

-- ============================================
-- 3. 创建 refresh_tokens 表
-- ============================================
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    device_id VARCHAR(100),
    user_agent VARCHAR(500),
    ip_address VARCHAR(50),
    is_revoked BOOLEAN DEFAULT false,
    revoked_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_token ON refresh_tokens(token);

-- 创建普通索引
CREATE INDEX IF NOT EXISTS idx_user_id_token ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_deleted_at_token ON refresh_tokens(deleted_at);

-- 添加表注释
COMMENT ON TABLE refresh_tokens IS '刷新令牌表';

-- ============================================
-- 4. 创建更新时间戳触发器函数
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为所有表创建更新触发器
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_polish_records_updated_at
BEFORE UPDATE ON polish_records
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_refresh_tokens_updated_at
BEFORE UPDATE ON refresh_tokens
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 迁移完成
-- ============================================
