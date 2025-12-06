-- 删除 selected_version 字段
ALTER TABLE polish_records
DROP COLUMN IF EXISTS selected_version;
