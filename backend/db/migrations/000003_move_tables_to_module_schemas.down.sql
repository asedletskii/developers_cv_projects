ALTER TABLE analytics.analytics_visitors SET SCHEMA public;
ALTER TABLE analytics.analytics_clicks SET SCHEMA public;
ALTER TABLE analytics.analytics_events SET SCHEMA public;
ALTER TABLE analytics.analytics_views SET SCHEMA public;

ALTER TABLE resume_link.resume_links SET SCHEMA public;

ALTER TABLE resume.resume_experience_tasks SET SCHEMA public;
ALTER TABLE resume.resume_project_contributors SET SCHEMA public;
ALTER TABLE resume.resume_project_technologies SET SCHEMA public;
ALTER TABLE resume.resume_project_screenshots SET SCHEMA public;
ALTER TABLE resume.resume_projects SET SCHEMA public;
ALTER TABLE resume.resume_activities SET SCHEMA public;
ALTER TABLE resume.resume_courses SET SCHEMA public;
ALTER TABLE resume.resume_educations SET SCHEMA public;
ALTER TABLE resume.resume_experience_stack SET SCHEMA public;
ALTER TABLE resume.resume_experience_achievements SET SCHEMA public;
ALTER TABLE resume.resume_experiences SET SCHEMA public;
ALTER TABLE resume.resume_technologies SET SCHEMA public;
ALTER TABLE resume.resume_spoken_languages SET SCHEMA public;
ALTER TABLE resume.resume_skill_soft SET SCHEMA public;
ALTER TABLE resume.resume_skill_hard SET SCHEMA public;
ALTER TABLE resume.resume_skills SET SCHEMA public;
ALTER TABLE resume.resume_contacts_text SET SCHEMA public;
ALTER TABLE resume.resume_contact_hidden SET SCHEMA public;
ALTER TABLE resume.resume_contacts SET SCHEMA public;
ALTER TABLE resume.resume_profile_work_formats SET SCHEMA public;
ALTER TABLE resume.resume_profiles SET SCHEMA public;
ALTER TABLE resume.resume_overview_sections SET SCHEMA public;
ALTER TABLE resume.resume_overviews SET SCHEMA public;
ALTER TABLE resume.resumes SET SCHEMA public;

ALTER TABLE auth.auth_transactions SET SCHEMA public;
ALTER TABLE auth.auth_sessions SET SCHEMA public;
ALTER TABLE auth.user_identities SET SCHEMA public;

ALTER TABLE users.user_profile_localizations SET SCHEMA public;
ALTER TABLE users.user_profiles SET SCHEMA public;
ALTER TABLE users.users SET SCHEMA public;

DROP SCHEMA IF EXISTS analytics;
DROP SCHEMA IF EXISTS resume_link;
DROP SCHEMA IF EXISTS resume;
DROP SCHEMA IF EXISTS auth;
DROP SCHEMA IF EXISTS users;
