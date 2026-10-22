package readers

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"developers-cv/backend/internal/modules/analytics/application/queries"
	"developers-cv/backend/internal/modules/analytics/infrastructure/persistence/gormmodel"
)

type dailyStatisticsRow struct {
	Date                         string
	Views, Reads, UniqueVisitors int64
}

type sectionStatisticsRow struct {
	Section        string
	Reached, Exits int64
}

type technologyStatisticsRow struct {
	Name                 string
	Selections, Visitors int64
}

func readDeviceStatistics(
	db *gorm.DB,
	linkID uuid.UUID,
	from, to time.Time,
) (map[string]int64, error) {
	var rows []deviceStatisticsRow
	err := db.Model(&gormmodel.AnalyticsView{}).
		Select("device, COUNT(*) AS count").
		Where("link_id = ?", linkID).
		Where("traffic = ?", "human").
		Where("started_at >= ? AND started_at < ?", from, to).
		Group("device").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return deviceStatisticsMap(rows), nil
}

func readResumeDeviceStatistics(
	db *gorm.DB,
	resumeID uuid.UUID,
	from, to time.Time,
) (map[string]int64, error) {
	var rows []deviceStatisticsRow
	err := db.Table("analytics.analytics_views AS view").
		Select("view.device, COUNT(*) AS count").
		Where("view.resume_id = ?", resumeID).
		Where("view.traffic = ?", "human").
		Where("view.started_at >= ? AND view.started_at < ?", from, to).
		Group("view.device").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return deviceStatisticsMap(rows), nil
}

type deviceStatisticsRow struct {
	Device string
	Count  int64
}

func deviceStatisticsMap(rows []deviceStatisticsRow) map[string]int64 {
	devices := map[string]int64{
		"desktop": 0,
		"mobile":  0,
		"tablet":  0,
		"unknown": 0,
	}
	for _, row := range rows {
		devices[row.Device] = row.Count
	}
	return devices
}

func readReferrerStatistics(
	db *gorm.DB,
	linkID uuid.UUID,
	from, to time.Time,
) ([]queries.ReferrerStat, error) {
	var rows []referrerStatisticsRow
	err := db.Model(&gormmodel.AnalyticsView{}).
		Select("NULLIF(referrer_host, '') AS host, COUNT(*) AS count").
		Where("link_id = ?", linkID).
		Where("traffic = ?", "human").
		Where("started_at >= ? AND started_at < ?", from, to).
		Group("referrer_host").
		Order("count DESC, host ASC NULLS LAST").
		Limit(maxReferrerStatistics).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return mapReferrerStatistics(rows), nil
}

func readResumeReferrerStatistics(
	db *gorm.DB,
	resumeID uuid.UUID,
	from, to time.Time,
) ([]queries.ReferrerStat, error) {
	var rows []referrerStatisticsRow
	err := db.Table("analytics.analytics_views AS view").
		Select("NULLIF(view.referrer_host, '') AS host, COUNT(*) AS count").
		Where("view.resume_id = ?", resumeID).
		Where("view.traffic = ?", "human").
		Where("view.started_at >= ? AND view.started_at < ?", from, to).
		Group("view.referrer_host").
		Order("count DESC, host ASC NULLS LAST").
		Limit(maxReferrerStatistics).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return mapReferrerStatistics(rows), nil
}

const maxReferrerStatistics = 20

type referrerStatisticsRow struct {
	Host  *string
	Count int64
}

func mapReferrerStatistics(rows []referrerStatisticsRow) []queries.ReferrerStat {
	result := make([]queries.ReferrerStat, 0, len(rows))
	for _, row := range rows {
		result = append(result, queries.ReferrerStat{
			Host: row.Host, Views: row.Count,
		})
	}
	return result
}

func readFilteredStatistics(
	db *gorm.DB,
	linkID uuid.UUID,
	from, to time.Time,
) (queries.FilteredStatistics, error) {
	var result queries.FilteredStatistics
	err := db.Model(&gormmodel.AnalyticsView{}).
		Select(`
			COUNT(*) FILTER (WHERE traffic = 'bot') AS bot,
			COUNT(*) FILTER (WHERE traffic = 'owner') AS owner,
			COUNT(*) FILTER (WHERE traffic = 'suspect') AS suspect
		`).
		Where("link_id = ?", linkID).
		Where("started_at >= ? AND started_at < ?", from, to).
		Scan(&result).Error
	return result, err
}

func readResumeFilteredStatistics(
	db *gorm.DB,
	resumeID uuid.UUID,
	from, to time.Time,
) (queries.FilteredStatistics, error) {
	var result queries.FilteredStatistics
	err := db.Table("analytics.analytics_views AS view").
		Select(`
			COUNT(*) FILTER (WHERE view.traffic = 'bot') AS bot,
			COUNT(*) FILTER (WHERE view.traffic = 'owner') AS owner,
			COUNT(*) FILTER (WHERE view.traffic = 'suspect') AS suspect
		`).
		Where("view.resume_id = ?", resumeID).
		Where("view.started_at >= ? AND view.started_at < ?", from, to).
		Scan(&result).Error
	return result, err
}
