-- 添加 selected_version 字段到 polish_records 表（如果不存在）
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'polish_records' AND column_name = 'selected_version'
    ) THEN
        ALTER TABLE polish_records
        ADD COLUMN selected_version VARCHAR(20);

        COMMENT ON COLUMN polish_records.selected_version IS '用户选择的版本类型(多版本模式下使用)';
    END IF;
END $$;
