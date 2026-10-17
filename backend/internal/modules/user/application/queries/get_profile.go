package queries

import (
	"context"
	"fmt"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type GetProfileQuery struct{ UserID sharedvalueobjects.ID }

type LocalizedProfileView struct {
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
	Position   *string `json:"position,omitempty"`
	Location   *string `json:"location,omitempty"`
}

type ProfileView struct {
	ID       string               `json:"userId"`
	RU       LocalizedProfileView `json:"ru"`
	EN       LocalizedProfileView `json:"en"`
	PhotoURL *string              `json:"photoPath,omitempty"`
}

type GetProfileReader interface {
	GetProfile(context.Context, GetProfileQuery) (ProfileView, error)
}

type GetProfileHandler struct{ reader GetProfileReader }

func NewGetProfileHandler(reader GetProfileReader) *GetProfileHandler {
	return &GetProfileHandler{reader: reader}
}

func (h *GetProfileHandler) Handle(ctx context.Context, query GetProfileQuery) (ProfileView, error) {
	if h == nil || h.reader == nil {
		return ProfileView{}, fmt.Errorf("get profile handler dependencies are nil")
	}
	if query.UserID.IsZero() {
		return ProfileView{}, fmt.Errorf("invalid profile user id")
	}
	return h.reader.GetProfile(ctx, query)
}
