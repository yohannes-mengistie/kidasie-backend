ALTER TABLE liturgies
ADD COLUMN status TEXT NOT NULL DEFAULT 'draft',
ADD COLUMN published_at TIMESTAMPTZ;

ALTER TABLE liturgies
ADD CONSTRAINT liturgies_status_valid
CHECK (
    status IN (
        'draft',
        'in_review',
        'published',
        'archived'
    )
);

ALTER TABLE liturgies
ADD CONSTRAINT liturgies_published_at_required
CHECK (
    status <> 'published'
    OR published_at IS NOT NULL
);

UPDATE liturgies
SET
    status = 'published',
    published_at = NOW();
