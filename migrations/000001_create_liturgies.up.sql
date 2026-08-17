CREATE TABLE liturgies (
      id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
      slug TEXT NOT NULL UNIQUE,
      name TEXT NOT NULL,
      name_am TEXT NOT NULL,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

      CONSTRAINT liturgies_slug_not_empty CHECK (slug <> ''),
      CONSTRAINT liturgies_name_not_empty CHECK (name <> ''),
      CONSTRAINT liturgies_name_am_not_empty CHECK (name_am <> '')
  );
