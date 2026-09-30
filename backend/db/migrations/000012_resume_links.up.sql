CREATE TABLE resume_links (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    name text NOT NULL,
    status text NOT NULL DEFAULT 'active',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT resume_links_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_links_status_check CHECK (status IN ('active', 'inactive')),
    CONSTRAINT resume_links_version_check CHECK (version > 0)
);

CREATE INDEX idx_resume_links_resume_id ON resume_links(resume_id);

CREATE TABLE resume_link_events (
    id uuid PRIMARY KEY,
    link_id uuid NOT NULL REFERENCES resume_links(id) ON DELETE CASCADE,
    event_type text NOT NULL,
    occurred_at timestamptz NOT NULL,
    CONSTRAINT resume_link_events_type_check CHECK (event_type IN ('view', 'contact_view'))
);

CREATE INDEX idx_resume_link_events_link_id ON resume_link_events(link_id);
CREATE INDEX idx_resume_link_events_link_id_occurred_at ON resume_link_events(link_id, occurred_at);
