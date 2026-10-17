CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id uuid PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE user_profiles (
    user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    photo_path text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT user_profiles_photo_path_check CHECK (photo_path IS NULL OR (octet_length(photo_path) <= 2048 AND photo_path ~ '^https://[^/@?#[:space:]]+([/?#].*)?$'))
);

CREATE TABLE user_profile_localizations (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language text NOT NULL,
    name text,
    surname text,
    patronymic text,
    position text,
    location text,
    CONSTRAINT user_profile_localizations_pkey PRIMARY KEY (user_id, language),
    CONSTRAINT user_profile_localizations_language_check CHECK (language IN ('ru', 'en')),
    CONSTRAINT user_profile_localizations_name_check CHECK (name IS NULL OR (length(btrim(name)) > 0 AND char_length(name) <= 100)),
    CONSTRAINT user_profile_localizations_surname_check CHECK (surname IS NULL OR (length(btrim(surname)) > 0 AND char_length(surname) <= 100)),
    CONSTRAINT user_profile_localizations_patronymic_length_check CHECK (patronymic IS NULL OR char_length(patronymic) <= 100),
    CONSTRAINT user_profile_localizations_position_length_check CHECK (position IS NULL OR char_length(position) <= 200),
    CONSTRAINT user_profile_localizations_location_length_check CHECK (location IS NULL OR char_length(location) <= 200)
);


CREATE TABLE resumes (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title text NOT NULL,
    language text NOT NULL,
    status text NOT NULL DEFAULT 'draft',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT resumes_title_check CHECK (char_length(btrim(title)) BETWEEN 1 AND 200),
    CONSTRAINT resumes_language_check CHECK (language IN ('ru', 'en')),
    CONSTRAINT resumes_status_check CHECK (status IN ('draft', 'opened', 'closed')),
    CONSTRAINT resumes_version_check CHECK (version > 0)
);
CREATE INDEX idx_resumes_user_id ON resumes(user_id);

CREATE TABLE resume_overviews (
    resume_id uuid PRIMARY KEY REFERENCES resumes(id) ON DELETE CASCADE,
    greeting_first text NOT NULL,
    greeting_second text NOT NULL,
    greeting_highlight text NOT NULL,
    greeting_third text NOT NULL,
    availability_status text,
    availability_from date,
    availability_note text,
    resume_link text,
    theme_primary text,
    theme_secondary text,
    CONSTRAINT resume_overviews_availability_status_check CHECK (availability_status IS NULL OR availability_status IN ('open', 'considering', 'busy')),
    CONSTRAINT resume_overviews_availability_fields_check CHECK (availability_status IS NOT NULL OR (availability_from IS NULL AND availability_note IS NULL)),
    CONSTRAINT resume_overviews_greeting_first_length_check CHECK (char_length(greeting_first) <= 200),
    CONSTRAINT resume_overviews_greeting_second_length_check CHECK (char_length(greeting_second) <= 200),
    CONSTRAINT resume_overviews_greeting_highlight_length_check CHECK (char_length(greeting_highlight) <= 120),
    CONSTRAINT resume_overviews_greeting_third_length_check CHECK (char_length(greeting_third) <= 200),
    CONSTRAINT resume_overviews_availability_note_length_check CHECK (availability_note IS NULL OR char_length(availability_note) <= 500),
    CONSTRAINT resume_overviews_resume_link_check CHECK (resume_link IS NULL OR (octet_length(resume_link) <= 2048 AND resume_link ~ '^https://[^/@?#[:space:]]+([/?#].*)?$')),
    CONSTRAINT resume_overviews_theme_check CHECK (
        (theme_primary IS NULL OR theme_primary ~ '^#[0-9A-Fa-f]{3}([0-9A-Fa-f]{3})?$')
        AND (theme_secondary IS NULL OR theme_secondary ~ '^#[0-9A-Fa-f]{3}([0-9A-Fa-f]{3})?$')
    )
);

