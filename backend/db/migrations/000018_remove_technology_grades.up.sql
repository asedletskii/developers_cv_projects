ALTER TABLE resume_technologies
    DROP CONSTRAINT IF EXISTS resume_technologies_grade_check,
    DROP COLUMN IF EXISTS grade;

ALTER TABLE resume_experience_stack
    DROP CONSTRAINT IF EXISTS resume_experience_stack_grade_check,
    DROP COLUMN IF EXISTS grade;

ALTER TABLE resume_project_technologies
    DROP CONSTRAINT IF EXISTS resume_project_technologies_grade_check,
    DROP COLUMN IF EXISTS grade;
