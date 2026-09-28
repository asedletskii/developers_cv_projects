ALTER TABLE resumes
    DROP CONSTRAINT IF EXISTS resumes_status_check;

ALTER TABLE resumes
    DROP COLUMN IF EXISTS status;
