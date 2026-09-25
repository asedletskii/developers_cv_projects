CREATE TABLE resume_roots (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title text NOT NULL,
    language text NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    CONSTRAINT resume_roots_title_check CHECK (length(btrim(title)) > 0),
    CONSTRAINT resume_roots_user_language_unique UNIQUE (user_id, language)
);
CREATE INDEX idx_resume_roots_user_id ON resume_roots(user_id);

CREATE TABLE resume_presentations (
    resume_id uuid PRIMARY KEY REFERENCES resume_roots(id) ON DELETE CASCADE,
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
    CONSTRAINT resume_presentations_availability_status_check CHECK (availability_status IS NULL OR availability_status IN ('open', 'considering', 'busy')),
    CONSTRAINT resume_presentations_resume_link_check CHECK (resume_link IS NULL OR resume_link ~ '^/[^/].*$'),
    CONSTRAINT resume_presentations_theme_check CHECK (
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
    CONSTRAINT resume_presentation_sections_id_check CHECK (section_id IN ('about', 'techStack', 'experience', 'skills', 'education', 'courses', 'projects', 'activities', 'contacts'))
);

CREATE TABLE resume_profiles (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL UNIQUE REFERENCES resume_roots(id) ON DELETE CASCADE,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    position text NOT NULL,
    photo_path text NOT NULL,
    date_of_birth date NOT NULL,
    location text NOT NULL,
    relocation boolean NOT NULL,
    CONSTRAINT resume_profiles_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_profiles_surname_check CHECK (length(btrim(surname)) > 0),
    CONSTRAINT resume_profiles_position_check CHECK (length(btrim(position)) > 0),
    CONSTRAINT resume_profiles_photo_path_check CHECK (photo_path ~ '^/[^/].*$'),
    CONSTRAINT resume_profiles_location_check CHECK (length(btrim(location)) > 0)
);

CREATE TABLE resume_profile_work_formats (
    profile_id uuid NOT NULL REFERENCES resume_profiles(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (profile_id, value),
    UNIQUE (profile_id, position),
    CONSTRAINT resume_profile_work_formats_value_check CHECK (value IN ('fullTime', 'partTime', 'office', 'remote', 'hybrid')),
    CONSTRAINT resume_profile_work_formats_position_check CHECK (position >= 0)
);

CREATE TABLE resume_contacts (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL UNIQUE REFERENCES resume_roots(id) ON DELETE CASCADE,
    phone text,
    telegram text,
    email text,
    CONSTRAINT resume_contacts_presence_check CHECK (phone IS NOT NULL OR telegram IS NOT NULL OR email IS NOT NULL)
);

CREATE TABLE resume_contact_hidden (
    contacts_id uuid NOT NULL REFERENCES resume_contacts(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (contacts_id, value),
    UNIQUE (contacts_id, position),
    CONSTRAINT resume_contact_hidden_value_check CHECK (value IN ('telegram', 'email', 'phone')),
    CONSTRAINT resume_contact_hidden_position_check CHECK (position >= 0)
);

CREATE TABLE resume_contacts_text (
    contacts_id uuid PRIMARY KEY REFERENCES resume_contacts(id) ON DELETE CASCADE,
    title text,
    subtitle text,
    CONSTRAINT resume_contacts_text_presence_check CHECK (title IS NOT NULL OR subtitle IS NOT NULL)
);

CREATE TABLE resume_skills (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL UNIQUE REFERENCES resume_roots(id) ON DELETE CASCADE
);

CREATE TABLE resume_skill_hard (
    skills_id uuid NOT NULL REFERENCES resume_skills(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (skills_id, value),
    UNIQUE (skills_id, position),
    CONSTRAINT resume_skill_hard_value_check CHECK (length(btrim(value)) > 0),
    CONSTRAINT resume_skill_hard_position_check CHECK (position >= 0)
);

CREATE TABLE resume_skill_soft (
    skills_id uuid NOT NULL REFERENCES resume_skills(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (skills_id, value),
    UNIQUE (skills_id, position),
    CONSTRAINT resume_skill_soft_value_check CHECK (length(btrim(value)) > 0),
    CONSTRAINT resume_skill_soft_position_check CHECK (position >= 0)
);

CREATE TABLE resume_languages (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    language text NOT NULL,
    position integer NOT NULL,
    grade text,
    UNIQUE (resume_id, language),
    UNIQUE (resume_id, position),
    CONSTRAINT resume_languages_position_check CHECK (position >= 0),
    CONSTRAINT resume_languages_language_check CHECK (language IN ('ru', 'en', 'fr', 'de', 'kz', 'zh', 'es')),
    CONSTRAINT resume_languages_grade_check CHECK (grade IS NULL OR grade IN ('A1', 'A2', 'B1', 'B2', 'C1', 'C2'))
);

CREATE TABLE resume_technologies (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    name text NOT NULL,
    position integer NOT NULL,
    grade text,
    UNIQUE (resume_id, position),
    CONSTRAINT resume_technologies_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_technologies_position_check CHECK (position >= 0),
    CONSTRAINT resume_technologies_grade_check CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'))
);

CREATE TABLE resume_experiences (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    company text NOT NULL,
    position text NOT NULL,
    location text,
    start_date date NOT NULL,
    end_date date,
    description text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_experiences_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_experiences_company_check CHECK (length(btrim(company)) > 0),
    CONSTRAINT resume_experiences_position_check CHECK (length(btrim(position)) > 0),
    CONSTRAINT resume_experiences_dates_check CHECK (end_date IS NULL OR end_date >= start_date)
);
CREATE INDEX idx_resume_experiences_resume_id ON resume_experiences(resume_id);

CREATE TABLE resume_experience_achievements (
    experience_id uuid NOT NULL REFERENCES resume_experiences(id) ON DELETE CASCADE,
    position integer NOT NULL,
    value text NOT NULL,
    PRIMARY KEY (experience_id, position),
    CONSTRAINT resume_experience_achievements_position_check CHECK (position >= 0),
    CONSTRAINT resume_experience_achievements_value_check CHECK (length(btrim(value)) > 0)
);

CREATE TABLE resume_experience_technologies (
    experience_id uuid NOT NULL REFERENCES resume_experiences(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    grade text,
    PRIMARY KEY (experience_id, position),
    CONSTRAINT resume_experience_technologies_position_check CHECK (position >= 0),
    CONSTRAINT resume_experience_technologies_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_experience_technologies_grade_check CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'))
);

CREATE TABLE resume_educations (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
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
    CONSTRAINT resume_educations_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_educations_location_check CHECK (length(btrim(location)) > 0),
    CONSTRAINT resume_educations_university_check CHECK (length(btrim(university)) > 0),
    CONSTRAINT resume_educations_grade_check CHECK (grade BETWEEN 0 AND 6),
    CONSTRAINT resume_educations_program_name_check CHECK (length(btrim(program_name)) > 0),
    CONSTRAINT resume_educations_faculty_check CHECK (length(btrim(faculty)) > 0),
    CONSTRAINT resume_educations_course_check CHECK (course IS NULL OR course > 0)
);
CREATE INDEX idx_resume_educations_resume_id ON resume_educations(resume_id);

CREATE TABLE resume_courses (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    name text NOT NULL,
    organization text NOT NULL,
    end_date date NOT NULL,
    hours integer,
    certificate_link text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_courses_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_courses_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_courses_organization_check CHECK (length(btrim(organization)) > 0),
    CONSTRAINT resume_courses_hours_check CHECK (hours IS NULL OR hours > 0),
    CONSTRAINT resume_courses_certificate_link_check CHECK (certificate_link IS NULL OR certificate_link ~ '^https?://[^[:space:]]+$')
);
CREATE INDEX idx_resume_courses_resume_id ON resume_courses(resume_id);

CREATE TABLE resume_activities (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    type text NOT NULL,
    title text NOT NULL,
    description text,
    date date,
    link text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_activities_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_activities_type_check CHECK (type IN ('openSource', 'hackathon', 'speaker', 'volunteering', 'mentoring', 'community')),
    CONSTRAINT resume_activities_title_check CHECK (length(btrim(title)) > 0),
    CONSTRAINT resume_activities_link_check CHECK (link IS NULL OR link ~ '^https?://[^[:space:]]+$')
);
CREATE INDEX idx_resume_activities_resume_id ON resume_activities(resume_id);

CREATE TABLE resume_projects (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL REFERENCES resume_roots(id) ON DELETE CASCADE,
    sort_order integer NOT NULL,
    name text NOT NULL,
    description text,
    role text,
    link text,
    git_link text,
    UNIQUE (resume_id, sort_order),
    CONSTRAINT resume_projects_sort_order_check CHECK (sort_order >= 0),
    CONSTRAINT resume_projects_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_projects_link_check CHECK (link IS NULL OR link ~ '^https?://[^[:space:]]+$'),
    CONSTRAINT resume_projects_git_link_check CHECK (git_link IS NULL OR git_link ~ '^https?://[^[:space:]]+$')
);
CREATE INDEX idx_resume_projects_resume_id ON resume_projects(resume_id);

CREATE TABLE resume_project_screenshots (
    project_id uuid NOT NULL REFERENCES resume_projects(id) ON DELETE CASCADE,
    position integer NOT NULL,
    path text NOT NULL,
    PRIMARY KEY (project_id, position),
    CONSTRAINT resume_project_screenshots_position_check CHECK (position >= 0),
    CONSTRAINT resume_project_screenshots_path_check CHECK (path ~ '^/[^/].*$')
);

CREATE TABLE resume_project_technologies (
    project_id uuid NOT NULL REFERENCES resume_projects(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    grade text,
    PRIMARY KEY (project_id, position),
    CONSTRAINT resume_project_technologies_position_check CHECK (position >= 0),
    CONSTRAINT resume_project_technologies_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_project_technologies_grade_check CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'))
);

CREATE TABLE resume_project_contributors (
    id uuid PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES resume_projects(id) ON DELETE CASCADE,
    position integer NOT NULL,
    name text NOT NULL,
    surname text NOT NULL,
    role text NOT NULL,
    phone text,
    telegram text,
    email text,
    UNIQUE (project_id, position),
    CONSTRAINT resume_project_contributors_position_check CHECK (position >= 0),
    CONSTRAINT resume_project_contributors_name_check CHECK (length(btrim(name)) > 0),
    CONSTRAINT resume_project_contributors_surname_check CHECK (length(btrim(surname)) > 0),
    CONSTRAINT resume_project_contributors_role_check CHECK (length(btrim(role)) > 0)
);
CREATE INDEX idx_resume_project_contributors_project_id ON resume_project_contributors(project_id);
