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

