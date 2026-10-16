import { AvailabilityStatus } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

export const useHeroOptions = () => {
  const translate = useTranslate('CurriculumVitae.Availability')

  return {
    availabilityOptions: Object.values(AvailabilityStatus).map((status) => ({
      value: status,
      label: translate(status),
    })),
  }
}