CREATE TABLE resume_overview_sections (
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    section_id text NOT NULL,
    position integer NOT NULL,
    enabled boolean NOT NULL,
    collapsed boolean NOT NULL,
    PRIMARY KEY (resume_id, section_id),
    UNIQUE (resume_id, position),
    CONSTRAINT resume_overview_sections_position_check CHECK (position >= 0 AND position <= 8),
    CONSTRAINT resume_overview_sections_id_check CHECK (section_id IN ('about', 'techStack', 'experience', 'skills', 'education', 'courses', 'projects', 'activities', 'contacts'))
);

CREATE TABLE resume_profiles (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL UNIQUE REFERENCES resumes(id) ON DELETE CASCADE,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    position text NOT NULL,
    photo_url text NOT NULL,
    date_of_birth date NOT NULL,
    location text NOT NULL,
    relocation boolean NOT NULL,
    CONSTRAINT resume_profiles_name_check CHECK (length(btrim(name)) > 0 AND char_length(name) <= 100),
    CONSTRAINT resume_profiles_surname_check CHECK (length(btrim(surname)) > 0 AND char_length(surname) <= 100),
    CONSTRAINT resume_profiles_position_check CHECK (length(btrim(position)) > 0 AND char_length(position) <= 120),
    CONSTRAINT resume_profiles_photo_url_check CHECK ((octet_length(photo_url) <= 2048 AND photo_url ~ '^https://[^/@?#[:space:]]+([/?#].*)?$')),
    CONSTRAINT resume_profiles_location_check CHECK (length(btrim(location)) > 0 AND char_length(location) <= 200),
    CONSTRAINT resume_profiles_patronymic_length_check CHECK (patronymic IS NULL OR char_length(patronymic) <= 100)
);

