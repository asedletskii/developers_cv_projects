package reader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type ResumeReader struct{ pool *pgxpool.Pool }

func NewResumeReader(pool *pgxpool.Pool) *ResumeReader { return &ResumeReader{pool: pool} }

func (r *ResumeReader) GetResume(
	ctx context.Context,
	query queries.GetResumeQuery,
) (queries.Resume, error) {
	if query.UserID.IsZero() {
		return queries.Resume{}, fmt.Errorf("invalid user id")
	}
	if query.ID == uuid.Nil {
		return queries.Resume{}, fmt.Errorf("invalid resume id")
	}
	if query.Lang != queries.LanguageRU && query.Lang != queries.LanguageEN {
		return queries.Resume{}, fmt.Errorf("unsupported language %q", query.Lang)
	}
	var p queries.Resume
	var patronymic, resumeLink *string
	var updatedAt time.Time
	var availabilityJSON, themeJSON, sectionsJSON, hiddenJSON, contactsTextJSON,
		contactsJSON, languagesJSON, experienceJSON, educationJSON, coursesJSON,
		activitiesJSON, hardSkillsJSON, softSkillsJSON, technologiesJSON, projectsJSON []byte
	var greetingJSON []byte

	row := r.pool.QueryRow(ctx, resumeQuerySQL, query.ID, query.UserID.UUID(), string(query.Lang))
	err := row.Scan(
		&p.ID, &p.Name, &p.Surname, &patronymic, &p.Position, &updatedAt, &resumeLink, &p.PhotoPath,
		&p.DateOfBirth, &p.Location, &p.Relocation, &greetingJSON, &availabilityJSON, &themeJSON,
		&sectionsJSON, &hiddenJSON, &contactsTextJSON, &contactsJSON, &languagesJSON, &experienceJSON,
		&educationJSON, &coursesJSON, &activitiesJSON, &hardSkillsJSON, &softSkillsJSON, &technologiesJSON, &projectsJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return queries.Resume{}, queries.ErrResumeNotFound
		}
		return queries.Resume{}, fmt.Errorf("read resume: %w", err)
	}

	p.Patronymic = patronymic
	p.UpdatedAt = &updatedAt
	p.ResumeLink = resumeLink
	if err := unmarshal(greetingJSON, &p.Greeting); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshalOptional(availabilityJSON, &p.Availability); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshalOptional(themeJSON, &p.Theme); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(sectionsJSON), &p.Sections); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(hiddenJSON), &p.HiddenContacts); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshalOptional(contactsTextJSON, &p.ContactsText); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(contactsJSON, &p.Contacts); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(languagesJSON), &p.Languages); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(experienceJSON), &p.Experience); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(educationJSON), &p.Education); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(coursesJSON), &p.Courses); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(activitiesJSON), &p.Activities); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(hardSkillsJSON), &p.HardSkills); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(softSkillsJSON), &p.SoftSkills); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(technologiesJSON), &p.Technologies); err != nil {
		return queries.Resume{}, err
	}
	if err := unmarshal(ensureArray(projectsJSON), &p.Projects); err != nil {
		return queries.Resume{}, err
	}
	return p, nil
}

func unmarshal(data []byte, target any) error {
	if len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode read model: %w", err)
	}
	return nil
}

func unmarshalOptional(data []byte, target any) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	return unmarshal(data, target)
}

func ensureArray(data []byte) []byte {
	if len(data) == 0 || string(data) == "null" {
		return []byte("[]")
	}
	return data
}

