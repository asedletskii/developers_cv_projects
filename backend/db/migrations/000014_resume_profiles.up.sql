CREATE TABLE resume_profiles (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL UNIQUE REFERENCES resume_roots(id) ON DELETE CASCADE,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    position text NOT NULL,
    photo_path text NOT NULL,
    date_of_birth date NOT NULL,
    location text NOT NULL,
    relocation boolean NOT NULL,
    CONSTRAINT resume_profiles_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_profiles_surname_check CHECK (length(btrim(surname)) > 0),
    CONSTRAINT resume_profiles_position_check CHECK (length(btrim(position)) > 0),
    CONSTRAINT resume_profiles_photo_path_check CHECK (photo_path ~ '^/[^/].*$'),
    CONSTRAINT resume_profiles_location_check CHECK (length(btrim(location)) > 0)
);
CREATE INDEX idx_resume_profiles_resume_id ON resume_profiles(resume_id);

CREATE TABLE resume_profile_work_formats (
    profile_id uuid NOT NULL REFERENCES resume_profiles(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (profile_id, value),
    UNIQUE (profile_id, position),
    CONSTRAINT resume_profile_work_formats_value_check CHECK (value IN ('fullTime','partTime','office','remote','hybrid')),
    CONSTRAINT resume_profile_work_formats_position_check CHECK (position >= 0)
);
CREATE INDEX idx_resume_profile_work_formats_profile_position ON resume_profile_work_formats(profile_id, position);

INSERT INTO resume_profiles (id, resume_id, name, surname, patronymic, position, photo_path, date_of_birth, location, relocation)
SELECT gen_random_uuid(), r.id, r.name, r.surname, r.patronymic, r.position, r.photo_path, r.date_of_birth, r.location, r.relocation
FROM resumes r
JOIN resume_roots rr ON rr.id = r.id
ON CONFLICT (resume_id) DO NOTHING;

INSERT INTO resume_profile_work_formats (profile_id, value, position)
SELECT p.id, w.value, w.position
FROM resume_work_formats w
JOIN resume_profiles p ON p.resume_id = w.resume_id
ON CONFLICT (profile_id, value) DO NOTHING;
