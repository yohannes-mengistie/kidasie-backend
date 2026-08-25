ALTER TABLE liturgies
ADD COLUMN audio_url TEXT,
ADD COLUMN audio_duration_ms INTEGER,
ADD COLUMN audio_size_bytes BIGINT,
ADD COLUMN audio_mime_type TEXT,
ADD COLUMN audio_sha256 TEXT;

ALTER TABLE liturgies
ADD CONSTRAINT liturgies_audio_url_not_empty
CHECK (audio_url IS NULL OR audio_url <> ''),
ADD CONSTRAINT liturgies_audio_duration_positive
CHECK (audio_duration_ms IS NULL OR audio_duration_ms > 0),
ADD CONSTRAINT liturgies_audio_size_positive
CHECK (audio_size_bytes IS NULL OR audio_size_bytes > 0),
ADD CONSTRAINT liturgies_audio_mime_type_valid
CHECK (audio_mime_type IS NULL OR audio_mime_type LIKE 'audio/%'),
ADD CONSTRAINT liturgies_audio_sha256_valid
CHECK (audio_sha256 IS NULL OR audio_sha256 ~ '^[0-9a-f]{64}$'),
ADD CONSTRAINT liturgies_audio_metadata_requires_url
CHECK (
    audio_url IS NOT NULL
    OR (
        audio_duration_ms IS NULL
        AND audio_size_bytes IS NULL
        AND audio_mime_type IS NULL
        AND audio_sha256 IS NULL
    )
);

CREATE OR REPLACE FUNCTION bump_liturgy_version_on_update()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF ROW(
        OLD.slug,
        OLD.name,
        OLD.name_am,
        OLD.audio_url,
        OLD.audio_duration_ms,
        OLD.audio_size_bytes,
        OLD.audio_mime_type,
        OLD.audio_sha256,
        OLD.status,
        OLD.published_at
    ) IS DISTINCT FROM ROW(
        NEW.slug,
        NEW.name,
        NEW.name_am,
        NEW.audio_url,
        NEW.audio_duration_ms,
        NEW.audio_size_bytes,
        NEW.audio_mime_type,
        NEW.audio_sha256,
        NEW.status,
        NEW.published_at
    ) THEN
        NEW.content_version := OLD.content_version + 1;
        NEW.updated_at := NOW();
    END IF;

    RETURN NEW;
END;
$$;

DROP TRIGGER liturgies_bump_content_version ON liturgies;

CREATE TRIGGER liturgies_bump_content_version
BEFORE UPDATE OF
    slug,
    name,
    name_am,
    audio_url,
    audio_duration_ms,
    audio_size_bytes,
    audio_mime_type,
    audio_sha256,
    status,
    published_at
ON liturgies
FOR EACH ROW
EXECUTE FUNCTION bump_liturgy_version_on_update();
