INSERT INTO user_profiles (user_id, name, surname, patronymic, position, location, photo_path)
VALUES
    ('10000000-0000-0000-0000-000000000001', 'Антон', 'Профиль', 'Демо', 'Profile Owner', 'Новосибирск', '/profiles/anton.svg'),
    ('10000000-0000-0000-0000-000000000002', 'Дмитрий', 'Профиль', NULL, 'Backend Engineer', 'Томск', '/profiles/dmitry.svg')
ON CONFLICT (user_id) DO NOTHING;
