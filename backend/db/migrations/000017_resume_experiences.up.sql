CREATE TABLE resume_experiences_v2 (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    company text NOT NULL,
    position text NOT NULL,
    location text,
    start_date date NOT NULL,
    end_date date,
    description text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_experiences_v2_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_experiences_v2_company_check CHECK (length(btrim(company)) > 0),
    CONSTRAINT resume_experiences_v2_position_check CHECK (length(btrim(position)) > 0),
    CONSTRAINT resume_experiences_v2_dates_check CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE INDEX idx_resume_experiences_v2_resume_id
    ON resume_experiences_v2(resume_id);

CREATE TABLE resume_experience_achievements_v2 (
    experience_id uuid NOT NULL REFERENCES resume_experiences_v2(id) ON DELETE CASCADE,
    position integer NOT NULL,
    value text NOT NULL,
    PRIMARY KEY (experience_id, position),
    CONSTRAINT resume_experience_achievements_v2_position_check CHECK (position >= 0),
    CONSTRAINT resume_experience_achievements_v2_value_check CHECK (length(btrim(value)) > 0)
);

CREATE TABLE resume_experience_technologies_v2 (
    experience_id uuid NOT NULL REFERENCES resume_experiences_v2(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    grade text,
    PRIMARY KEY (experience_id, position),
    CONSTRAINT resume_experience_technologies_v2_position_check CHECK (position >= 0),
    CONSTRAINT resume_experience_technologies_v2_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_experience_technologies_v2_grade_check CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'))
);

INSERT INTO resume_experiences_v2 (
    id,
    resume_id,
    sort_order,
    company,
    position,
    location,
    start_date,
    end_date,
    description
)
SELECT
    e.id,
    e.resume_id,
    e.position,
    e.company,
    e.job_position,
    e.location,
    e.start_date,
    e.end_date,
    e.description
FROM experiences e
JOIN resume_roots r ON r.id = e.resume_id;

INSERT INTO resume_experience_achievements_v2 (experience_id, position, value)
SELECT experience_id, position, value
FROM experience_achievements;

INSERT INTO resume_experience_technologies_v2 (experience_id, position, name, grade)
SELECT experience_id, position, name, grade
FROM experience_technologies;
