ALTER TABLE resume_link.resume_links
    DROP CONSTRAINT resume_links_resume_id_fkey;

ALTER TABLE resume.resumes
    DROP CONSTRAINT resumes_user_id_fkey;

ALTER TABLE auth.user_identities
    DROP CONSTRAINT user_identities_user_id_fkey;

ALTER TABLE auth.auth_sessions
    DROP CONSTRAINT auth_sessions_user_id_fkey;
