ALTER TABLE resume_link_views
    ADD COLUMN last_state_seq bigint NOT NULL DEFAULT 0;

CREATE UNIQUE INDEX ux_resume_link_events_section_view
    ON resume_link_events (
        view_id,
        (payload->>'section')
    )
    WHERE event_type = 'section_view';
