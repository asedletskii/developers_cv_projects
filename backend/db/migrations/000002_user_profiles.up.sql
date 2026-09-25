CREATE TABLE user_profiles (
    user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    position text,
    location text,
    photo_path text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT user_profiles_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT user_profiles_surname_check CHECK (length(btrim(surname)) > 0),
    CONSTRAINT user_profiles_photo_path_check CHECK (photo_path IS NULL OR photo_path ~ '^/[^/].*$')
);

INSERT INTO user_profiles (user_id, name, surname, patronymic, position, location, photo_path)
VALUES
    ('10000000-0000-0000-0000-000000000001', 'Антон', 'Профиль', 'Демо', 'Profile Owner', 'Новосибирск', '/profiles/anton.svg'),
    ('10000000-0000-0000-0000-000000000002', 'Дмитрий', 'Профиль', NULL, 'Backend Engineer', 'Томск', '/profiles/dmitry.svg');
