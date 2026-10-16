import { WorkFormat } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

export const useProfileOptions = () => {
  const translate = useTranslate('CurriculumVitae.About')

  return {
    workFormatOptions: Object.values(WorkFormat).map((format) => ({
      value: format,
      label: translate(`WorkFormat.${format}`),
    })),
  }
}
