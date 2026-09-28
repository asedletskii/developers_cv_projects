ALTER TABLE resumes
    DROP CONSTRAINT IF EXISTS resumes_version_check,
    DROP COLUMN IF EXISTS version;
