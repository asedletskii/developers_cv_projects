-- Phase 9: replace the last cross-module foreign keys with explicit cleanup
-- commands. Run the migration only after the cleanup code is deployed and the
-- orphan checks report zero rows.
ALTER TABLE resume_link.resume_links
    DROP CONSTRAINT resume_links_resume_id_fkey;

ALTER TABLE resume.resumes
    DROP CONSTRAINT resumes_user_id_fkey;

ALTER TABLE auth.user_identities
    DROP CONSTRAINT user_identities_user_id_fkey;

ALTER TABLE auth.auth_sessions
    DROP CONSTRAINT auth_sessions_user_id_fkey;
