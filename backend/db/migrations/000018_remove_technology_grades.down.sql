ALTER TABLE resume_technologies
    ADD COLUMN grade text,
    ADD CONSTRAINT resume_technologies_grade_check
        CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'));

ALTER TABLE resume_experience_stack
    ADD COLUMN grade text,
    ADD CONSTRAINT resume_experience_stack_grade_check
        CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'));

ALTER TABLE resume_project_technologies
    ADD COLUMN grade text,
    ADD CONSTRAINT resume_project_technologies_grade_check
        CHECK (grade IS NULL OR grade IN ('intern', 'junior', 'middle', 'senior'));
