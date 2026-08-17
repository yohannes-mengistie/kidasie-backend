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

