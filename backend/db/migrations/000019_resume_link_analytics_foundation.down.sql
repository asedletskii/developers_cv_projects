DROP TABLE IF EXISTS resume_link_events;
DROP TABLE IF EXISTS resume_link_views;

CREATE TABLE resume_link_events (
    id uuid PRIMARY KEY,
    link_id uuid NOT NULL REFERENCES resume_links(id) ON DELETE CASCADE,
    event_type text NOT NULL,
    occurred_at timestamptz NOT NULL,
    CONSTRAINT resume_link_events_type_check CHECK (event_type IN ('view', 'contact_view'))
);

CREATE INDEX idx_resume_link_events_link_id ON resume_link_events(link_id);
CREATE INDEX idx_resume_link_events_link_id_occurred_at ON resume_link_events(link_id, occurred_at);
