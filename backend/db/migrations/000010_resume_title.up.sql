ALTER TABLE resumes
    ADD COLUMN title text;

UPDATE resumes
SET title = position
WHERE title IS NULL;

ALTER TABLE resumes
    ALTER COLUMN title SET NOT NULL;
