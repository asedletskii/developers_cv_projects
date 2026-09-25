package readers

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

type PublicResumeReader struct {
	pool *pgxpool.Pool
}

func NewPublicResumeReader(pool *pgxpool.Pool) *PublicResumeReader {
	return &PublicResumeReader{pool: pool}
}

func (r *PublicResumeReader) GetPublicResume(
	ctx context.Context,
	query queries.PublicResumeQuery,
) (queries.PublicResumeView, error) {
	if r == nil || r.pool == nil {
		return queries.PublicResumeView{}, fmt.Errorf("get public resume: reader dependencies are nil")
	}
	if query.ID == uuid.Nil {
		return queries.PublicResumeView{}, fmt.Errorf("invalid resume id")
	}

	var view queries.PublicResumeView
	var patronymic, resumeLink *string
	var updatedAt time.Time
	var availabilityJSON, themeJSON, sectionsJSON, hiddenJSON, contactsTextJSON,
		workFormatJSON, contactsJSON, languagesJSON, experienceJSON, educationJSON, coursesJSON,
		activitiesJSON, hardSkillsJSON, softSkillsJSON, technologiesJSON, projectsJSON []byte
	var greetingJSON []byte

	row := r.pool.QueryRow(ctx, publicResumeQuerySQL, query.ID)
	err := row.Scan(
		&view.ID, &view.Name, &view.Surname, &patronymic, &view.Position, &view.Title, &updatedAt, &resumeLink,
		&view.PhotoPath, &view.DateOfBirth, &view.Location, &view.Relocation, &workFormatJSON, &greetingJSON,
		&availabilityJSON, &themeJSON, &sectionsJSON, &hiddenJSON, &contactsTextJSON, &contactsJSON,
		&languagesJSON, &experienceJSON, &educationJSON, &coursesJSON, &activitiesJSON, &hardSkillsJSON,
		&softSkillsJSON, &technologiesJSON, &projectsJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return queries.PublicResumeView{}, queries.ErrResumeNotFound
		}
		return queries.PublicResumeView{}, fmt.Errorf("read public resume: %w", err)
	}

	view.Patronymic = patronymic
	view.UpdatedAt = &updatedAt
	view.ResumeLink = resumeLink
	if err := unmarshal(greetingJSON, &view.Greeting); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(workFormatJSON), &view.WorkFormat); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshalOptional(availabilityJSON, &view.Availability); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshalOptional(themeJSON, &view.Theme); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(sectionsJSON), &view.Sections); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(hiddenJSON), &view.HiddenContacts); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshalOptional(contactsTextJSON, &view.ContactsText); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(contactsJSON, &view.Contacts); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(languagesJSON), &view.Languages); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(experienceJSON), &view.Experience); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(educationJSON), &view.Education); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(coursesJSON), &view.Courses); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(activitiesJSON), &view.Activities); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(hardSkillsJSON), &view.HardSkills); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(softSkillsJSON), &view.SoftSkills); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(technologiesJSON), &view.Technologies); err != nil {
		return queries.PublicResumeView{}, err
	}
	if err := unmarshal(ensureArray(projectsJSON), &view.Projects); err != nil {
		return queries.PublicResumeView{}, err
	}
	return view, nil
}
