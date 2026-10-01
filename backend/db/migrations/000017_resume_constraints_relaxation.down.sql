ALTER TABLE resume_skill_soft
    ADD CONSTRAINT resume_skill_soft_value_enum_check
    CHECK (value IN ('communication', 'teamwork', 'problemSolving', 'timeManagement', 'mentoring', 'adaptability', 'ownership'));

ALTER TABLE resume_skill_hard
    ADD CONSTRAINT resume_skill_hard_value_enum_check
    CHECK (value IN ('architecture', 'testing', 'cicd', 'performance', 'accessibility', 'codeReview', 'uiDesign', 'api'));

ALTER TABLE resume_languages
    DROP CONSTRAINT IF EXISTS resume_languages_language_check;

ALTER TABLE resume_languages
    ADD CONSTRAINT resume_languages_language_check
    CHECK (language IN ('ru', 'en', 'fr', 'de', 'kz', 'zh', 'es'));

ALTER TABLE resumes
    DROP CONSTRAINT IF EXISTS resumes_language_check;

ALTER TABLE resumes
    ADD CONSTRAINT resumes_language_check
    CHECK (language IN ('ru', 'en', 'fr', 'de', 'kz', 'zh', 'es'));
