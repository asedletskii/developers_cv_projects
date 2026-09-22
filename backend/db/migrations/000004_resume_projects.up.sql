CREATE TABLE projects (
    id uuid PRIMARY KEY,
    person_id uuid NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    description text,
    role text,
    link text,
    git_link text,
    UNIQUE (person_id, position)
);

CREATE TABLE project_screenshots (
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    position integer NOT NULL,
    path text NOT NULL,
    PRIMARY KEY (project_id, position)
);

CREATE TABLE project_technologies (
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    grade text,
    PRIMARY KEY (project_id, position)
);

CREATE TABLE project_contributors (
    id uuid PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    surname text NOT NULL,
    role text NOT NULL,
    phone text,
    telegram text,
    email text,
    UNIQUE (project_id, position)
);
