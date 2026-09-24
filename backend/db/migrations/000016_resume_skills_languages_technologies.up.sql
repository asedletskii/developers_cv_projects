CREATE TABLE resume_skills_v2 (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL UNIQUE REFERENCES resume_roots(id) ON DELETE CASCADE
);

CREATE INDEX idx_resume_skills_v2_resume_id ON resume_skills_v2(resume_id);

CREATE TABLE resume_skill_hard_v2 (
    skills_id uuid NOT NULL REFERENCES resume_skills_v2(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (skills_id, value),
    UNIQUE (skills_id, position),
    CONSTRAINT resume_skill_hard_v2_value_check CHECK (length(btrim(value)) > 0),
    CONSTRAINT resume_skill_hard_v2_position_check CHECK (position >= 0)
);

CREATE TABLE resume_skill_soft_v2 (
    skills_id uuid NOT NULL REFERENCES resume_skills_v2(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (skills_id, value),
    UNIQUE (skills_id, position),
    CONSTRAINT resume_skill_soft_v2_value_check CHECK (length(btrim(value)) > 0),
    CONSTRAINT resume_skill_soft_v2_position_check CHECK (position >= 0)
);

CREATE TABLE resume_languages_v2 (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    language text NOT NULL,
    position integer NOT NULL,
    grade text,
    UNIQUE (resume_id, language),
    UNIQUE (resume_id, position),
    CONSTRAINT resume_languages_v2_position_check CHECK (position >= 0),
    CONSTRAINT resume_languages_v2_language_check CHECK (language IN ('ru', 'en', 'fr', 'de', 'kz', 'zh', 'es')),
    CONSTRAINT resume_languages_v2_grade_check CHECK (grade IS NULL OR grade IN ('A1', 'A2', 'B1', 'B2', 'C1', 'C2'))
);

CREATE INDEX idx_resume_languages_v2_resume_position
    ON resume_languages_v2(resume_id, position);

CREATE TABLE resume_technologies_v2 (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    name text NOT NULL,
    position integer NOT NULL,
    grade text,
    UNIQUE (resume_id, position),
    CONSTRAINT resume_technologies_v2_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_technologies_v2_position_check CHECK (position >= 0),
    CONSTRAINT resume_technologies_v2_grade_check CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'))
);

CREATE INDEX idx_resume_technologies_v2_resume_position
    ON resume_technologies_v2(resume_id, position);

INSERT INTO resume_skills_v2 (id, resume_id)
SELECT gen_random_uuid(), r.id
FROM resume_roots r
WHERE EXISTS (
    SELECT 1
    FROM resume_hard_skills h
    WHERE h.resume_id = r.id
)
OR EXISTS (
    SELECT 1
    FROM resume_soft_skills s
    WHERE s.resume_id = r.id
);

INSERT INTO resume_skill_hard_v2 (skills_id, value, position)
SELECT s.id, h.value, h.position
FROM resume_hard_skills h
JOIN resume_skills_v2 s ON s.resume_id = h.resume_id;

INSERT INTO resume_skill_soft_v2 (skills_id, value, position)
SELECT s.id, h.value, h.position
FROM resume_soft_skills h
JOIN resume_skills_v2 s ON s.resume_id = h.resume_id;

INSERT INTO resume_languages_v2 (id, resume_id, language, position, grade)
SELECT gen_random_uuid(), resume_id, language, position, grade
FROM resume_languages;

INSERT INTO resume_technologies_v2 (id, resume_id, name, position, grade)
SELECT gen_random_uuid(), resume_id, name, position, grade
FROM resume_technologies;
