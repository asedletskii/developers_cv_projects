package commands

import (
	"context"
	"developers-cv/backend/internal/modules/resume/application/ports"
	"fmt"

	"developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume/domain/enums"
	"developers-cv/backend/internal/modules/resume/domain/valueobjects"
	useridentity "developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type LocalizedInput struct{ RU, EN string }
type GreetingInput struct{ First, Second, SecondHighlight, Third LocalizedInput }
type AvailabilityInput struct {
	Status string
	From   *string
	Note   *LocalizedInput
}
type ThemeInput struct{ Primary, Secondary *string }
type SectionInput struct {
	ID                 string
	Enabled, Collapsed bool
}
type PresentationInput struct {
	Theme        *ThemeInput
	Greeting     GreetingInput
	Availability *AvailabilityInput
	ResumeLink   *string
	Sections     []SectionInput
}

func (i PresentationInput) toDomain() (domain.ResumePresentation, error) {
	loc := func(x LocalizedInput) (valueobjects.LocalizedText, error) {
		return valueobjects.NewLocalizedText(valueobjects.LocalizedTextParams{RU: x.RU, EN: x.EN})
	}
	f, e := loc(i.Greeting.First)
	if e != nil {
		return domain.ResumePresentation{}, e
	}
	sec, e := loc(i.Greeting.Second)
	if e != nil {
		return domain.ResumePresentation{}, e
	}
	hi, e := loc(i.Greeting.SecondHighlight)
	if e != nil {
		return domain.ResumePresentation{}, e
	}
	th, e := loc(i.Greeting.Third)
	if e != nil {
		return domain.ResumePresentation{}, e
	}
	greeting := valueobjects.NewGreeting(valueobjects.GreetingParams{First: f, Second: sec, SecondHighlight: hi, Third: th})
	var theme *valueobjects.ThemeSettings
	if i.Theme != nil {
		var a, b *valueobjects.Color
		if i.Theme.Primary != nil {
			v, e := valueobjects.NewColor(*i.Theme.Primary)
			if e != nil {
				return domain.ResumePresentation{}, e
			}
			a = &v
		}
		if i.Theme.Secondary != nil {
			v, e := valueobjects.NewColor(*i.Theme.Secondary)
			if e != nil {
				return domain.ResumePresentation{}, e
			}
			b = &v
		}
		v, e := valueobjects.NewThemeSettings(valueobjects.ThemeSettingsParams{Primary: a, Secondary: b})
		if e != nil {
			return domain.ResumePresentation{}, e
		}
		theme = &v
	}
	var av *valueobjects.Availability
	if i.Availability != nil {
		var from *valueobjects.Date
		if i.Availability.From != nil {
			v, e := valueobjects.ParseDate(*i.Availability.From)
			if e != nil {
				return domain.ResumePresentation{}, e
			}
			from = &v
		}
		var note *valueobjects.LocalizedText
		if i.Availability.Note != nil {
			v, e := loc(*i.Availability.Note)
			if e != nil {
				return domain.ResumePresentation{}, e
			}
			note = &v
		}
		v, e := valueobjects.NewAvailability(valueobjects.AvailabilityParams{Status: enums.AvailabilityStatus(i.Availability.Status), From: from, Note: note})
		if e != nil {
			return domain.ResumePresentation{}, e
		}
		av = &v
	}
	var link *valueobjects.AssetPath
	if i.ResumeLink != nil {
		v, e := valueobjects.NewAssetPath(*i.ResumeLink)
		if e != nil {
			return domain.ResumePresentation{}, e
		}
		link = &v
	}
	sections := make([]valueobjects.SectionSettings, len(i.Sections))
	for n, v := range i.Sections {
		x, e := valueobjects.NewSectionSettings(valueobjects.SectionSettingsParams{ID: enums.SectionID(v.ID), Enabled: v.Enabled, Collapsed: v.Collapsed})
		if e != nil {
			return domain.ResumePresentation{}, e
		}
		sections[n] = x
	}
	return domain.ResumePresentation{Theme: theme, Greeting: greeting, Availability: av, ResumeLink: link, Sections: sections}, nil
}

type CreateResumeRootCommand struct {
	ID           valueobjects.ID
	UserID       useridentity.ID
	Title        string
	Presentation PresentationInput
}
type UpdateResumeRootCommand struct {
	ResumeID     valueobjects.ID
	UserID       useridentity.ID
	Title        string
	Presentation PresentationInput
}

type GetResumeRootCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}
type DeleteResumeRootCommand struct {
	ResumeID valueobjects.ID
	UserID   useridentity.ID
}

