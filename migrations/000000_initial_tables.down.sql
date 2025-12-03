-- ============================================
-- 回滚初始表结构脚本
-- 版本: 000000
-- ============================================

-- 删除触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_polish_records_updated_at ON polish_records;
DROP TRIGGER IF EXISTS update_refresh_tokens_updated_at ON refresh_tokens;

-- 删除触发器函数
DROP FUNCTION IF EXISTS update_updated_at_column();

-- 删除表（按依赖顺序）
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS polish_records CASCADE;
DROP TABLE IF EXISTS users CASCADE;
