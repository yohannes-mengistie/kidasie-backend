ALTER TABLE liturgies
DROP CONSTRAINT liturgies_content_version_positive;

ALTER TABLE liturgies
DROP COLUMN content_version;