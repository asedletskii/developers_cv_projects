package readers

import (
	"context"
	"errors"
	"fmt"

	"developers-cv/backend/internal/modules/resume/application/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResumeRootReader struct{ pool *pgxpool.Pool }

func NewResumeRootReader(pool *pgxpool.Pool) *ResumeRootReader { return &ResumeRootReader{pool: pool} }

func (r *ResumeRootReader) GetResumeRoot(ctx context.Context, q queries.GetResumeRootQuery) (queries.ResumeRootView, error) {
	if r == nil || r.pool == nil {
		return queries.ResumeRootView{}, fmt.Errorf("get resume root: reader dependencies are nil")
	}
	return r.getResumeRootTyped(ctx, q)
}

func (r *ResumeRootReader) getResumeRootTyped(ctx context.Context, q queries.GetResumeRootQuery) (queries.ResumeRootView, error) {
	var v queries.ResumeRootView
	var availabilityFrom, noteRU, noteEN, resumeLink, themePrimary, themeSecondary, availabilityStatus *string
	var rootID, title, firstRU, firstEN, secondRU, secondEN, highlightRU, highlightEN, thirdRU, thirdEN string
	if err := r.pool.QueryRow(ctx, `
		SELECT r.id::text, r.title, r.created_at, r.updated_at,
		       p.greeting_first_ru, p.greeting_first_en, p.greeting_second_ru, p.greeting_second_en,
		       p.greeting_highlight_ru, p.greeting_highlight_en, p.greeting_third_ru, p.greeting_third_en,
		       p.availability_status, p.availability_from::text, p.availability_note_ru, p.availability_note_en,
		       p.resume_link, p.theme_primary, p.theme_secondary
		FROM resume_roots r JOIN resume_presentations p ON p.resume_id = r.id
		WHERE r.id = $1 AND r.user_id = $2`, q.ResumeID.UUID(), q.UserID.UUID()).Scan(
		&rootID, &title, &v.CreatedAt, &v.UpdatedAt, &firstRU, &firstEN, &secondRU, &secondEN,
		&highlightRU, &highlightEN, &thirdRU, &thirdEN, &availabilityStatus, &availabilityFrom,
		&noteRU, &noteEN, &resumeLink, &themePrimary, &themeSecondary); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return queries.ResumeRootView{}, queries.ErrResumeNotFound
		}
		return queries.ResumeRootView{}, fmt.Errorf("get resume root: %w", err)
	}
	v.ID, v.Title = rootID, title
	v.Presentation.Greeting = queries.ResumeRootGreeting{
		First: queries.ResumeRootLocalizedText{RU: firstRU, EN: firstEN}, Second: queries.ResumeRootLocalizedText{RU: secondRU, EN: secondEN},
		SecondHighlight: queries.ResumeRootLocalizedText{RU: highlightRU, EN: highlightEN}, Third: queries.ResumeRootLocalizedText{RU: thirdRU, EN: thirdEN},
	}
	if availabilityStatus != nil {
		v.Presentation.Availability = &queries.ResumeRootAvailability{Status: *availabilityStatus, From: availabilityFrom}
		if noteRU != nil || noteEN != nil {
			v.Presentation.Availability.Note = &queries.ResumeRootLocalizedText{RU: valueOrEmpty(noteRU), EN: valueOrEmpty(noteEN)}
		}
	}
	if resumeLink != nil {
		v.Presentation.ResumeLink = resumeLink
	}
	if themePrimary != nil || themeSecondary != nil {
		v.Presentation.Theme = &queries.ResumeRootThemeSettings{Primary: valueOrEmpty(themePrimary), Secondary: valueOrEmpty(themeSecondary)}
	}
	rows, err := r.pool.Query(ctx, `SELECT section_id, enabled, collapsed FROM resume_presentation_sections WHERE resume_id=$1 ORDER BY position`, q.ResumeID.UUID())
	if err != nil {
		return queries.ResumeRootView{}, fmt.Errorf("get resume root sections: %w", err)
	}
	defer rows.Close()
	v.Presentation.Sections = make([]queries.ResumeRootSectionSettings, 0)
	for rows.Next() {
		var s queries.ResumeRootSectionSettings
		if err := rows.Scan(&s.ID, &s.Enabled, &s.Collapsed); err != nil {
			return queries.ResumeRootView{}, fmt.Errorf("scan resume root section: %w", err)
		}
		v.Presentation.Sections = append(v.Presentation.Sections, s)
	}
	if err := rows.Err(); err != nil {
		return queries.ResumeRootView{}, fmt.Errorf("iterate resume root sections: %w", err)
	}
	return v, nil
}

func valueOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
