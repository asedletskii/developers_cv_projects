import { Person } from '@features/curriculumVitae/types'
import { useFormContext, useWatch } from 'react-hook-form'

/** Данные предпросмотра первого экрана: значения формы «как сейчас набрано». */
export const useHeroPreview = () => {
  const { control } = useFormContext<Person>()

  return {
    greeting: useWatch({ control, name: 'greeting' }),
    availability: useWatch({ control, name: 'availability' }),
    resumeLink: useWatch({ control, name: 'resumeLink' }),
    theme: useWatch({ control, name: 'theme' }),
  }
}
