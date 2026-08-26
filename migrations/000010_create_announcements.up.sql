CREATE TABLE announcements (
    id BIGSERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    version BIGINT NOT NULL DEFAULT 1,
    title_am TEXT NOT NULL,
    title_en TEXT NOT NULL,
    body_am TEXT NOT NULL,
    body_en TEXT NOT NULL,
    kind TEXT NOT NULL DEFAULT 'general',
    action_type TEXT NOT NULL DEFAULT 'none',
    action_value TEXT,
    priority INTEGER NOT NULL DEFAULT 0,
    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    status TEXT NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT announcements_slug_valid CHECK (
        slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'
    ),
    CONSTRAINT announcements_version_positive CHECK (version > 0),
    CONSTRAINT announcements_title_am_not_empty CHECK (BTRIM(title_am) <> ''),
    CONSTRAINT announcements_title_en_not_empty CHECK (BTRIM(title_en) <> ''),
    CONSTRAINT announcements_body_am_not_empty CHECK (BTRIM(body_am) <> ''),
    CONSTRAINT announcements_body_en_not_empty CHECK (BTRIM(body_en) <> ''),
    CONSTRAINT announcements_kind_valid CHECK (
        kind IN ('general', 'content', 'audio', 'app_update', 'important')
    ),
    CONSTRAINT announcements_action_type_valid CHECK (
        action_type IN ('none', 'open_liturgy', 'download_apk', 'open_url')
    ),
    CONSTRAINT announcements_action_value_valid CHECK (
        (action_type = 'none' AND action_value IS NULL)
        OR (
            action_type <> 'none'
            AND action_value IS NOT NULL
            AND BTRIM(action_value) <> ''
        )
    ),
    CONSTRAINT announcements_priority_valid CHECK (
        priority BETWEEN 0 AND 100
    ),
    CONSTRAINT announcements_status_valid CHECK (
        status IN ('draft', 'published', 'archived')
    ),
    CONSTRAINT announcements_publication_valid CHECK (
        (status = 'published' AND published_at IS NOT NULL)
        OR status <> 'published'
    ),
    CONSTRAINT announcements_expiration_valid CHECK (
        expires_at IS NULL
        OR published_at IS NULL
        OR expires_at > published_at
    )
);

CREATE INDEX announcements_public_feed_idx
ON announcements (is_pinned DESC, priority DESC, published_at DESC, id DESC)
WHERE status = 'published';

CREATE OR REPLACE FUNCTION bump_announcement_version_on_update()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF ROW(
        OLD.slug,
        OLD.title_am,
        OLD.title_en,
        OLD.body_am,
        OLD.body_en,
        OLD.kind,
        OLD.action_type,
        OLD.action_value,
        OLD.priority,
        OLD.is_pinned,
        OLD.status,
        OLD.published_at,
        OLD.expires_at
    ) IS DISTINCT FROM ROW(
        NEW.slug,
        NEW.title_am,
        NEW.title_en,
        NEW.body_am,
        NEW.body_en,
        NEW.kind,
        NEW.action_type,
        NEW.action_value,
        NEW.priority,
        NEW.is_pinned,
        NEW.status,
        NEW.published_at,
        NEW.expires_at
    ) THEN
        NEW.version := OLD.version + 1;
        NEW.updated_at := NOW();
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER announcements_bump_version
BEFORE UPDATE ON announcements
FOR EACH ROW
EXECUTE FUNCTION bump_announcement_version_on_update();
