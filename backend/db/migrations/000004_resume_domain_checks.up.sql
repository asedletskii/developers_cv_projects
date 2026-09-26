ALTER TABLE resumes
    ADD CONSTRAINT resumes_language_check
    CHECK (language IN ('ru', 'en', 'fr', 'de', 'kz', 'zh', 'es'));

ALTER TABLE resume_skill_hard
    ADD CONSTRAINT resume_skill_hard_value_enum_check
    CHECK (value IN ('architecture', 'testing', 'cicd', 'performance', 'accessibility', 'codeReview', 'uiDesign', 'api'));

ALTER TABLE resume_skill_soft
    ADD CONSTRAINT resume_skill_soft_value_enum_check
    CHECK (value IN ('communication', 'teamwork', 'problemSolving', 'timeManagement', 'mentoring', 'adaptability', 'ownership'));
