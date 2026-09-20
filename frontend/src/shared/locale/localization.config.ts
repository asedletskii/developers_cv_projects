import { appLocaleDict } from '@app/locale'
import { Languages } from '@shared/types'
import { getInitialLanguage } from '@shared/utils/getInitialLanguage'
import i18next, { i18n } from 'i18next'

let _instance: i18n | null = null

export const initI18n = (): i18n => {
  if (_instance) return _instance

  const initialLng = getInitialLanguage(Languages.Ru)

  i18next.init({
    lng: initialLng,
    fallbackLng: Languages.Ru,
    resources: Object.fromEntries(
      Object.entries(appLocaleDict).map(([lng, dict]) => [lng, { translation: dict }]),
    ) as any,

    interpolation: {
      escapeValue: false,
    },

    returnNull: false,
    returnEmptyString: false,
    parseMissingKeyHandler(key) {
      return key
    },
  })

  _instance = i18next
  return i18next
}
