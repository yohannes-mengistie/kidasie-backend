INSERT INTO liturgies (
      slug,
      name,
      name_am
  )
  VALUES
      (
          'apostles',
          'Anaphora of the Apostles',
          'የሐዋርያት ቅዳሴ'
      ),
      (
          'our-lord',
          'Anaphora of Our Lord',
          'የጌታችን ቅዳሴ'
      )
  ON CONFLICT (slug) DO UPDATE
  SET
      name = EXCLUDED.name,
      name_am = EXCLUDED.name_am,
      updated_at = NOW();


INSERT INTO sections (
      liturgy_id,
      sort_order,
      title,
      title_am,
      audio_url
  )
  SELECT
      id,
      1,
      'Preparation',
      'ዝግጅት',
      NULL
  FROM liturgies
  WHERE slug = 'apostles'
  ON CONFLICT (liturgy_id, sort_order) DO UPDATE
  SET
      title = EXCLUDED.title,
      title_am = EXCLUDED.title_am,
      audio_url = EXCLUDED.audio_url,
      updated_at = NOW();

  INSERT INTO sections (
      liturgy_id,
      sort_order,
      title,
      title_am,
      audio_url
  )
  SELECT
      id,
      2,
      'Opening Prayer',
      'የመክፈቻ ጸሎት',
      NULL
  FROM liturgies
  WHERE slug = 'apostles'
  ON CONFLICT (liturgy_id, sort_order) DO UPDATE
  SET
      title = EXCLUDED.title,
      title_am = EXCLUDED.title_am,
      audio_url = EXCLUDED.audio_url,
      updated_at = NOW();

 INSERT INTO verses (
      section_id,
      sort_order,
      text_geez,
      text_am,
      text_en,
      role,
      start_ms,
      end_ms
  )
  SELECT
      s.id,
      values_to_insert.sort_order,
      values_to_insert.text_geez,
      values_to_insert.text_am,
      values_to_insert.text_en,
      values_to_insert.role,
      values_to_insert.start_ms,
      values_to_insert.end_ms
  FROM sections AS s
  INNER JOIN liturgies AS l
      ON l.id = s.liturgy_id
  CROSS JOIN (
      VALUES
          (
              1,
              '[DEVELOPMENT GEʽEZ TEXT 1]',
              '[የሙከራ አማርኛ ጽሑፍ ፩]',
              'Development English text 1',
              'priest',
              0,
              5000
          ),
          (
              2,
              '[DEVELOPMENT GEʽEZ TEXT 2]',
              '[የሙከራ አማርኛ ጽሑፍ ፪]',
              'Development English text 2',
              'congregation',
              5000,
              9000
          )
  ) AS values_to_insert (
      sort_order,
      text_geez,
      text_am,
      text_en,
      role,
      start_ms,
      end_ms
  )
  WHERE
      l.slug = 'apostles'
      AND s.sort_order = 1
  ON CONFLICT (section_id, sort_order) DO UPDATE
  SET
      text_geez = EXCLUDED.text_geez,
      text_am = EXCLUDED.text_am,
      text_en = EXCLUDED.text_en,
      role = EXCLUDED.role,
      start_ms = EXCLUDED.start_ms,
      end_ms = EXCLUDED.end_ms,
      updated_at = NOW();
