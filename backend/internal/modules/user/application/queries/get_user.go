package queries

import (
	"context"
	"fmt"

	"developers-cv/backend/internal/modules/user/domain/valueobjects"
)

type GetUserQuery struct {
	UserID valueobjects.ID
}

type GetUserReader interface {
	GetUser(ctx context.Context, query GetUserQuery) (User, error)
}

type GetUserHandler struct {
	reader GetUserReader
}

func NewGetUserHandler(reader GetUserReader) *GetUserHandler {
	return &GetUserHandler{reader: reader}
}

func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (User, error) {
	if h == nil || h.reader == nil {
		return User{}, fmt.Errorf("get user handler dependencies are nil")
	}
	if query.UserID.IsZero() {
		return User{}, fmt.Errorf("invalid user id")
	}
	return h.reader.GetUser(ctx, query)
}
