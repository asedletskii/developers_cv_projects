ALTER TABLE user_profiles
    ADD COLUMN name text,
    ADD COLUMN surname text,
    ADD COLUMN patronymic text,
    ADD COLUMN position text,
    ADD COLUMN location text;

UPDATE user_profiles p
SET name = COALESCE(
        (SELECT l.name FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'ru'),
        (SELECT l.name FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'en')
    ),
    surname = COALESCE(
        (SELECT l.surname FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'ru'),
        (SELECT l.surname FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'en')
    ),
    patronymic = COALESCE(
        (SELECT l.patronymic FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'ru'),
        (SELECT l.patronymic FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'en')
    ),
    position = COALESCE(
        (SELECT l.position FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'ru'),
        (SELECT l.position FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'en')
    ),
    location = COALESCE(
        (SELECT l.location FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'ru'),
        (SELECT l.location FROM user_profile_localizations l
         WHERE l.user_id = p.user_id AND l.language = 'en')
    );

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM user_profiles WHERE name IS NULL OR surname IS NULL) THEN
        RAISE EXCEPTION 'cannot rollback 000011: one or more profiles have no name or surname in RU/EN localizations';
    END IF;
END $$;

ALTER TABLE user_profiles
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN surname SET NOT NULL,
    ADD CONSTRAINT user_profiles_name_check CHECK (length(btrim(name)) > 0),
    ADD CONSTRAINT user_profiles_surname_check CHECK (length(btrim(surname)) > 0);

DROP TABLE user_profile_localizations;
