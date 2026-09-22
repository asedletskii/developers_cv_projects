CREATE TABLE experiences (
    id uuid PRIMARY KEY,
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    position integer NOT NULL,
    company text NOT NULL,
    job_position text NOT NULL,
    location text,
    start_date date NOT NULL,
    end_date date,
    description text,
    UNIQUE (person_id, position)
);

CREATE TABLE experience_achievements (
    experience_id uuid NOT NULL REFERENCES experiences(id) ON DELETE CASCADE,
    position integer NOT NULL,
    value text NOT NULL,
    PRIMARY KEY (experience_id, position)
);

CREATE TABLE experience_technologies (
    experience_id uuid NOT NULL REFERENCES experiences(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    grade text,
    PRIMARY KEY (experience_id, position)
);

CREATE TABLE educations (
    id uuid PRIMARY KEY,
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    position integer NOT NULL,
    location text NOT NULL,
    university text NOT NULL,
    grade integer NOT NULL,
    program_name text NOT NULL,
    profile_name text,
    faculty text NOT NULL,
    course integer,
    end_date date,
    UNIQUE (person_id, position)
);

CREATE TABLE courses (
    id uuid PRIMARY KEY,
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    organization text NOT NULL,
    end_date date NOT NULL,
    hours integer,
    certificate_link text,
    UNIQUE (person_id, position)
);

CREATE TABLE activities (
    id uuid PRIMARY KEY,
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    position integer NOT NULL,
    type text NOT NULL,
    title text NOT NULL,
    description text,
    date date,
    link text,
    UNIQUE (person_id, position)
);
