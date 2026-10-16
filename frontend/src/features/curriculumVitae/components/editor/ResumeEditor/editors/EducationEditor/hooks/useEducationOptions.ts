import { useTranslate } from '@shared/hooks'
import { EDUCATION_LEVELS } from '../constants'

export const useEducationOptions = () => {
  const translate = useTranslate('CurriculumVitae.Education')

  return {
    levelOptions: EDUCATION_LEVELS.map((level) => ({
      value: level,
      label: translate(`Levels.level${level}`),
    })),
  }
}
