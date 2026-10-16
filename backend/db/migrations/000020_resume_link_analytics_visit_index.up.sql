CREATE INDEX idx_resume_link_views_ip_started
    ON resume_link_views (ip_hash, started_at DESC);
