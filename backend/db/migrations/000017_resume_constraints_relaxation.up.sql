ALTER TABLE resumes
    DROP CONSTRAINT IF EXISTS resumes_language_check;

ALTER TABLE resumes
    ADD CONSTRAINT resumes_language_check
    CHECK (language IN ('ru', 'en'));

ALTER TABLE resume_languages
    DROP CONSTRAINT IF EXISTS resume_languages_language_check;

ALTER TABLE resume_languages
    ADD CONSTRAINT resume_languages_language_check
    CHECK (language IN ('ru', 'en'));

ALTER TABLE resume_skill_hard
    DROP CONSTRAINT IF EXISTS resume_skill_hard_value_enum_check;

ALTER TABLE resume_skill_soft
    DROP CONSTRAINT IF EXISTS resume_skill_soft_value_enum_check;
