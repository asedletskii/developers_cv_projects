package queries

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	application "developers-cv/backend/internal/modules/resume_link/application"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestActivePublicResumeLinkResolverResolvesResumeID(t *testing.T) {
	linkID := sharedvalueobjects.NewID()
	resumeID := sharedvalueobjects.NewID()
	reader := &resumeLinkReaderStub{view: ResumeLinkPublicView{
		ID: linkID.String(), ResumeID: resumeID.String(), Status: string(domain.ResumeLinkStatusActive),
	}}

	resolvedID, err := NewActivePublicResumeLinkResolver(reader).Resolve(context.Background(), linkID)

	require.NoError(t, err)
	require.Equal(t, resumeID, resolvedID)
	require.Equal(t, linkID, reader.query.LinkID)
}

func TestActivePublicResumeLinkResolverRejectsUnavailableLinksBeforeResumeRead(t *testing.T) {
	tests := []struct {
		name   string
		reader *resumeLinkReaderStub
		want   error
	}{
		{
			name: "missing",
			reader: &resumeLinkReaderStub{
				err: domain.ErrEntityNotFound,
			},
			want: application.ErrResumeLinkNotFound,
		},
		{
			name: "inactive",
			reader: &resumeLinkReaderStub{
				view: ResumeLinkPublicView{Status: string(domain.ResumeLinkStatusInactive)},
			},
			want: domain.ErrResumeLinkInactive,
		},
		{
			name: "unknown status",
			reader: &resumeLinkReaderStub{
				view: ResumeLinkPublicView{Status: "unknown"},
			},
			want: application.ErrResumeLinkNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resumeID, err := NewActivePublicResumeLinkResolver(tt.reader).Resolve(
				context.Background(),
				sharedvalueobjects.NewID(),
			)

			require.ErrorIs(t, err, tt.want)
			require.True(t, resumeID.IsZero())
		})
	}
}

func TestActivePublicResumeLinkResolverRejectsMalformedResolvedResumeID(t *testing.T) {
	resolver := NewActivePublicResumeLinkResolver(&resumeLinkReaderStub{
		view: ResumeLinkPublicView{ResumeID: "not-an-id", Status: string(domain.ResumeLinkStatusActive)},
	})

	_, err := resolver.Resolve(context.Background(), sharedvalueobjects.NewID())

	require.Error(t, err)
}
