DROP TRIGGER IF EXISTS announcements_bump_version ON announcements;
DROP FUNCTION IF EXISTS bump_announcement_version_on_update();
DROP TABLE IF EXISTS announcements;
