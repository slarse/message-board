ALTER TABLE message DROP COLUMN revisions;
ALTER TABLE message ADD COLUMN content TEXT NOT NULL;
ALTER TABLE message ADD COLUMN title TEXT NOT NULL;
DROP FUNCTION get_effective_username(INT, JSONB);