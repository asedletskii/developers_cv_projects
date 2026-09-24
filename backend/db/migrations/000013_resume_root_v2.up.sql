CREATE TABLE resume_roots (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT resume_roots_title_check CHECK (length(btrim(title)) > 0)
);

CREATE INDEX idx_resume_roots_user_id ON resume_roots(user_id);

CREATE TABLE resume_presentations (
    resume_id uuid PRIMARY KEY REFERENCES resume_roots(id) ON DELETE CASCADE,
    greeting_first_ru text NOT NULL,
    greeting_first_en text NOT NULL,
    greeting_second_ru text NOT NULL,
    greeting_second_en text NOT NULL,
    greeting_highlight_ru text NOT NULL,
    greeting_highlight_en text NOT NULL,
    greeting_third_ru text NOT NULL,
    greeting_third_en text NOT NULL,
    availability_status text,
    availability_from date,
    availability_note_ru text,
    availability_note_en text,
    resume_link text,
    theme_primary text,
    theme_secondary text,
    CONSTRAINT resume_presentations_availability_status_check
        CHECK (availability_status IS NULL OR availability_status IN ('open', 'considering', 'busy')),
    CONSTRAINT resume_presentations_resume_link_check
        CHECK (resume_link IS NULL OR resume_link ~ '^/[^/].*$'),
    CONSTRAINT resume_presentations_theme_check
        CHECK (
            (theme_primary IS NULL OR theme_primary ~ '^#[0-9A-Fa-f]{3}([0-9A-Fa-f]{3})?$')
            AND (theme_secondary IS NULL OR theme_secondary ~ '^#[0-9A-Fa-f]{3}([0-9A-Fa-f]{3})?$')
        )
);

CREATE TABLE resume_presentation_sections (
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    section_id text NOT NULL,
    position integer NOT NULL,
    enabled boolean NOT NULL,
    collapsed boolean NOT NULL,
    PRIMARY KEY (resume_id, section_id),
    UNIQUE (resume_id, position),
    CONSTRAINT resume_presentation_sections_position_check CHECK (position >= 0),
    CONSTRAINT resume_presentation_sections_id_check CHECK (
        section_id IN ('about', 'techStack', 'experience', 'skills', 'education', 'courses', 'projects', 'activities', 'contacts')
    )
);

CREATE INDEX idx_resume_presentation_sections_resume_position
    ON resume_presentation_sections(resume_id, position);

INSERT INTO resume_roots (id, user_id, title, created_at, updated_at)
SELECT id, user_id, title, updated_at, updated_at
FROM resumes;

INSERT INTO resume_presentations (
    resume_id,
    greeting_first_ru, greeting_first_en,
    greeting_second_ru, greeting_second_en,
    greeting_highlight_ru, greeting_highlight_en,
    greeting_third_ru, greeting_third_en,
    availability_status, availability_from, availability_note_ru, availability_note_en,
    resume_link, theme_primary, theme_secondary
)
SELECT
    r.id,
    g.first_ru, g.first_en,
    g.second_ru, g.second_en,
    g.second_highlight_ru, g.second_highlight_en,
    g.third_ru, g.third_en,
    a.status, a.available_from, a.note_ru, a.note_en,
    r.resume_link, t.primary_color, t.secondary_color
FROM resumes r
JOIN resume_greetings g ON g.resume_id = r.id
LEFT JOIN resume_availability a ON a.resume_id = r.id
LEFT JOIN resume_themes t ON t.resume_id = r.id;

INSERT INTO resume_presentation_sections (resume_id, section_id, position, enabled, collapsed)
SELECT resume_id, section_id, position, enabled, collapsed
FROM resume_sections;
