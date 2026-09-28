package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type GetProfileQuery struct{ UserID valueobjects.ID }

type LocalizedProfile struct {
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
	Position   *string `json:"position,omitempty"`
	Location   *string `json:"location,omitempty"`
}

type Profile struct {
	ID        string           `json:"userId"`
	RU        LocalizedProfile `json:"ru"`
	EN        LocalizedProfile `json:"en"`
	PhotoPath *string          `json:"photoPath,omitempty"`
}

type GetProfileReader interface {
	GetProfile(ctx context.Context, query GetProfileQuery) (Profile, error)
}

type GetProfileHandler struct{ reader GetProfileReader }

func NewGetProfileHandler(reader GetProfileReader) *GetProfileHandler {
	return &GetProfileHandler{reader: reader}
}

func (h *GetProfileHandler) Handle(ctx context.Context, query GetProfileQuery) (Profile, error) {
	if h == nil || h.reader == nil {
		return Profile{}, fmt.Errorf("get profile handler dependencies are nil")
	}
	if query.UserID.IsZero() {
		return Profile{}, fmt.Errorf("invalid profile user id")
	}
	return h.reader.GetProfile(ctx, query)
}
