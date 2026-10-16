DROP TABLE IF EXISTS resume_link_events;

CREATE TABLE resume_link_views (
    id uuid PRIMARY KEY,
    link_id uuid NOT NULL REFERENCES resume_links(id) ON DELETE CASCADE,
    visitor_id uuid NOT NULL,
    ip_hash bytea NOT NULL,
    device text NOT NULL,
    browser text NOT NULL,
    os text NOT NULL,
    referrer_host text,
    traffic text NOT NULL DEFAULT 'human',
    traffic_reason text,
    layout jsonb NOT NULL,
    started_at timestamptz NOT NULL,
    confirmed_at timestamptz,
    last_seen_at timestamptz,
    active_ms bigint NOT NULL DEFAULT 0,
    max_section text,
    exit_section text,
    ended_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT resume_link_views_device_check CHECK (device IN ('desktop', 'mobile', 'tablet', 'unknown')),
    CONSTRAINT resume_link_views_browser_check CHECK (browser IN ('chrome', 'safari', 'firefox', 'edge', 'yandex', 'other')),
    CONSTRAINT resume_link_views_os_check CHECK (os IN ('windows', 'macos', 'linux', 'ios', 'android', 'other')),
    CONSTRAINT resume_link_views_traffic_check CHECK (traffic IN ('human', 'bot', 'owner', 'suspect')),
    CONSTRAINT resume_link_views_active_check CHECK (active_ms >= 0),
    CONSTRAINT resume_link_views_layout_check CHECK (jsonb_typeof(layout) = 'object')
);

CREATE INDEX idx_resume_link_views_link_started
    ON resume_link_views (link_id, started_at DESC);
CREATE INDEX idx_resume_link_views_link_visitor_started
    ON resume_link_views (link_id, visitor_id, started_at);
CREATE INDEX idx_resume_link_views_link_traffic_started
    ON resume_link_views (link_id, traffic, started_at);

CREATE TABLE resume_link_events (
    id uuid PRIMARY KEY,
    view_id uuid NOT NULL REFERENCES resume_link_views(id) ON DELETE CASCADE,
    link_id uuid NOT NULL REFERENCES resume_links(id) ON DELETE CASCADE,
    client_event_id uuid,
    event_type text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}',
    occurred_at timestamptz NOT NULL,
    CONSTRAINT resume_link_events_type_check CHECK (event_type IN (
        'contact_view',
        'section_view',
        'download',
        'tech_filter',
        'external_click',
        'print',
        'contact_copy',
        'contact_click'
    )),
    CONSTRAINT resume_link_events_payload_check CHECK (jsonb_typeof(payload) = 'object')
);

CREATE UNIQUE INDEX ux_resume_link_events_client_id
    ON resume_link_events (view_id, client_event_id)
    WHERE client_event_id IS NOT NULL;
CREATE INDEX idx_resume_link_events_view_id
    ON resume_link_events (view_id);
CREATE INDEX idx_resume_link_events_link_occurred_at
    ON resume_link_events (link_id, occurred_at);
