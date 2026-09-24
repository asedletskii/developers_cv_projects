package reader

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

type GetResumeReader struct{ pool *pgxpool.Pool }

func NewGetResumeReader(pool *pgxpool.Pool) *GetResumeReader { return &GetResumeReader{pool: pool} }

func (r *GetResumeReader) GetResume(
	ctx context.Context,
	query queries.GetResumeQuery,
) (queries.Resume, error) {
	var userID *uuid.UUID
	if query.UserID != nil {
		value := query.UserID.UUID()
		userID = &value
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
		workFormatJSON, contactsJSON, languagesJSON, experienceJSON, educationJSON, coursesJSON,
		activitiesJSON, hardSkillsJSON, softSkillsJSON, technologiesJSON, projectsJSON []byte
	var greetingJSON []byte

	row := r.pool.QueryRow(ctx, resumeQuerySQL, query.ID, userID, string(query.Lang))
	err := row.Scan(
		&p.ID, &p.Name, &p.Surname, &patronymic, &p.Position, &p.Title, &updatedAt, &resumeLink, &p.PhotoPath,
		&p.DateOfBirth, &p.Location, &p.Relocation, &workFormatJSON, &greetingJSON, &availabilityJSON, &themeJSON,
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
	if err := unmarshal(ensureArray(workFormatJSON), &p.WorkFormat); err != nil {
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

var _ queries.GetResumeReader = (*GetResumeReader)(nil)
