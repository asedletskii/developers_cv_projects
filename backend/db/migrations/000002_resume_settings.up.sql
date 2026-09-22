CREATE TABLE person_sections (
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    section_id text NOT NULL,
    position integer NOT NULL,
    enabled boolean NOT NULL,
    collapsed boolean NOT NULL,
    PRIMARY KEY (person_id, section_id),
    UNIQUE (person_id, position)
);

CREATE TABLE person_work_formats (
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (person_id, value),
    UNIQUE (person_id, position)
);

CREATE TABLE person_hidden_contacts (
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (person_id, value),
    UNIQUE (person_id, position)
);

CREATE TABLE person_hard_skills (
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (person_id, value),
    UNIQUE (person_id, position)
);

CREATE TABLE person_soft_skills (
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (person_id, value),
    UNIQUE (person_id, position)
);

CREATE TABLE person_languages (
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    language text NOT NULL,
    grade text,
    position integer NOT NULL,
    PRIMARY KEY (person_id, language),
    UNIQUE (person_id, position)
);

CREATE TABLE person_technologies (
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    name text NOT NULL,
    grade text,
    position integer NOT NULL,
    PRIMARY KEY (person_id, name),
    UNIQUE (person_id, position)
);
