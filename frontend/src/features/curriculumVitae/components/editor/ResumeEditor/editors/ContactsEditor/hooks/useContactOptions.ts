import { ContactKey } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

export const useContactOptions = () => {
  const translate = useTranslate('Lk.Edit.Contacts')

  return {
    hiddenContactOptions: Object.values(ContactKey).map((key) => ({
      value: key,
      label: translate(key),
    })),
  }
}
