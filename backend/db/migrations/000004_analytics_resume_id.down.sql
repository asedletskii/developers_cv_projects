DROP INDEX analytics.idx_analytics_views_resume_visitor_started;
DROP INDEX analytics.idx_analytics_views_resume_started;

ALTER TABLE analytics.analytics_views
    DROP COLUMN resume_id;

ALTER TABLE analytics.analytics_views
    ADD CONSTRAINT analytics_views_link_id_fkey
    FOREIGN KEY (link_id) REFERENCES resume_link.resume_links(id) ON DELETE CASCADE;
ALTER TABLE analytics.analytics_events
    ADD CONSTRAINT analytics_events_link_id_fkey
    FOREIGN KEY (link_id) REFERENCES resume_link.resume_links(id) ON DELETE CASCADE;
ALTER TABLE analytics.analytics_clicks
    ADD CONSTRAINT analytics_clicks_link_id_fkey
    FOREIGN KEY (link_id) REFERENCES resume_link.resume_links(id) ON DELETE CASCADE;
ALTER TABLE analytics.analytics_visitors
    ADD CONSTRAINT analytics_visitors_link_id_fkey
    FOREIGN KEY (link_id) REFERENCES resume_link.resume_links(id) ON DELETE CASCADE;
