CREATE TABLE resume_link_clicks (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    view_id      uuid NOT NULL REFERENCES resume_link_views(id) ON DELETE CASCADE,
    link_id      uuid NOT NULL REFERENCES resume_links(id) ON DELETE CASCADE,
    section      text NOT NULL,
    x_pct        smallint NOT NULL,
    y_pct        smallint NOT NULL,
    viewport_w   smallint NOT NULL,
    layout_hash  text NOT NULL,
    target       text,
    occurred_at  timestamptz NOT NULL,
    CONSTRAINT resume_link_clicks_x_pct_check CHECK (x_pct BETWEEN 0 AND 10000),
    CONSTRAINT resume_link_clicks_y_pct_check CHECK (y_pct BETWEEN 0 AND 10000),
    CONSTRAINT resume_link_clicks_viewport_w_check CHECK (viewport_w BETWEEN 1 AND 32767),
    CONSTRAINT resume_link_clicks_layout_hash_check CHECK (layout_hash ~ '^[0-9a-f]{16}$'),
    CONSTRAINT resume_link_clicks_target_check CHECK (target IS NULL OR char_length(target) <= 64)
);

CREATE INDEX idx_resume_link_clicks_link_section
    ON resume_link_clicks (link_id, section, layout_hash, occurred_at);
CREATE INDEX idx_resume_link_clicks_view_id
    ON resume_link_clicks (view_id);
