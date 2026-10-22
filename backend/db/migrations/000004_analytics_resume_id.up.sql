ALTER TABLE analytics.analytics_views
    ADD COLUMN resume_id uuid;

UPDATE analytics.analytics_views AS view
SET resume_id = link.resume_id
FROM resume_link.resume_links AS link
WHERE link.id = view.link_id;

ALTER TABLE analytics.analytics_views
    ALTER COLUMN resume_id SET NOT NULL;

CREATE INDEX idx_analytics_views_resume_started
    ON analytics.analytics_views (resume_id, started_at DESC);
CREATE INDEX idx_analytics_views_resume_visitor_started
    ON analytics.analytics_views (resume_id, visitor_id, started_at);

ALTER TABLE analytics.analytics_views
    DROP CONSTRAINT analytics_views_link_id_fkey;
ALTER TABLE analytics.analytics_events
    DROP CONSTRAINT analytics_events_link_id_fkey;
ALTER TABLE analytics.analytics_clicks
    DROP CONSTRAINT analytics_clicks_link_id_fkey;
ALTER TABLE analytics.analytics_visitors
    DROP CONSTRAINT analytics_visitors_link_id_fkey;
