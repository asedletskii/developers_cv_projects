import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { getAnalyticsRange } from '@features/curriculumVitae/utils'
import { useAnalyticsRequest } from './useAnalyticsRequest'

const SUMMARY_PERIOD_DAYS = 30

export const useResumeStatistics = (resumeId?: string) =>
  useAnalyticsRequest(
    async (signal) =>
      resumeId
        ? curriculumVitaeApi.getResumeStatistics(
            resumeId,
            getAnalyticsRange(SUMMARY_PERIOD_DAYS),
            signal,
          )
        : undefined,
    [resumeId],
  )
