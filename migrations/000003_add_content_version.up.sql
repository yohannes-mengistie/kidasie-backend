ALTER TABLE liturgies
ADD COLUMN content_version BIGINT NOT NULL DEFAULT 1;

ALTER TABLE liturgies
ADD CONSTRAINT liturgies_content_version_positive
CHECK (content_version > 0);