DROP INDEX IF EXISTS idx_resumes_user_id;

ALTER TABLE resumes
    DROP CONSTRAINT IF EXISTS resumes_user_id_fkey;

ALTER TABLE resumes
    DROP COLUMN IF EXISTS user_id;
