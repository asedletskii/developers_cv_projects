CREATE SCHEMA IF NOT EXISTS users;
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS resume;
CREATE SCHEMA IF NOT EXISTS resume_link;
CREATE SCHEMA IF NOT EXISTS analytics;

ALTER TABLE public.users SET SCHEMA users;
ALTER TABLE public.user_profiles SET SCHEMA users;
ALTER TABLE public.user_profile_localizations SET SCHEMA users;

ALTER TABLE public.user_identities SET SCHEMA auth;
ALTER TABLE public.auth_sessions SET SCHEMA auth;
ALTER TABLE public.auth_transactions SET SCHEMA auth;

ALTER TABLE public.resumes SET SCHEMA resume;
ALTER TABLE public.resume_overviews SET SCHEMA resume;
ALTER TABLE public.resume_overview_sections SET SCHEMA resume;
ALTER TABLE public.resume_profiles SET SCHEMA resume;
ALTER TABLE public.resume_profile_work_formats SET SCHEMA resume;
ALTER TABLE public.resume_contacts SET SCHEMA resume;
ALTER TABLE public.resume_contact_hidden SET SCHEMA resume;
ALTER TABLE public.resume_contacts_text SET SCHEMA resume;
ALTER TABLE public.resume_skills SET SCHEMA resume;
ALTER TABLE public.resume_skill_hard SET SCHEMA resume;
ALTER TABLE public.resume_skill_soft SET SCHEMA resume;
ALTER TABLE public.resume_spoken_languages SET SCHEMA resume;
ALTER TABLE public.resume_technologies SET SCHEMA resume;
ALTER TABLE public.resume_experiences SET SCHEMA resume;
ALTER TABLE public.resume_experience_achievements SET SCHEMA resume;
ALTER TABLE public.resume_experience_stack SET SCHEMA resume;
ALTER TABLE public.resume_educations SET SCHEMA resume;
ALTER TABLE public.resume_courses SET SCHEMA resume;
ALTER TABLE public.resume_activities SET SCHEMA resume;
ALTER TABLE public.resume_projects SET SCHEMA resume;
ALTER TABLE public.resume_project_screenshots SET SCHEMA resume;
ALTER TABLE public.resume_project_technologies SET SCHEMA resume;
ALTER TABLE public.resume_project_contributors SET SCHEMA resume;
ALTER TABLE public.resume_experience_tasks SET SCHEMA resume;

ALTER TABLE public.resume_links SET SCHEMA resume_link;

ALTER TABLE public.analytics_views SET SCHEMA analytics;
ALTER TABLE public.analytics_events SET SCHEMA analytics;
ALTER TABLE public.analytics_clicks SET SCHEMA analytics;
ALTER TABLE public.analytics_visitors SET SCHEMA analytics;
