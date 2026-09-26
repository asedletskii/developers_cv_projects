ALTER TABLE user_profiles
    DROP CONSTRAINT IF EXISTS user_profiles_photo_path_check;

ALTER TABLE resume_profiles
    DROP CONSTRAINT IF EXISTS resume_profiles_photo_path_check;

ALTER TABLE resume_overviews
    DROP CONSTRAINT IF EXISTS resume_overviews_resume_link_check;

ALTER TABLE resume_project_screenshots
    DROP CONSTRAINT IF EXISTS resume_project_screenshots_path_check;


UPDATE user_profiles
SET photo_path = CASE photo_path
    WHEN '/profiles/anton.svg' THEN 'https://example.com/profiles/anton.svg'
    WHEN '/profiles/dmitry.svg' THEN 'https://example.com/profiles/dmitry.svg'
    ELSE photo_path
END
WHERE photo_path IN ('/profiles/anton.svg', '/profiles/dmitry.svg');

ALTER TABLE user_profiles
    ADD CONSTRAINT user_profiles_photo_path_check
        CHECK (photo_path IS NULL OR photo_path ~ '^https?://[^[:space:]]+$');

ALTER TABLE resume_profiles
    ADD CONSTRAINT resume_profiles_photo_path_check
        CHECK (photo_path ~ '^https?://[^[:space:]]+$');

ALTER TABLE resume_overviews
    ADD CONSTRAINT resume_overviews_resume_link_check
        CHECK (resume_link IS NULL OR resume_link ~ '^https?://[^[:space:]]+$');

ALTER TABLE resume_project_screenshots
    ADD CONSTRAINT resume_project_screenshots_path_check
        CHECK (path ~ '^https?://[^[:space:]]+$');
