DROP TRIGGER liturgies_bump_content_version ON liturgies;

CREATE OR REPLACE FUNCTION bump_liturgy_version_on_update()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF ROW(
        OLD.slug,
        OLD.name,
        OLD.name_am,
        OLD.status,
        OLD.published_at
    ) IS DISTINCT FROM ROW(
        NEW.slug,
        NEW.name,
        NEW.name_am,
        NEW.status,
        NEW.published_at
    ) THEN
        NEW.content_version := OLD.content_version + 1;
        NEW.updated_at := NOW();
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER liturgies_bump_content_version
BEFORE UPDATE OF slug, name, name_am, status, published_at
ON liturgies
FOR EACH ROW
EXECUTE FUNCTION bump_liturgy_version_on_update();

ALTER TABLE liturgies
DROP CONSTRAINT liturgies_audio_metadata_requires_url,
DROP CONSTRAINT liturgies_audio_sha256_valid,
DROP CONSTRAINT liturgies_audio_mime_type_valid,
DROP CONSTRAINT liturgies_audio_size_positive,
DROP CONSTRAINT liturgies_audio_duration_positive,
DROP CONSTRAINT liturgies_audio_url_not_empty,
DROP COLUMN audio_sha256,
DROP COLUMN audio_mime_type,
DROP COLUMN audio_size_bytes,
DROP COLUMN audio_duration_ms,
DROP COLUMN audio_url;
