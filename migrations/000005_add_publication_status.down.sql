ALTER TABLE liturgies
DROP CONSTRAINT liturgies_published_at_required,
DROP CONSTRAINT liturgies_status_valid;

ALTER TABLE liturgies
DROP COLUMN published_at,
DROP COLUMN status;

