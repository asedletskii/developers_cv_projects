ALTER TABLE person_availability
    ADD CONSTRAINT person_availability_status_check CHECK (status IN ('open', 'considering', 'busy'));

ALTER TABLE person_themes
    ADD CONSTRAINT person_themes_color_check CHECK (
        (primary_color IS NULL OR primary_color ~ '^#[0-9A-Fa-f]{3}([0-9A-Fa-f]{3})?$')
        AND (secondary_color IS NULL OR secondary_color ~ '^#[0-9A-Fa-f]{3}([0-9A-Fa-f]{3})?$')
    );

ALTER TABLE person_sections
    ADD CONSTRAINT person_sections_id_check CHECK (section_id IN ('about', 'techStack', 'experience', 'skills', 'education', 'courses', 'projects', 'activities', 'contacts'));

ALTER TABLE person_work_formats
    ADD CONSTRAINT person_work_formats_value_check CHECK (value IN ('fullTime', 'partTime', 'office', 'remote', 'hybrid'));

ALTER TABLE person_hidden_contacts
    ADD CONSTRAINT person_hidden_contacts_value_check CHECK (value IN ('telegram', 'email', 'phone'));

ALTER TABLE person_hard_skills
    ADD CONSTRAINT person_hard_skills_value_check CHECK (value IN ('architecture', 'testing', 'cicd', 'performance', 'accessibility', 'codeReview', 'uiDesign', 'api'));

ALTER TABLE person_soft_skills
    ADD CONSTRAINT person_soft_skills_value_check CHECK (value IN ('communication', 'teamwork', 'problemSolving', 'timeManagement', 'mentoring', 'adaptability', 'ownership'));

ALTER TABLE person_languages
    ADD CONSTRAINT person_languages_value_check CHECK (language IN ('ru', 'en', 'fr', 'de', 'kz', 'zh', 'es')),
    ADD CONSTRAINT person_languages_grade_check CHECK (grade IS NULL OR grade IN ('A1', 'A2', 'B1', 'B2', 'C1', 'C2'));

ALTER TABLE person_technologies
    ADD CONSTRAINT person_technologies_grade_check CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'));

ALTER TABLE experiences
    ADD CONSTRAINT experiences_dates_check CHECK (end_date IS NULL OR end_date >= start_date),
    ADD CONSTRAINT experiences_position_check CHECK (position >= 0);

ALTER TABLE experience_technologies
    ADD CONSTRAINT experience_technologies_grade_check CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'));

ALTER TABLE educations
    ADD CONSTRAINT educations_grade_check CHECK (grade BETWEEN 0 AND 6),
    ADD CONSTRAINT educations_course_check CHECK (course IS NULL OR course >= 1);

ALTER TABLE courses
    ADD CONSTRAINT courses_hours_check CHECK (hours IS NULL OR hours > 0);

ALTER TABLE activities
    ADD CONSTRAINT activities_type_check CHECK (type IN ('openSource', 'hackathon', 'speaker', 'volunteering', 'mentoring', 'community'));

ALTER TABLE projects
    ADD CONSTRAINT projects_link_check CHECK (link IS NULL OR link ~ '^https?://[^[:space:]]+$'),
    ADD CONSTRAINT projects_git_link_check CHECK (git_link IS NULL OR git_link ~ '^https?://[^[:space:]]+$');

ALTER TABLE courses
    ADD CONSTRAINT courses_certificate_link_check CHECK (certificate_link IS NULL OR certificate_link ~ '^https?://[^[:space:]]+$');

ALTER TABLE activities
    ADD CONSTRAINT activities_link_check CHECK (link IS NULL OR link ~ '^https?://[^[:space:]]+$');

ALTER TABLE people
    ADD CONSTRAINT people_resume_link_check CHECK (resume_link IS NULL OR resume_link ~ '^/[^/].*$'),
    ADD CONSTRAINT people_photo_path_check CHECK (photo_path ~ '^/[^/].*$');

ALTER TABLE project_screenshots
    ADD CONSTRAINT project_screenshots_path_check CHECK (path ~ '^/[^/].*$');
