import { Resume } from '@features/curriculumVitae/types'
import { useFormContext, useWatch } from 'react-hook-form'

export const useHeroPreview = () => {
  const { control } = useFormContext<Resume>()

  return {
    greeting: useWatch({ control, name: 'overview.greeting' }),
    availability: useWatch({ control, name: 'overview.availability' }),
    theme: useWatch({ control, name: 'overview.theme' }),
  }
}
