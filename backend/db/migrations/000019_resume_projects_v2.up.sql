CREATE TABLE resume_projects_v2 (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    name text NOT NULL,
    description text,
    role text,
    link text,
    git_link text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_projects_v2_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_projects_v2_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_projects_v2_link_check CHECK (link IS NULL OR link ~ '^https?://[^[:space:]]+$'),
    CONSTRAINT resume_projects_v2_git_link_check CHECK (git_link IS NULL OR git_link ~ '^https?://[^[:space:]]+$')
);

CREATE INDEX idx_resume_projects_v2_resume_id ON resume_projects_v2(resume_id);

CREATE TABLE resume_project_screenshots_v2 (
    project_id uuid NOT NULL REFERENCES resume_projects_v2(id) ON DELETE CASCADE,
    position integer NOT NULL,
    path text NOT NULL,
    PRIMARY KEY (project_id, position),
    CONSTRAINT resume_project_screenshots_v2_position_check CHECK (position >= 0),
    CONSTRAINT resume_project_screenshots_v2_path_check CHECK (path ~ '^/[^/].*$')
);

CREATE TABLE resume_project_technologies_v2 (
    project_id uuid NOT NULL REFERENCES resume_projects_v2(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    grade text,
    PRIMARY KEY (project_id, position),
    CONSTRAINT resume_project_technologies_v2_position_check CHECK (position >= 0),
    CONSTRAINT resume_project_technologies_v2_name_check CHECK (length(btrim(name)) > 0)
);

CREATE TABLE resume_project_contributors_v2 (
    id uuid PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES resume_projects_v2(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    surname text NOT NULL,
    role text NOT NULL,
    phone text,
    telegram text,
    email text,
    UNIQUE (project_id, position),
    CONSTRAINT resume_project_contributors_v2_position_check CHECK (position >= 0),
    CONSTRAINT resume_project_contributors_v2_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_project_contributors_v2_surname_check CHECK (length(btrim(surname)) > 0),
    CONSTRAINT resume_project_contributors_v2_role_check CHECK (length(btrim(role)) > 0)
);

CREATE INDEX idx_resume_project_contributors_v2_project_id ON resume_project_contributors_v2(project_id);

INSERT INTO resume_projects_v2 (id, resume_id, sort_order, name, description, role, link, git_link)
SELECT id, resume_id, position, name, description, role, link, git_link
FROM projects p
JOIN resume_roots r ON r.id = p.resume_id;

INSERT INTO resume_project_screenshots_v2 (project_id, position, path)
SELECT s.project_id, s.position, s.path
FROM project_screenshots s
JOIN resume_projects_v2 p ON p.id = s.project_id;

INSERT INTO resume_project_technologies_v2 (project_id, position, name, grade)
SELECT t.project_id, t.position, t.name, t.grade
FROM project_technologies t
JOIN resume_projects_v2 p ON p.id = t.project_id;

INSERT INTO resume_project_contributors_v2 (id, project_id, position, name, surname, role, phone, telegram, email)
SELECT c.id, c.project_id, c.position, c.name, c.surname, c.role, c.phone, c.telegram, c.email
FROM project_contributors c
JOIN resume_projects_v2 p ON p.id = c.project_id;
