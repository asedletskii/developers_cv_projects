ALTER TABLE resumes
    ADD COLUMN status TEXT NOT NULL DEFAULT 'draft';

ALTER TABLE resumes
    ADD CONSTRAINT resumes_status_check
    CHECK (status IN ('draft', 'opened', 'closed'));
