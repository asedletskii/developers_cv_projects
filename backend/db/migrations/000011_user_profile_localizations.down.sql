ALTER TABLE user_profiles
    ADD COLUMN name text,
    ADD COLUMN surname text,
    ADD COLUMN patronymic text,
    ADD COLUMN position text,
    ADD COLUMN location text;

UPDATE user_profiles p
SET name = l.name,
    surname = l.surname,
    patronymic = l.patronymic,
    position = l.position,
    location = l.location
FROM user_profile_localizations l
WHERE l.user_id = p.user_id AND l.language = 'ru';

UPDATE user_profiles p
SET name = l.name,
    surname = l.surname,
    patronymic = l.patronymic,
    position = l.position,
    location = l.location
FROM user_profile_localizations l
WHERE l.user_id = p.user_id AND l.language = 'en' AND p.name IS NULL AND p.surname IS NULL;

DELETE FROM user_profiles
WHERE name IS NULL OR surname IS NULL;

ALTER TABLE user_profiles
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN surname SET NOT NULL,
    ADD CONSTRAINT user_profiles_name_check CHECK (length(btrim(name)) > 0),
    ADD CONSTRAINT user_profiles_surname_check CHECK (length(btrim(surname)) > 0);

DROP TABLE user_profile_localizations;
