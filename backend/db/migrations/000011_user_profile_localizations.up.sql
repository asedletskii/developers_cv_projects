CREATE TABLE user_profile_localizations (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language text NOT NULL,
    name text,
    surname text,
    patronymic text,
    position text,
    location text,
    CONSTRAINT user_profile_localizations_pkey PRIMARY KEY (user_id, language),
    CONSTRAINT user_profile_localizations_language_check CHECK (language IN ('ru', 'en')),
    CONSTRAINT user_profile_localizations_name_check CHECK (name IS NULL OR length(btrim(name)) > 0),
    CONSTRAINT user_profile_localizations_surname_check CHECK (surname IS NULL OR length(btrim(surname)) > 0)
);

INSERT INTO user_profile_localizations (user_id, language, name, surname, patronymic, position, location)
SELECT user_id, 'ru', name, surname, patronymic, position, location
FROM user_profiles;

ALTER TABLE user_profiles
    DROP CONSTRAINT user_profiles_name_check,
    DROP CONSTRAINT user_profiles_surname_check,
    DROP CONSTRAINT user_profiles_photo_path_check,
    DROP COLUMN name,
    DROP COLUMN surname,
    DROP COLUMN patronymic,
    DROP COLUMN position,
    DROP COLUMN location;

ALTER TABLE user_profiles
    ADD CONSTRAINT user_profiles_photo_path_check CHECK (photo_path IS NULL OR photo_path ~ '^/[^/].*$');
