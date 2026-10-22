package contracts

import (
	"context"

	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

type PublicInteractionEventRecorder interface {
	RecordContactView(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID, string) error
	RecordDownload(context.Context, sharedvalueobjects.ID, sharedvalueobjects.ID) error
}
