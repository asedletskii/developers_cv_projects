ALTER TABLE resume_skill_soft
    DROP CONSTRAINT IF EXISTS resume_skill_soft_value_enum_check;

ALTER TABLE resume_skill_hard
    DROP CONSTRAINT IF EXISTS resume_skill_hard_value_enum_check;

ALTER TABLE resumes
    DROP CONSTRAINT IF EXISTS resumes_language_check;
