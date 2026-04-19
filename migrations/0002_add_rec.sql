ALTER TABLE tasks 
ADD COLUMN IF NOT EXISTS recurrence_type VARCHAR(20) DEFAULT 'none',
ADD COLUMN IF NOT EXISTS recurrence_rule JSONB DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_recurrence_type ON tasks(recurrence_type);