CREATE TABLE resume_profile_work_formats (
    profile_id uuid NOT NULL REFERENCES resume_profiles(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (profile_id, value),
    UNIQUE (profile_id, position),
    CONSTRAINT resume_profile_work_formats_value_check CHECK (value IN ('fullTime', 'partTime', 'office', 'remote', 'hybrid')),
    CONSTRAINT resume_profile_work_formats_position_check CHECK (position >= 0 AND position <= 4)
);

CREATE TABLE resume_contacts (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL UNIQUE REFERENCES resumes(id) ON DELETE CASCADE,
    phone text,
    telegram text,
    email text,
    CONSTRAINT resume_contacts_presence_check CHECK (NULLIF(BTRIM(phone), '') IS NOT NULL OR NULLIF(BTRIM(telegram), '') IS NOT NULL OR NULLIF(BTRIM(email), '') IS NOT NULL),
    CONSTRAINT resume_contacts_phone_length_check CHECK (phone IS NULL OR char_length(phone) <= 32),
    CONSTRAINT resume_contacts_telegram_length_check CHECK (telegram IS NULL OR char_length(telegram) <= 64),
    CONSTRAINT resume_contacts_email_length_check CHECK (email IS NULL OR char_length(email) <= 254)
);

CREATE TABLE resume_contact_hidden (
    contacts_id uuid NOT NULL REFERENCES resume_contacts(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (contacts_id, value),
    UNIQUE (contacts_id, position),
    CONSTRAINT resume_contact_hidden_value_check CHECK (value IN ('telegram', 'email', 'phone')),
    CONSTRAINT resume_contact_hidden_position_check CHECK (position >= 0 AND position <= 2)
);

CREATE TABLE resume_contacts_text (
    contacts_id uuid PRIMARY KEY REFERENCES resume_contacts(id) ON DELETE CASCADE,
    title text,
    subtitle text,
    CONSTRAINT resume_contacts_text_presence_check CHECK (NULLIF(BTRIM(title), '') IS NOT NULL OR NULLIF(BTRIM(subtitle), '') IS NOT NULL),
    CONSTRAINT resume_contacts_text_title_length_check CHECK (title IS NULL OR char_length(title) <= 200),
    CONSTRAINT resume_contacts_text_subtitle_length_check CHECK (subtitle IS NULL OR char_length(subtitle) <= 200)
);

CREATE TABLE resume_skills (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL UNIQUE REFERENCES resumes(id) ON DELETE CASCADE
);

CREATE TABLE resume_skill_hard (
    skills_id uuid NOT NULL REFERENCES resume_skills(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (skills_id, value),
    UNIQUE (skills_id, position),
    CONSTRAINT resume_skill_hard_value_check CHECK (length(btrim(value)) > 0 AND char_length(value) <= 60),
    CONSTRAINT resume_skill_hard_position_check CHECK (position >= 0 AND position <= 49)
);

CREATE TABLE resume_skill_soft (
    skills_id uuid NOT NULL REFERENCES resume_skills(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (skills_id, value),
    UNIQUE (skills_id, position),
    CONSTRAINT resume_skill_soft_value_check CHECK (length(btrim(value)) > 0 AND char_length(value) <= 60),
    CONSTRAINT resume_skill_soft_position_check CHECK (position >= 0 AND position <= 49)
);

CREATE TABLE resume_spoken_languages (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    language text NOT NULL,
    position integer NOT NULL,
    grade text,
    UNIQUE (resume_id, language),
    UNIQUE (resume_id, position),
    CONSTRAINT resume_spoken_languages_position_check CHECK (position >= 0 AND position <= 6),
    CONSTRAINT resume_spoken_languages_language_check CHECK (language IN ('ru', 'en', 'fr', 'de', 'kk', 'zh', 'es')),
    CONSTRAINT resume_spoken_languages_grade_check CHECK (grade IS NULL OR grade IN ('A1', 'A2', 'B1', 'B2', 'C1', 'C2'))
);

CREATE TABLE resume_technologies (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    name text NOT NULL,
    position integer NOT NULL,
    UNIQUE (resume_id, position),
    CONSTRAINT resume_technologies_resume_id_name_key UNIQUE (resume_id, name),
    CONSTRAINT resume_technologies_name_check CHECK (length(btrim(name)) > 0 AND char_length(name) <= 60),
    CONSTRAINT resume_technologies_position_check CHECK (position >= 0 AND position <= 49)
);

CREATE TABLE resume_experiences (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    company text NOT NULL,
    position text NOT NULL,
    location text,
    start_date date NOT NULL,
    end_date date,
    description text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_experiences_sort_order_check CHECK (sort_order >= 0 AND sort_order <= 29),
    CONSTRAINT resume_experiences_company_check CHECK (length(btrim(company)) > 0 AND char_length(company) <= 120),
    CONSTRAINT resume_experiences_position_check CHECK (length(btrim(position)) > 0 AND char_length(position) <= 120),
    CONSTRAINT resume_experiences_dates_check CHECK (end_date IS NULL OR end_date >= start_date),
    CONSTRAINT resume_experiences_location_length_check CHECK (location IS NULL OR char_length(location) <= 200),
    CONSTRAINT resume_experiences_description_length_check CHECK (description IS NULL OR char_length(description) <= 2000)
);

CREATE TABLE resume_experience_achievements (
    experience_id uuid NOT NULL REFERENCES resume_experiences(id) ON DELETE CASCADE,
    position integer NOT NULL,
    value text NOT NULL,
    PRIMARY KEY (experience_id, position),
    CONSTRAINT resume_experience_achievements_position_check CHECK (position >= 0 AND position <= 29),
    CONSTRAINT resume_experience_achievements_value_check CHECK (length(btrim(value)) > 0 AND char_length(value) <= 500)
);

CREATE TABLE resume_experience_stack (
    experience_id uuid NOT NULL REFERENCES resume_experiences(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    PRIMARY KEY (experience_id, position),
    CONSTRAINT resume_experience_stack_experience_id_name_key UNIQUE (experience_id, name),
    CONSTRAINT resume_experience_stack_position_check CHECK (position >= 0 AND position <= 49),
    CONSTRAINT resume_experience_stack_name_check CHECK (length(btrim(name)) > 0 AND char_length(name) <= 60)
);

CREATE TABLE resume_educations (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    location text NOT NULL,
    university text NOT NULL,
    grade integer NOT NULL,
    program_name text NOT NULL,
    profile_name text,
    faculty text NOT NULL,
    course integer,
    end_date date,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_educations_sort_order_check CHECK (sort_order >= 0 AND sort_order <= 29),
    CONSTRAINT resume_educations_location_check CHECK (length(btrim(location)) > 0 AND char_length(location) <= 200),
    CONSTRAINT resume_educations_university_check CHECK (length(btrim(university)) > 0 AND char_length(university) <= 200),
    CONSTRAINT resume_educations_grade_check CHECK (grade BETWEEN 0 AND 6),
    CONSTRAINT resume_educations_program_name_check CHECK (length(btrim(program_name)) > 0 AND char_length(program_name) <= 200),
    CONSTRAINT resume_educations_faculty_check CHECK (length(btrim(faculty)) > 0 AND char_length(faculty) <= 200),
    CONSTRAINT resume_educations_course_check CHECK (course IS NULL OR course > 0),
    CONSTRAINT resume_educations_profile_name_length_check CHECK (profile_name IS NULL OR char_length(profile_name) <= 200)
);

CREATE TABLE resume_courses (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    name text NOT NULL,
    organization text NOT NULL,
    end_date date NOT NULL,
    hours integer,
    certificate_link text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_courses_sort_order_check CHECK (sort_order >= 0 AND sort_order <= 29),
    CONSTRAINT resume_courses_name_check CHECK (length(btrim(name)) > 0 AND char_length(name) <= 200),
    CONSTRAINT resume_courses_organization_check CHECK (length(btrim(organization)) > 0 AND char_length(organization) <= 200),
    CONSTRAINT resume_courses_hours_check CHECK (hours IS NULL OR hours > 0),
    CONSTRAINT resume_courses_certificate_link_check CHECK (certificate_link IS NULL OR (octet_length(certificate_link) <= 2048 AND certificate_link ~ '^https://[^/@?#[:space:]]+([/?#].*)?$'))
);

CREATE TABLE resume_activities (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    type text NOT NULL,
    title text NOT NULL,
    description text,
    date date,
    link text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_activities_sort_order_check CHECK (sort_order >= 0 AND sort_order <= 29),
    CONSTRAINT resume_activities_type_check CHECK (type IN ('openSource', 'hackathon', 'speaker', 'volunteering', 'mentoring', 'community')),
    CONSTRAINT resume_activities_title_check CHECK (length(btrim(title)) > 0 AND char_length(title) <= 200),
    CONSTRAINT resume_activities_link_check CHECK (link IS NULL OR (octet_length(link) <= 2048 AND link ~ '^https://[^/@?#[:space:]]+([/?#].*)?$')),
    CONSTRAINT resume_activities_description_length_check CHECK (description IS NULL OR char_length(description) <= 2000)
);

CREATE TABLE resume_projects (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    name text NOT NULL,
    description text,
    role text,
    link text,
    git_link text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_projects_sort_order_check CHECK (sort_order >= 0 AND sort_order <= 29),
    CONSTRAINT resume_projects_name_check CHECK (length(btrim(name)) > 0 AND char_length(name) <= 200),
    CONSTRAINT resume_projects_link_check CHECK (link IS NULL OR (octet_length(link) <= 2048 AND link ~ '^https://[^/@?#[:space:]]+([/?#].*)?$')),
    CONSTRAINT resume_projects_git_link_check CHECK (git_link IS NULL OR (octet_length(git_link) <= 2048 AND git_link ~ '^https://[^/@?#[:space:]]+([/?#].*)?$')),
    CONSTRAINT resume_projects_description_length_check CHECK (description IS NULL OR char_length(description) <= 2000),
    CONSTRAINT resume_projects_role_length_check CHECK (role IS NULL OR char_length(role) <= 200)
);

CREATE TABLE resume_project_screenshots (
    project_id uuid NOT NULL REFERENCES resume_projects(id) ON DELETE CASCADE,
    position integer NOT NULL,
    path text NOT NULL,
    PRIMARY KEY (project_id, position),
    CONSTRAINT resume_project_screenshots_position_check CHECK (position >= 0 AND position <= 9),
    CONSTRAINT resume_project_screenshots_path_check CHECK ((octet_length(path) <= 2048 AND path ~ '^https://[^/@?#[:space:]]+([/?#].*)?$'))
);

CREATE TABLE resume_project_technologies (
    project_id uuid NOT NULL REFERENCES resume_projects(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    PRIMARY KEY (project_id, position),
    CONSTRAINT resume_project_technologies_project_id_name_key UNIQUE (project_id, name),
    CONSTRAINT resume_project_technologies_position_check CHECK (position >= 0 AND position <= 49),
    CONSTRAINT resume_project_technologies_name_check CHECK (length(btrim(name)) > 0 AND char_length(name) <= 60)
);

CREATE TABLE resume_project_contributors (
    project_id uuid NOT NULL REFERENCES resume_projects(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    surname text NOT NULL,
    role text NOT NULL,
    phone text,
    telegram text,
    email text,
    PRIMARY KEY (project_id, position),
    CONSTRAINT resume_project_contributors_position_check CHECK (position >= 0 AND position <= 19),
    CONSTRAINT resume_project_contributors_name_check CHECK (length(btrim(name)) > 0 AND char_length(name) <= 100),
    CONSTRAINT resume_project_contributors_surname_check CHECK (length(btrim(surname)) > 0 AND char_length(surname) <= 100),
    CONSTRAINT resume_project_contributors_role_check CHECK (length(btrim(role)) > 0 AND char_length(role) <= 120),
    CONSTRAINT resume_project_contributors_contacts_check CHECK (NULLIF(BTRIM(phone), '') IS NOT NULL OR NULLIF(BTRIM(telegram), '') IS NOT NULL OR NULLIF(BTRIM(email), '') IS NOT NULL),
    CONSTRAINT resume_project_contributors_phone_length_check CHECK (phone IS NULL OR char_length(phone) <= 32),
    CONSTRAINT resume_project_contributors_telegram_length_check CHECK (telegram IS NULL OR char_length(telegram) <= 64),
    CONSTRAINT resume_project_contributors_email_length_check CHECK (email IS NULL OR char_length(email) <= 254)
);



CREATE TABLE resume_experience_tasks (
    experience_id uuid NOT NULL REFERENCES resume_experiences(id) ON DELETE CASCADE,
    position integer NOT NULL,
    value text NOT NULL,
    PRIMARY KEY (experience_id, position),
    CONSTRAINT resume_experience_tasks_position_check CHECK (position >= 0 AND position <= 29),
    CONSTRAINT resume_experience_tasks_value_check CHECK (length(btrim(value)) > 0 AND char_length(value) <= 500)
);

CREATE TABLE resume_links (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    name text NOT NULL,
    status text NOT NULL DEFAULT 'active',
    version bigint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT resume_links_name_length_check CHECK (char_length(btrim(name)) BETWEEN 1 AND 128),
    CONSTRAINT resume_links_status_check CHECK (status IN ('active', 'inactive')),
    CONSTRAINT resume_links_version_check CHECK (version > 0)
);
CREATE INDEX idx_resume_links_resume_id ON resume_links(resume_id);
CREATE UNIQUE INDEX ux_resume_links_resume_id_name ON resume_links (resume_id, name);

CREATE TABLE user_identities (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider text NOT NULL,
    provider_subject text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    last_login_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT user_identities_provider_check CHECK (length(provider) > 0),
    CONSTRAINT user_identities_subject_check CHECK (length(provider_subject) > 0),
    CONSTRAINT user_identities_timestamps_check CHECK (updated_at >= created_at AND last_login_at >= created_at),
    CONSTRAINT user_identities_provider_subject_unique UNIQUE (provider, provider_subject)
);

CREATE TABLE auth_sessions (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash bytea NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    last_seen_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    user_agent text,
    ip_address inet,
    CONSTRAINT auth_sessions_token_hash_unique UNIQUE (token_hash),
    CONSTRAINT auth_sessions_token_hash_length CHECK (octet_length(token_hash) = 32),
    CONSTRAINT auth_sessions_timestamps_check CHECK (last_seen_at >= created_at AND expires_at > created_at AND (revoked_at IS NULL OR revoked_at >= created_at))
);
CREATE INDEX auth_sessions_user_id_idx ON auth_sessions(user_id);
CREATE INDEX auth_sessions_expires_at_idx ON auth_sessions(expires_at);

CREATE TABLE auth_transactions (
    id uuid PRIMARY KEY,
    provider text NOT NULL,
    state_hash bytea NOT NULL,
    nonce_hash bytea NOT NULL,
    browser_binding_hash bytea NOT NULL,
    code_verifier text NOT NULL,
    redirect_uri text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    consumed_at timestamptz,
    CONSTRAINT auth_transactions_provider_check CHECK (length(provider) > 0),
    CONSTRAINT auth_transactions_state_hash_length CHECK (octet_length(state_hash) = 32),
    CONSTRAINT auth_transactions_nonce_hash_length CHECK (octet_length(nonce_hash) = 32),
    CONSTRAINT auth_transactions_browser_binding_hash_length CHECK (octet_length(browser_binding_hash) = 32),
    CONSTRAINT auth_transactions_code_verifier_check CHECK (length(code_verifier) > 0),
    CONSTRAINT auth_transactions_redirect_uri_check CHECK (length(redirect_uri) > 0),
    CONSTRAINT auth_transactions_timestamps_check CHECK (expires_at > created_at AND (consumed_at IS NULL OR consumed_at >= created_at)),
    CONSTRAINT auth_transactions_state_hash_unique UNIQUE (state_hash)
);
CREATE INDEX auth_transactions_expires_at_idx ON auth_transactions(expires_at);

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
    last_state_seq bigint NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT resume_link_views_device_check CHECK (device IN ('desktop', 'mobile', 'tablet', 'unknown')),
    CONSTRAINT resume_link_views_browser_check CHECK (browser IN ('chrome', 'safari', 'firefox', 'edge', 'yandex', 'other')),
    CONSTRAINT resume_link_views_os_check CHECK (os IN ('windows', 'macos', 'linux', 'ios', 'android', 'other')),
    CONSTRAINT resume_link_views_traffic_check CHECK (traffic IN ('human', 'bot', 'owner', 'suspect')),
    CONSTRAINT resume_link_views_referrer_host_length_check CHECK (referrer_host IS NULL OR octet_length(referrer_host) <= 253),
    CONSTRAINT resume_link_views_traffic_reason_length_check CHECK (traffic_reason IS NULL OR octet_length(traffic_reason) <= 100),
    CONSTRAINT resume_link_views_active_check CHECK (active_ms BETWEEN 0 AND 43230000),
    CONSTRAINT resume_link_views_max_section_check CHECK (max_section IS NULL OR max_section = '' OR max_section IN ('about', 'experience', 'skills', 'education', 'courses', 'projects', 'activities', 'contacts')),
    CONSTRAINT resume_link_views_exit_section_check CHECK (exit_section IS NULL OR exit_section = '' OR exit_section IN ('about', 'experience', 'skills', 'education', 'courses', 'projects', 'activities', 'contacts')),
    CONSTRAINT resume_link_views_layout_check CHECK (jsonb_typeof(layout) = 'object')
);
CREATE INDEX idx_resume_link_views_link_started ON resume_link_views (link_id, started_at DESC);
CREATE INDEX idx_resume_link_views_link_visitor_started ON resume_link_views (link_id, visitor_id, started_at);
CREATE INDEX idx_resume_link_views_link_traffic_started ON resume_link_views (link_id, traffic, started_at);
CREATE INDEX idx_resume_link_views_ip_started ON resume_link_views (ip_hash, started_at DESC);
CREATE INDEX idx_resume_link_views_started_at_retention ON resume_link_views (started_at, id);

CREATE TABLE resume_link_events (
    id uuid PRIMARY KEY,
    view_id uuid NOT NULL REFERENCES resume_link_views(id) ON DELETE CASCADE,
    link_id uuid NOT NULL REFERENCES resume_links(id) ON DELETE CASCADE,
    client_event_id uuid,
    event_type text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}',
    occurred_at timestamptz NOT NULL,
    CONSTRAINT resume_link_events_type_check CHECK (event_type IN ('contact_view', 'section_view', 'download', 'tech_filter', 'external_click', 'print', 'contact_copy', 'contact_click')),
    CONSTRAINT resume_link_events_payload_check CHECK (jsonb_typeof(payload) = 'object')
);
CREATE UNIQUE INDEX ux_resume_link_events_client_id ON resume_link_events (view_id, client_event_id) WHERE client_event_id IS NOT NULL;
CREATE UNIQUE INDEX ux_resume_link_events_section_view ON resume_link_events (view_id, (payload->>'section')) WHERE event_type = 'section_view';
CREATE INDEX idx_resume_link_events_view_id ON resume_link_events (view_id);
CREATE INDEX idx_resume_link_events_link_occurred_at ON resume_link_events (link_id, occurred_at);

CREATE TABLE resume_link_clicks (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    view_id uuid NOT NULL REFERENCES resume_link_views(id) ON DELETE CASCADE,
    link_id uuid NOT NULL REFERENCES resume_links(id) ON DELETE CASCADE,
    section text NOT NULL,
    x_pct smallint NOT NULL,
    y_pct smallint NOT NULL,
    viewport_w smallint NOT NULL,
    layout_hash text NOT NULL,
    target text,
    occurred_at timestamptz NOT NULL,
    CONSTRAINT resume_link_clicks_x_pct_check CHECK (x_pct BETWEEN 0 AND 10000),
    CONSTRAINT resume_link_clicks_y_pct_check CHECK (y_pct BETWEEN 0 AND 10000),
    CONSTRAINT resume_link_clicks_viewport_w_check CHECK (viewport_w BETWEEN 1 AND 32767),
    CONSTRAINT resume_link_clicks_layout_hash_check CHECK (layout_hash ~ '^[0-9a-f]{16}$'),
    CONSTRAINT resume_link_clicks_target_check CHECK (target IS NULL OR char_length(target) <= 64)
);
CREATE INDEX idx_resume_link_clicks_link_section ON resume_link_clicks (link_id, section, layout_hash, occurred_at);
CREATE INDEX idx_resume_link_clicks_view_id ON resume_link_clicks (view_id);
CREATE INDEX idx_resume_link_clicks_occurred_at_retention ON resume_link_clicks (occurred_at, id);

CREATE TABLE resume_link_visitors (
    link_id uuid NOT NULL REFERENCES resume_links(id) ON DELETE CASCADE,
    visitor_id uuid NOT NULL,
    first_seen_at timestamptz NOT NULL,
    first_view_id uuid NOT NULL,
    ordinal bigint,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT resume_link_visitors_ordinal_check CHECK (ordinal IS NULL OR ordinal > 0),
    CONSTRAINT resume_link_visitors_link_visitor_unique UNIQUE (link_id, visitor_id),
    CONSTRAINT resume_link_visitors_link_ordinal_unique UNIQUE (link_id, ordinal) DEFERRABLE INITIALLY DEFERRED
);
CREATE INDEX idx_resume_link_visitors_link_first_seen ON resume_link_visitors (link_id, first_seen_at, first_view_id);
