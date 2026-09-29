ALTER TABLE resume_links
    ADD CONSTRAINT resume_links_name_length_check
    CHECK (char_length(btrim(name)) BETWEEN 1 AND 128);

CREATE UNIQUE INDEX ux_resume_links_resume_id_name
    ON resume_links (resume_id, name);
