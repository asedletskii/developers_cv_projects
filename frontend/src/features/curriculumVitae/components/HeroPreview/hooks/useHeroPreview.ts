import { Resume } from '@features/curriculumVitae/types'
import { useFormContext, useWatch } from 'react-hook-form'

/** Данные предпросмотра первого экрана: значения формы «как сейчас набрано». */
export const useHeroPreview = () => {
  const { control } = useFormContext<Resume>()

  return {
    greeting: useWatch({ control, name: 'overview.greeting' }),
    availability: useWatch({ control, name: 'overview.availability' }),
    resumeLink: useWatch({ control, name: 'overview.resumeLink' }),
    theme: useWatch({ control, name: 'overview.theme' }),
  }
}
