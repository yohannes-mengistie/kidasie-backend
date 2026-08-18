ALTER TABLE sections
DROP CONSTRAINT sections_audio_metadata_requires_url,
DROP CONSTRAINT sections_audio_sha256_valid,
DROP CONSTRAINT sections_audio_mime_type_valid,
DROP CONSTRAINT sections_audio_size_positive,
DROP CONSTRAINT sections_audio_duration_positive;

ALTER TABLE sections
DROP COLUMN audio_sha256,
DROP COLUMN audio_mime_type,
DROP COLUMN audio_size_bytes,
DROP COLUMN audio_duration_ms;
