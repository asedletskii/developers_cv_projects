package readers

const publicResumeQuerySQL = `SELECT
    r.id::text,
    p.name,
    p.surname,
    p.patronymic,
    p.position,
    r.title,
    r.language,
    r.updated_at,
    rp.resume_link,
    p.photo_path,
    to_char(p.date_of_birth, 'YYYY-MM-DD'),
    p.location,
    p.relocation,
    COALESCE((SELECT jsonb_agg(w.value ORDER BY w.position) FROM resume_profile_work_formats w WHERE w.profile_id = p.id), '[]'::jsonb),
    jsonb_build_object(
        'first', rp.greeting_first,
        'second', rp.greeting_second,
        'secondHighlight', rp.greeting_highlight,
        'third', rp.greeting_third
    ),
    CASE WHEN rp.availability_status IS NULL THEN NULL ELSE jsonb_build_object(
        'status', rp.availability_status,
        'from', CASE WHEN rp.availability_from IS NULL THEN NULL ELSE to_char(rp.availability_from, 'YYYY-MM-DD') END,
        'note', rp.availability_note
    ) END,
    CASE WHEN rp.theme_primary IS NULL AND rp.theme_secondary IS NULL THEN NULL ELSE jsonb_build_object(
        'primary', COALESCE(rp.theme_primary, ''),
        'secondary', COALESCE(rp.theme_secondary, '')
    ) END,
    COALESCE((SELECT jsonb_agg(jsonb_build_object('id', s.section_id, 'enabled', s.enabled, 'collapsed', s.collapsed) ORDER BY s.position) FROM resume_overview_sections s WHERE s.resume_id = r.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(h.value ORDER BY h.position) FROM resume_contact_hidden h WHERE h.contacts_id = c.id), '[]'::jsonb),
    CASE WHEN ct.contacts_id IS NULL THEN NULL ELSE jsonb_build_object(
        'title', ct.title,
        'subtitle', ct.subtitle
    ) END,
    jsonb_build_object('phone', c.phone, 'telegram', c.telegram, 'email', c.email),
    COALESCE((SELECT jsonb_agg(jsonb_build_object('language', l.language, 'grade', l.grade) ORDER BY l.position) FROM resume_languages l WHERE l.resume_id = r.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', e.id::text,
        'company', e.company,
        'position', e.position,
        'location', e.location,
        'startDate', to_char(e.start_date, 'YYYY-MM'),
        'endDate', CASE WHEN e.end_date IS NULL THEN NULL ELSE to_char(e.end_date, 'YYYY-MM') END,
        'description', e.description,
        'achievements', COALESCE((SELECT jsonb_agg(a.value ORDER BY a.position) FROM resume_experience_achievements a WHERE a.experience_id = e.id), '[]'::jsonb),
        'stack', COALESCE((SELECT jsonb_agg(jsonb_build_object('name', x.name) ORDER BY x.position) FROM resume_experience_technologies x WHERE x.experience_id = e.id), '[]'::jsonb)
    ) ORDER BY e.sort_order) FROM resume_experiences e WHERE e.resume_id = r.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', e.id::text, 'location', e.location, 'university', e.university, 'grade', e.grade,
        'programName', e.program_name, 'profileName', e.profile_name, 'faculty', e.faculty,
        'course', e.course, 'endDate', CASE WHEN e.end_date IS NULL THEN NULL ELSE to_char(e.end_date, 'YYYY') END
    ) ORDER BY e.sort_order) FROM resume_educations e WHERE e.resume_id = r.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', c1.id::text, 'name', c1.name, 'organization', c1.organization,
        'endDate', to_char(c1.end_date, 'YYYY-MM'), 'hours', c1.hours, 'certificateLink', c1.certificate_link
    ) ORDER BY c1.sort_order) FROM resume_courses c1 WHERE c1.resume_id = r.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', a1.id::text, 'type', a1.type, 'title', a1.title, 'description', a1.description,
        'date', CASE WHEN a1.date IS NULL THEN NULL ELSE to_char(a1.date, 'YYYY-MM') END, 'link', a1.link
    ) ORDER BY a1.sort_order) FROM resume_activities a1 WHERE a1.resume_id = r.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(h.value ORDER BY h.position) FROM resume_skill_hard h WHERE h.skills_id = sk.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(s.value ORDER BY s.position) FROM resume_skill_soft s WHERE s.skills_id = sk.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object('name', t.name) ORDER BY t.position) FROM resume_technologies t WHERE t.resume_id = r.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', pr.id::text, 'name', pr.name, 'description', pr.description,
        'stack', COALESCE((SELECT jsonb_agg(jsonb_build_object('name', pt.name) ORDER BY pt.position) FROM resume_project_technologies pt WHERE pt.project_id = pr.id), '[]'::jsonb),
        'link', pr.link
    ) ORDER BY pr.sort_order) FROM resume_projects pr WHERE pr.resume_id = r.id), '[]'::jsonb)
FROM resumes r
JOIN resume_profiles p ON p.resume_id = r.id
JOIN resume_overviews rp ON rp.resume_id = r.id
LEFT JOIN resume_contacts c ON c.resume_id = r.id
LEFT JOIN resume_contacts_text ct ON ct.contacts_id = c.id
LEFT JOIN resume_skills sk ON sk.resume_id = r.id
WHERE r.id = $1
`
