-- ============================================
-- 多版本润色功能数据库回滚脚本
-- 版本: 001
-- 作者: Claude Code
-- 日期: 2025-12-03
-- ============================================

-- ============================================
-- 1. 删除触发器
-- ============================================
DROP TRIGGER IF EXISTS update_polish_prompts_updated_at ON polish_prompts;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- ============================================
-- 2. 删除从表 polish_versions
-- ============================================
DROP TABLE IF EXISTS polish_versions CASCADE;

-- ============================================
-- 3. 删除 Prompt 管理表 polish_prompts
-- ============================================
DROP TABLE IF EXISTS polish_prompts CASCADE;

-- ============================================
-- 4. 删除 users 表新增的字段
-- ============================================
ALTER TABLE users
DROP COLUMN IF EXISTS enable_multi_version,
DROP COLUMN IF EXISTS multi_version_quota;

-- ============================================
-- 5. 删除主表 polish_records 的 mode 字段
-- ============================================
DROP INDEX IF EXISTS idx_mode;
ALTER TABLE polish_records
DROP COLUMN IF EXISTS mode;

-- ============================================
-- 回滚完成
-- ============================================
