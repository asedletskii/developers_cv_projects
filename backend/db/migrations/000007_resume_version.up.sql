ALTER TABLE resumes
    ADD COLUMN version bigint NOT NULL DEFAULT 1,
    ADD CONSTRAINT resumes_version_check CHECK (version > 0);
