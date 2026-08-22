ALTER TABLE verses
    DROP CONSTRAINT verses_start_non_negative,
    DROP CONSTRAINT verses_end_after_start,
    DROP CONSTRAINT verses_role_valid;

ALTER TABLE verses
    ALTER COLUMN start_ms DROP NOT NULL,
    ALTER COLUMN end_ms DROP NOT NULL,
    ADD COLUMN source_page INTEGER,
    ADD COLUMN source_part TEXT,
    ADD COLUMN source_kind TEXT,
    ADD COLUMN source_note TEXT,
    ADD COLUMN source_needs_review BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE verses
    ADD CONSTRAINT verses_timing_pair_valid
    CHECK (
        (start_ms IS NULL AND end_ms IS NULL)
        OR (
            start_ms IS NOT NULL
            AND end_ms IS NOT NULL
            AND start_ms >= 0
            AND end_ms > start_ms
        )
    ),
    ADD CONSTRAINT verses_source_page_positive
    CHECK (source_page IS NULL OR source_page > 0),
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
            'rubric',
            'mixed'
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
        OLD.end_ms,
        OLD.source_page,
        OLD.source_part,
        OLD.source_kind,
        OLD.source_note,
        OLD.source_needs_review
    ) IS DISTINCT FROM ROW(
        NEW.section_id,
        NEW.sort_order,
        NEW.text_geez,
        NEW.text_am,
        NEW.text_en,
        NEW.role,
        NEW.start_ms,
        NEW.end_ms,
        NEW.source_page,
        NEW.source_part,
        NEW.source_kind,
        NEW.source_note,
        NEW.source_needs_review
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
