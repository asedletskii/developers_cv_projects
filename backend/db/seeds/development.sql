INSERT INTO users (id)
VALUES
    ('10000000-0000-0000-0000-000000000001'),
    ('10000000-0000-0000-0000-000000000002')
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_profiles (user_id, photo_path)
VALUES
    ('10000000-0000-0000-0000-000000000001', 'https://example.com/profiles/anton.svg'),
    ('10000000-0000-0000-0000-000000000002', 'https://example.com/profiles/dmitry.svg')
ON CONFLICT (user_id) DO UPDATE SET photo_path = EXCLUDED.photo_path;

INSERT INTO user_profile_localizations (user_id, language, name, surname, patronymic, position, location)
VALUES
    ('10000000-0000-0000-0000-000000000001', 'ru', 'Антон', 'Профиль', 'Демо', 'Profile Owner', 'Новосибирск'),
    ('10000000-0000-0000-0000-000000000002', 'ru', 'Дмитрий', 'Профиль', NULL, 'Backend Engineer', 'Томск')
ON CONFLICT (user_id, language) DO UPDATE SET
    name = EXCLUDED.name,
    surname = EXCLUDED.surname,
    patronymic = EXCLUDED.patronymic,
    position = EXCLUDED.position,
    location = EXCLUDED.location;
