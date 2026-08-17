CREATE TABLE sections (
      id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      liturgy_id BIGINT NOT NULL,
      sort_order INTEGER NOT NULL,
      title TEXT NOT NULL,
      title_am TEXT NOT NULL,
      audio_url TEXT,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

      CONSTRAINT sections_liturgy_fk
          FOREIGN KEY (liturgy_id)
          REFERENCES liturgies (id)
          ON DELETE CASCADE,

      CONSTRAINT sections_sort_order_positive
          CHECK (sort_order > 0),

      CONSTRAINT sections_title_not_empty
          CHECK (title <> ''),

      CONSTRAINT sections_title_am_not_empty
          CHECK (title_am <> ''),

      CONSTRAINT sections_liturgy_sort_order_unique
          UNIQUE (liturgy_id, sort_order)
  );

  CREATE TABLE verses (
      id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      section_id BIGINT NOT NULL,
      sort_order INTEGER NOT NULL,
      text_geez TEXT NOT NULL,
      text_am TEXT NOT NULL,
      text_en TEXT,
      role TEXT NOT NULL,
      start_ms INTEGER NOT NULL,
      end_ms INTEGER NOT NULL,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

      CONSTRAINT verses_section_fk
          FOREIGN KEY (section_id)
          REFERENCES sections (id)
          ON DELETE CASCADE,

      CONSTRAINT verses_sort_order_positive
          CHECK (sort_order > 0),

      CONSTRAINT verses_start_non_negative
          CHECK (start_ms >= 0),

      CONSTRAINT verses_end_after_start
          CHECK (end_ms > start_ms),

      CONSTRAINT verses_role_valid
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
          ),

      CONSTRAINT verses_section_sort_order_unique
          UNIQUE (section_id, sort_order)
  );
