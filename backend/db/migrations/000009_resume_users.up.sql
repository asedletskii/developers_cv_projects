ALTER TABLE resumes
    ADD COLUMN user_id uuid;

UPDATE resumes
SET user_id = '10000000-0000-0000-0000-000000000001'
WHERE user_id IS NULL;

ALTER TABLE resumes
    ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE resumes
    ADD CONSTRAINT resumes_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX idx_resumes_user_id ON resumes(user_id);
