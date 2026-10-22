-- Restore the cross-module foreign keys. Before running this rollback, verify
-- that the orphan checks are empty; otherwise the constraints cannot be added.
ALTER TABLE auth.user_identities
    ADD CONSTRAINT user_identities_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE;

ALTER TABLE auth.auth_sessions
    ADD CONSTRAINT auth_sessions_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE;

ALTER TABLE resume.resumes
    ADD CONSTRAINT resumes_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users.users(id) ON DELETE CASCADE;

ALTER TABLE resume_link.resume_links
    ADD CONSTRAINT resume_links_resume_id_fkey
    FOREIGN KEY (resume_id) REFERENCES resume.resumes(id) ON DELETE CASCADE;
