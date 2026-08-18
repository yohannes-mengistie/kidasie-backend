ALTER TABLE sections
ADD COLUMN audio_duration_ms INTEGER,
ADD COLUMN audio_size_bytes BIGINT,
ADD COLUMN audio_mime_type TEXT,
ADD COLUMN audio_sha256 TEXT;

ALTER TABLE sections
ADD CONSTRAINT sections_audio_duration_positive
CHECK (
    audio_duration_ms IS NULL
    OR audio_duration_ms > 0
);

ALTER TABLE sections
ADD CONSTRAINT sections_audio_size_positive
CHECK (
    audio_size_bytes IS NULL
    OR audio_size_bytes > 0
);

ALTER TABLE sections
ADD CONSTRAINT sections_audio_mime_type_valid
CHECK (
    audio_mime_type IS NULL
    OR audio_mime_type LIKE 'audio/%'
);

ALTER TABLE sections
ADD CONSTRAINT sections_audio_sha256_valid
CHECK (
    audio_sha256 IS NULL
    OR audio_sha256 ~ '^[0-9a-f]{64}$'
);

ALTER TABLE sections
ADD CONSTRAINT sections_audio_metadata_requires_url
CHECK (
    audio_url IS NOT NULL
    OR (
        audio_duration_ms IS NULL
        AND audio_size_bytes IS NULL
        AND audio_mime_type IS NULL
        AND audio_sha256 IS NULL
    )
);
