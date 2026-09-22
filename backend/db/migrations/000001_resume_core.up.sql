CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE people (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    position text NOT NULL,
    resume_link text,
    photo_path text NOT NULL,
    date_of_birth date NOT NULL,
    location text NOT NULL,
    relocation boolean NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE person_greetings (
    person_id uuid PRIMARY KEY REFERENCES people(id) ON DELETE CASCADE,
    first_ru text NOT NULL,
    first_en text NOT NULL,
    second_ru text NOT NULL,
    second_en text NOT NULL,
    second_highlight_ru text NOT NULL,
    second_highlight_en text NOT NULL,
    third_ru text NOT NULL,
    third_en text NOT NULL
);

CREATE TABLE person_availability (
    person_id uuid PRIMARY KEY REFERENCES people(id) ON DELETE CASCADE,
    status text NOT NULL,
    available_from date,
    note_ru text,
    note_en text
);

CREATE TABLE person_themes (
    person_id uuid PRIMARY KEY REFERENCES people(id) ON DELETE CASCADE,
    primary_color text,
    secondary_color text,
    CHECK (primary_color IS NOT NULL OR secondary_color IS NOT NULL)
);

CREATE TABLE person_contacts_text (
    person_id uuid PRIMARY KEY REFERENCES people(id) ON DELETE CASCADE,
    title_ru text,
    title_en text,
    subtitle_ru text,
    subtitle_en text,
    CHECK (title_ru IS NOT NULL OR title_en IS NOT NULL OR subtitle_ru IS NOT NULL OR subtitle_en IS NOT NULL)
);

CREATE TABLE person_contacts (
    person_id uuid PRIMARY KEY REFERENCES people(id) ON DELETE CASCADE,
    phone text,
    telegram text,
    email text,
    CHECK (phone IS NOT NULL OR telegram IS NOT NULL OR email IS NOT NULL)
);
