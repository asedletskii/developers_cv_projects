import { ActivityType } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

/** Виды активности с подписями на текущем языке. */
export const useActivityOptions = () => {
  const translate = useTranslate('CurriculumVitae.Activities')

  return {
    typeOptions: Object.values(ActivityType).map((type) => ({
      value: type,
      label: translate(`Types.${type}`),
    })),
  }
}
