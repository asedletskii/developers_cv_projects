package readers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/resume/application/queries"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

func TestListResumesReaderRejectsNilReceiver(t *testing.T) {
	var reader *ListResumesReader

	_, err := reader.ListResumes(context.Background(), queries.ListResumesQuery{UserID: sharedvalueobjects.NewID()})

	require.ErrorContains(t, err, "reader is nil")
}

func TestDateFormatting(t *testing.T) {
	value := time.Date(2026, time.September, 17, 13, 45, 0, 0, time.UTC)

	require.Equal(t, "2026-09", yearMonth(value))

	yearMonthValue := value
	require.Equal(t, "2026-09", *optionalYearMonth(&yearMonthValue))

	require.Equal(t, "2026", *optionalYear(&yearMonthValue))
	require.Nil(t, optionalYearMonth(nil))
	require.Nil(t, optionalYear(nil))
}
