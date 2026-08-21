ALTER TABLE verses
    DROP CONSTRAINT verses_role_valid;

ALTER TABLE verses
    ADD CONSTRAINT verses_role_valid
    CHECK (
        role IN (
            'priest',
            'assistant_priest',
            'deacon',
            'congregation',
            'chanter',
            'reader',
            'rubric'
        )
    );
