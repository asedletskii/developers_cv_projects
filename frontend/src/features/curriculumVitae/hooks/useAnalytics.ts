import { useContext } from 'react'
import { AnalyticsContext } from '@features/curriculumVitae/components/AnalyticsProvider'

export const useAnalytics = () => useContext(AnalyticsContext)
