import { ContactKey } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

/** Контакты, которые можно скрыть до клика, с подписями на текущем языке. */
export const useContactOptions = () => {
  const translate = useTranslate('Lk.Edit.Contacts')

  return {
    hiddenContactOptions: Object.values(ContactKey).map((key) => ({
      value: key,
      label: translate(key),
    })),
  }
}
