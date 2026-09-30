DROP INDEX IF EXISTS ux_resume_links_resume_id_name;

ALTER TABLE resume_links
    DROP CONSTRAINT IF EXISTS resume_links_name_length_check;