const resumeQuerySQL = `
SELECT
    p.id::text,
    p.name,
    p.surname,
    p.patronymic,
    p.position,
    p.updated_at,
    p.resume_link,
    p.photo_path,
    to_char(p.date_of_birth, 'YYYY-MM-DD'),
    p.location,
    p.relocation,
    jsonb_build_object(
        'first', CASE WHEN $3 = 'en' THEN COALESCE(NULLIF(g.first_en, ''), g.first_ru) ELSE g.first_ru END,
        'second', CASE WHEN $3 = 'en' THEN COALESCE(NULLIF(g.second_en, ''), g.second_ru) ELSE g.second_ru END,
        'secondHighlight', CASE WHEN $3 = 'en' THEN COALESCE(NULLIF(g.second_highlight_en, ''), g.second_highlight_ru) ELSE g.second_highlight_ru END,
        'third', CASE WHEN $3 = 'en' THEN COALESCE(NULLIF(g.third_en, ''), g.third_ru) ELSE g.third_ru END
    ),
    CASE WHEN a.resume_id IS NULL THEN NULL ELSE jsonb_build_object(
        'status', a.status,
        'from', CASE WHEN a.available_from IS NULL THEN NULL ELSE to_char(a.available_from, 'YYYY-MM-DD') END,
        'note', CASE WHEN $3 = 'en' THEN COALESCE(NULLIF(a.note_en, ''), a.note_ru) ELSE a.note_ru END
    ) END,
    CASE WHEN t.resume_id IS NULL THEN NULL ELSE jsonb_build_object(
        'primary', COALESCE(t.primary_color, ''),
        'secondary', COALESCE(t.secondary_color, '')
    ) END,
    COALESCE((SELECT jsonb_agg(jsonb_build_object('id', s.section_id, 'enabled', s.enabled, 'collapsed', s.collapsed) ORDER BY s.position) FROM resume_sections s WHERE s.resume_id = p.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(h.value ORDER BY h.position) FROM resume_hidden_contacts h WHERE h.resume_id = p.id), '[]'::jsonb),
    CASE WHEN ct.resume_id IS NULL THEN NULL ELSE jsonb_build_object(
        'title', CASE WHEN $3 = 'en' THEN COALESCE(NULLIF(ct.title_en, ''), ct.title_ru) ELSE ct.title_ru END,
        'subtitle', CASE WHEN $3 = 'en' THEN COALESCE(NULLIF(ct.subtitle_en, ''), ct.subtitle_ru) ELSE ct.subtitle_ru END
    ) END,
    jsonb_build_object('phone', c.phone, 'telegram', c.telegram, 'email', c.email),
    COALESCE((SELECT jsonb_agg(jsonb_build_object('language', l.language, 'grade', l.grade) ORDER BY l.position) FROM resume_languages l WHERE l.resume_id = p.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', e.id::text,
        'company', e.company,
        'position', e.job_position,
        'location', e.location,
        'startDate', to_char(e.start_date, 'YYYY-MM'),
        'endDate', CASE WHEN e.end_date IS NULL THEN NULL ELSE to_char(e.end_date, 'YYYY-MM') END,
        'description', e.description,
        'achievements', COALESCE((SELECT jsonb_agg(a.value ORDER BY a.position) FROM experience_achievements a WHERE a.experience_id = e.id), '[]'::jsonb),
        'stack', COALESCE((SELECT jsonb_agg(jsonb_build_object('name', x.name, 'grade', x.grade) ORDER BY x.position) FROM experience_technologies x WHERE x.experience_id = e.id), '[]'::jsonb)
    ) ORDER BY e.position) FROM experiences e WHERE e.resume_id = p.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', e.id::text, 'location', e.location, 'university', e.university, 'grade', e.grade,
        'programName', e.program_name, 'profileName', e.profile_name, 'faculty', e.faculty,
        'course', e.course, 'endDate', CASE WHEN e.end_date IS NULL THEN NULL ELSE to_char(e.end_date, 'YYYY') END
    ) ORDER BY e.position) FROM educations e WHERE e.resume_id = p.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', c1.id::text, 'name', c1.name, 'organization', c1.organization,
        'endDate', to_char(c1.end_date, 'YYYY-MM'), 'hours', c1.hours, 'certificateLink', c1.certificate_link
    ) ORDER BY c1.position) FROM courses c1 WHERE c1.resume_id = p.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', a1.id::text, 'type', a1.type, 'title', a1.title, 'description', a1.description,
        'date', CASE WHEN a1.date IS NULL THEN NULL ELSE to_char(a1.date, 'YYYY-MM') END, 'link', a1.link
    ) ORDER BY a1.position) FROM activities a1 WHERE a1.resume_id = p.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(h.value ORDER BY h.position) FROM resume_hard_skills h WHERE h.resume_id = p.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(s.value ORDER BY s.position) FROM resume_soft_skills s WHERE s.resume_id = p.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object('name', t1.name, 'grade', t1.grade) ORDER BY t1.position) FROM resume_technologies t1 WHERE t1.resume_id = p.id), '[]'::jsonb),
    COALESCE((SELECT jsonb_agg(jsonb_build_object(
        'id', pr.id::text, 'name', pr.name, 'description', pr.description,
        'stack', COALESCE((SELECT jsonb_agg(jsonb_build_object('name', pt.name, 'grade', pt.grade) ORDER BY pt.position) FROM project_technologies pt WHERE pt.project_id = pr.id), '[]'::jsonb),
        'link', pr.link
    ) ORDER BY pr.position) FROM projects pr WHERE pr.resume_id = p.id), '[]'::jsonb)
FROM resumes p
JOIN resume_greetings g ON g.resume_id = p.id
LEFT JOIN resume_availability a ON a.resume_id = p.id
LEFT JOIN resume_themes t ON t.resume_id = p.id
LEFT JOIN resume_contacts_text ct ON ct.resume_id = p.id
LEFT JOIN resume_contacts c ON c.resume_id = p.id
WHERE p.id = $1 AND p.user_id = $2
`
