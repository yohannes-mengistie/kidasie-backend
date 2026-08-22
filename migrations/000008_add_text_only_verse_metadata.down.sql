ALTER TABLE verses
    DROP CONSTRAINT verses_timing_pair_valid,
    DROP CONSTRAINT verses_source_page_positive,
    DROP CONSTRAINT verses_role_valid;

ALTER TABLE verses
    DROP COLUMN source_page,
    DROP COLUMN source_part,
    DROP COLUMN source_kind,
    DROP COLUMN source_note,
    DROP COLUMN source_needs_review,
    ALTER COLUMN start_ms SET NOT NULL,
    ALTER COLUMN end_ms SET NOT NULL;

ALTER TABLE verses
    ADD CONSTRAINT verses_start_non_negative
    CHECK (start_ms >= 0),
    ADD CONSTRAINT verses_end_after_start
    CHECK (end_ms > start_ms),
    ADD CONSTRAINT verses_role_valid
    CHECK (
        role IN (
            'priest',
            'assistant_priest',
            'deacon',
            'assistant_deacon',
            'congregation',
            'chanter',
            'reader',
            'rubric'
        )
    );

CREATE OR REPLACE FUNCTION bump_liturgy_version_from_verse()
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
