CREATE FUNCTION increment_liturgy_content_version(
    target_liturgy_id BIGINT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    IF target_liturgy_id IS NULL THEN
        RETURN;
    END IF;

    UPDATE liturgies
    SET
        content_version = content_version + 1,
        updated_at = NOW()
    WHERE id = target_liturgy_id;
END;
$$;

CREATE FUNCTION bump_liturgy_version_on_update()
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

CREATE FUNCTION bump_liturgy_version_from_section()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        PERFORM increment_liturgy_content_version(NEW.liturgy_id);
    ELSIF TG_OP = 'DELETE' THEN
        PERFORM increment_liturgy_content_version(OLD.liturgy_id);
    ELSIF ROW(
        OLD.liturgy_id,
        OLD.sort_order,
        OLD.title,
        OLD.title_am,
        OLD.audio_url,
        OLD.audio_duration_ms,
        OLD.audio_size_bytes,
        OLD.audio_mime_type,
        OLD.audio_sha256
    ) IS DISTINCT FROM ROW(
        NEW.liturgy_id,
        NEW.sort_order,
        NEW.title,
        NEW.title_am,
        NEW.audio_url,
        NEW.audio_duration_ms,
        NEW.audio_size_bytes,
        NEW.audio_mime_type,
        NEW.audio_sha256
    ) THEN
        PERFORM increment_liturgy_content_version(NEW.liturgy_id);

        IF OLD.liturgy_id IS DISTINCT FROM NEW.liturgy_id THEN
            PERFORM increment_liturgy_content_version(OLD.liturgy_id);
        END IF;
    END IF;

    RETURN NULL;
END;
$$;

CREATE TRIGGER sections_bump_content_version
AFTER INSERT OR UPDATE OR DELETE
ON sections
FOR EACH ROW
EXECUTE FUNCTION bump_liturgy_version_from_section();

CREATE FUNCTION bump_liturgy_version_from_verse()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    new_liturgy_id BIGINT;
    old_liturgy_id BIGINT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        SELECT liturgy_id
        INTO new_liturgy_id
        FROM sections
        WHERE id = NEW.section_id;

        PERFORM increment_liturgy_content_version(new_liturgy_id);
    ELSIF TG_OP = 'DELETE' THEN
        SELECT liturgy_id
        INTO old_liturgy_id
        FROM sections
        WHERE id = OLD.section_id;

        PERFORM increment_liturgy_content_version(old_liturgy_id);
    ELSIF ROW(
        OLD.section_id,
        OLD.sort_order,
        OLD.text_geez,
        OLD.text_am,
        OLD.text_en,
        OLD.role,
        OLD.start_ms,
        OLD.end_ms
    ) IS DISTINCT FROM ROW(
        NEW.section_id,
        NEW.sort_order,
        NEW.text_geez,
        NEW.text_am,
        NEW.text_en,
        NEW.role,
        NEW.start_ms,
        NEW.end_ms
    ) THEN
        SELECT liturgy_id
        INTO new_liturgy_id
        FROM sections
        WHERE id = NEW.section_id;

        PERFORM increment_liturgy_content_version(new_liturgy_id);

        IF OLD.section_id IS DISTINCT FROM NEW.section_id THEN
            SELECT liturgy_id
            INTO old_liturgy_id
            FROM sections
            WHERE id = OLD.section_id;

            IF old_liturgy_id IS DISTINCT FROM new_liturgy_id THEN
                PERFORM increment_liturgy_content_version(old_liturgy_id);
            END IF;
        END IF;
    END IF;

    RETURN NULL;
END;
$$;

CREATE TRIGGER verses_bump_content_version
AFTER INSERT OR UPDATE OR DELETE
ON verses
FOR EACH ROW
EXECUTE FUNCTION bump_liturgy_version_from_verse();
