ALTER TABLE user_profiles
    DROP CONSTRAINT IF EXISTS user_profiles_photo_path_check;

ALTER TABLE resume_profiles
    DROP CONSTRAINT IF EXISTS resume_profiles_photo_path_check;

ALTER TABLE resume_overviews
    DROP CONSTRAINT IF EXISTS resume_overviews_resume_link_check;

ALTER TABLE resume_courses
    DROP CONSTRAINT IF EXISTS resume_courses_certificate_link_check;

ALTER TABLE resume_activities
    DROP CONSTRAINT IF EXISTS resume_activities_link_check;

ALTER TABLE resume_projects
    DROP CONSTRAINT IF EXISTS resume_projects_link_check;

ALTER TABLE resume_projects
    DROP CONSTRAINT IF EXISTS resume_projects_git_link_check;

ALTER TABLE resume_project_screenshots
    DROP CONSTRAINT IF EXISTS resume_project_screenshots_path_check;

ALTER TABLE user_profiles
    ADD CONSTRAINT user_profiles_photo_path_check
        CHECK (photo_path IS NULL OR photo_path ~ '^https?://[^/?#[:space:]]+([/?#].*)?$');

ALTER TABLE resume_profiles
    ADD CONSTRAINT resume_profiles_photo_path_check
        CHECK (photo_path ~ '^https?://[^/?#[:space:]]+([/?#].*)?$');

ALTER TABLE resume_overviews
    ADD CONSTRAINT resume_overviews_resume_link_check
        CHECK (resume_link IS NULL OR resume_link ~ '^https?://[^/?#[:space:]]+([/?#].*)?$');

ALTER TABLE resume_courses
    ADD CONSTRAINT resume_courses_certificate_link_check
        CHECK (certificate_link IS NULL OR certificate_link ~ '^https?://[^/?#[:space:]]+([/?#].*)?$');

ALTER TABLE resume_activities
    ADD CONSTRAINT resume_activities_link_check
        CHECK (link IS NULL OR link ~ '^https?://[^/?#[:space:]]+([/?#].*)?$');

ALTER TABLE resume_projects
    ADD CONSTRAINT resume_projects_link_check
        CHECK (link IS NULL OR link ~ '^https?://[^/?#[:space:]]+([/?#].*)?$');

ALTER TABLE resume_projects
    ADD CONSTRAINT resume_projects_git_link_check
        CHECK (git_link IS NULL OR git_link ~ '^https?://[^/?#[:space:]]+([/?#].*)?$');

ALTER TABLE resume_project_screenshots
    ADD CONSTRAINT resume_project_screenshots_path_check
        CHECK (path ~ '^https?://[^/?#[:space:]]+([/?#].*)?$');

ALTER TABLE resume_experience_technologies
    RENAME TO resume_experience_stack;

ALTER TABLE resume_experience_stack
    RENAME CONSTRAINT resume_experience_technologies_pkey
    TO resume_experience_stack_pkey;

ALTER TABLE resume_experience_stack
    RENAME CONSTRAINT resume_experience_technologies_experience_id_fkey
    TO resume_experience_stack_experience_id_fkey;

ALTER TABLE resume_experience_stack
    RENAME CONSTRAINT resume_experience_technologies_position_check
    TO resume_experience_stack_position_check;

ALTER TABLE resume_experience_stack
    RENAME CONSTRAINT resume_experience_technologies_name_check
    TO resume_experience_stack_name_check;

ALTER TABLE resume_experience_stack
    RENAME CONSTRAINT resume_experience_technologies_grade_check
    TO resume_experience_stack_grade_check;

ALTER TABLE resume_project_contributors
    ADD CONSTRAINT resume_project_contributors_contacts_check
    CHECK (
        NULLIF(BTRIM(phone), '') IS NOT NULL
        OR NULLIF(BTRIM(telegram), '') IS NOT NULL
        OR NULLIF(BTRIM(email), '') IS NOT NULL
    );

UPDATE resume_contacts
SET phone = NULLIF(BTRIM(phone), ''),
    telegram = NULLIF(BTRIM(telegram), ''),
    email = NULLIF(BTRIM(email), '');

DELETE FROM resume_contacts
WHERE phone IS NULL
  AND telegram IS NULL
  AND email IS NULL;

UPDATE resume_contacts_text
SET title = NULLIF(BTRIM(title), ''),
    subtitle = NULLIF(BTRIM(subtitle), '');

DELETE FROM resume_contacts_text
WHERE title IS NULL
  AND subtitle IS NULL;

ALTER TABLE resume_contacts
    DROP CONSTRAINT resume_contacts_presence_check,
    ADD CONSTRAINT resume_contacts_presence_check CHECK (
        NULLIF(BTRIM(phone), '') IS NOT NULL
        OR NULLIF(BTRIM(telegram), '') IS NOT NULL
        OR NULLIF(BTRIM(email), '') IS NOT NULL
    );

ALTER TABLE resume_contacts_text
    DROP CONSTRAINT resume_contacts_text_presence_check,
    ADD CONSTRAINT resume_contacts_text_presence_check CHECK (
        NULLIF(BTRIM(title), '') IS NOT NULL
        OR NULLIF(BTRIM(subtitle), '') IS NOT NULL
    );
