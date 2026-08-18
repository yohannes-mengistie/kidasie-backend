DROP TRIGGER IF EXISTS verses_bump_content_version
ON verses;

DROP TRIGGER IF EXISTS sections_bump_content_version
ON sections;

DROP TRIGGER IF EXISTS liturgies_bump_content_version
ON liturgies;

DROP FUNCTION IF EXISTS bump_liturgy_version_from_verse();

DROP FUNCTION IF EXISTS bump_liturgy_version_from_section();

DROP FUNCTION IF EXISTS bump_liturgy_version_on_update();

DROP FUNCTION IF EXISTS increment_liturgy_content_version(BIGINT);