type CreateResumeRootRepository = ports.CreateResumeRootRepository
type CreateResumeRootRepositoryCommand = ports.CreateResumeRootRepositoryCommand
type GetResumeRootRepository = ports.GetResumeRootRepository
type UpdateResumeRootRepository = ports.UpdateResumeRootRepository
type UpdateResumeRootRepositoryCommand = ports.UpdateResumeRootRepositoryCommand
type DeleteResumeRootRepository = ports.DeleteResumeRootRepository

type CreateResumeRootHandler struct {
	repository CreateResumeRootRepository
	uow        UnitOfWork
}

func NewCreateResumeRootHandler(repository CreateResumeRootRepository, uow UnitOfWork) *CreateResumeRootHandler {
	return &CreateResumeRootHandler{repository: repository, uow: uow}
}
func (h *CreateResumeRootHandler) Handle(ctx context.Context, c CreateResumeRootCommand) (valueobjects.ID, error) {
	if h == nil || h.repository == nil || h.uow == nil {
		return valueobjects.ID{}, fmt.Errorf("create resume root handler dependencies are nil")
	}
	if c.ID.IsZero() {
		c.ID = valueobjects.NewID()
	}
	presentation, err := c.Presentation.toDomain()
	if err != nil {
		return valueobjects.ID{}, err
	}
	err = h.uow.Do(ctx, func(tx context.Context) error {
		resume, err := domain.NewResumeRoot(domain.ResumeRootParams{ID: c.ID, UserID: c.UserID, Title: c.Title, Presentation: presentation})
		if err != nil {
			return err
		}
		return h.repository.CreateRoot(tx, ports.CreateResumeRootRepositoryCommand{Resume: resume})
	})
	if err != nil {
		return valueobjects.ID{}, err
	}
	return c.ID, nil
}

type UpdateResumeRootHandler struct {
	repository GetResumeRootRepository
	updater    UpdateResumeRootRepository
	uow        UnitOfWork
}

func NewUpdateResumeRootHandler(repository GetResumeRootRepository, updater UpdateResumeRootRepository, uow UnitOfWork) *UpdateResumeRootHandler {
	return &UpdateResumeRootHandler{repository: repository, updater: updater, uow: uow}
}
func (h *UpdateResumeRootHandler) Handle(ctx context.Context, c UpdateResumeRootCommand) error {
	if h == nil || h.repository == nil || h.updater == nil || h.uow == nil {
		return fmt.Errorf("update resume root handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid update resume root command")
	}
	presentation, err := c.Presentation.toDomain()
	if err != nil {
		return err
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		current, err := h.repository.GetRoot(tx, ports.GetResumeRootRepositoryCommand{ResumeID: c.ResumeID, UserID: c.UserID})
		if err != nil {
			return err
		}
		if err := current.Replace(domain.ResumeRootParams{ID: current.ID(), UserID: current.UserID(), Title: c.Title, Presentation: presentation}); err != nil {
			return err
		}
		return h.updater.UpdateRoot(tx, ports.UpdateResumeRootRepositoryCommand{Resume: current})
	})
}

type DeleteResumeRootHandler struct {
	repository DeleteResumeRootRepository
	uow        UnitOfWork
}

func NewDeleteResumeRootHandler(repository DeleteResumeRootRepository, uow UnitOfWork) *DeleteResumeRootHandler {
	return &DeleteResumeRootHandler{repository: repository, uow: uow}
}
func (h *DeleteResumeRootHandler) Handle(ctx context.Context, c DeleteResumeRootCommand) error {
	if h == nil || h.repository == nil || h.uow == nil {
		return fmt.Errorf("delete resume root handler dependencies are nil")
	}
	if c.ResumeID.IsZero() || c.UserID.IsZero() {
		return fmt.Errorf("invalid delete resume root command")
	}
	return h.uow.Do(ctx, func(tx context.Context) error {
		return h.repository.DeleteRoot(tx, ports.DeleteResumeRootRepositoryCommand{ResumeID: c.ResumeID, UserID: c.UserID})
	})
}
