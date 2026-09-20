import { useContext } from 'react'
import { AnalyticsContext } from '@features/curriculumVitae/components/analytics/AnalyticsProvider'

export const useAnalytics = () => useContext(AnalyticsContext)
