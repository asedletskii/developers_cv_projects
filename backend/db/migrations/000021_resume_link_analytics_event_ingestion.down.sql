DROP INDEX IF EXISTS ux_resume_link_events_section_view;
ALTER TABLE resume_link_views DROP COLUMN IF EXISTS last_state_seq;
