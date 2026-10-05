CREATE INDEX idx_resume_link_views_started_at_retention
    ON resume_link_views (started_at, id);

CREATE INDEX idx_resume_link_clicks_occurred_at_retention
    ON resume_link_clicks (occurred_at, id);
