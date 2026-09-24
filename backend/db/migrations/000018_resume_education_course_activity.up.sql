CREATE TABLE resume_educations_v2 (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    location text NOT NULL,
    university text NOT NULL,
    grade integer NOT NULL,
    program_name text NOT NULL,
    profile_name text,
    faculty text NOT NULL,
    course integer,
    end_date date,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_educations_v2_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_educations_v2_location_check CHECK (length(btrim(location)) > 0),
    CONSTRAINT resume_educations_v2_university_check CHECK (length(btrim(university)) > 0),
    CONSTRAINT resume_educations_v2_grade_check CHECK (grade BETWEEN 0 AND 6),
    CONSTRAINT resume_educations_v2_program_name_check CHECK (length(btrim(program_name)) > 0),
    CONSTRAINT resume_educations_v2_faculty_check CHECK (length(btrim(faculty)) > 0),
    CONSTRAINT resume_educations_v2_course_check CHECK (course IS NULL OR course > 0)
);

CREATE INDEX idx_resume_educations_v2_resume_id
    ON resume_educations_v2(resume_id);

CREATE TABLE resume_courses_v2 (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    name text NOT NULL,
    organization text NOT NULL,
    end_date date NOT NULL,
    hours integer,
    certificate_link text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_courses_v2_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_courses_v2_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_courses_v2_organization_check CHECK (length(btrim(organization)) > 0),
    CONSTRAINT resume_courses_v2_hours_check CHECK (hours IS NULL OR hours > 0)
);

CREATE INDEX idx_resume_courses_v2_resume_id
    ON resume_courses_v2(resume_id);

CREATE TABLE resume_activities_v2 (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    type text NOT NULL,
    title text NOT NULL,
    description text,
    date date,
    link text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_activities_v2_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_activities_v2_type_check CHECK (type IN ('openSource', 'hackathon', 'speaker', 'volunteering', 'mentoring', 'community')),
    CONSTRAINT resume_activities_v2_title_check CHECK (length(btrim(title)) > 0)
);

CREATE INDEX idx_resume_activities_v2_resume_id
    ON resume_activities_v2(resume_id);

INSERT INTO resume_educations_v2 (
    id, resume_id, sort_order, location, university, grade, program_name, profile_name, faculty, course, end_date
)
SELECT
    e.id, e.resume_id, e.position, e.location, e.university, e.grade, e.program_name, e.profile_name, e.faculty, e.course, e.end_date
FROM educations e
JOIN resume_roots r ON r.id = e.resume_id;

INSERT INTO resume_courses_v2 (
    id, resume_id, sort_order, name, organization, end_date, hours, certificate_link
)
SELECT
    c.id, c.resume_id, c.position, c.name, c.organization, c.end_date, c.hours, c.certificate_link
FROM courses c
JOIN resume_roots r ON r.id = c.resume_id;

INSERT INTO resume_activities_v2 (
    id, resume_id, sort_order, type, title, description, date, link
)
SELECT
    a.id, a.resume_id, a.position, a.type, a.title, a.description, a.date, a.link
FROM activities a
JOIN resume_roots r ON r.id = a.resume_id;
