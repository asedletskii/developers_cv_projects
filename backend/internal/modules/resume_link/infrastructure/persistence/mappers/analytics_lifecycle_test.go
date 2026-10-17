package mappers

import (
	"testing"

	"github.com/stretchr/testify/require"

	resumedomain "developers-cv/backend/internal/modules/resume/domain"
	"developers-cv/backend/internal/modules/resume_link/domain"
)

func TestAnalyticsLifecycleToDomain(t *testing.T) {
	tests := []struct {
		name           string
		model          AnalyticsLifecycle
		wantLinkStatus domain.ResumeLinkStatus
		wantResume     resumedomain.ResumeStatus
		wantErr        bool
	}{
		{
			name:           "active link and opened resume",
			model:          AnalyticsLifecycle{LinkStatus: "active", ResumeStatus: "opened"},
			wantLinkStatus: domain.ResumeLinkStatusActive,
			wantResume:     resumedomain.ResumeStatusOpened,
		},
		{
			name:    "unknown link status",
			model:   AnalyticsLifecycle{LinkStatus: "unknown", ResumeStatus: "opened"},
			wantErr: true,
		},
		{
			name:    "unknown resume status",
			model:   AnalyticsLifecycle{LinkStatus: "active", ResumeStatus: "unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linkStatus, resumeStatus, err := AnalyticsLifecycleToDomain(tt.model)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantLinkStatus, linkStatus)
			require.Equal(t, tt.wantResume, resumeStatus)
		})
	}
}
