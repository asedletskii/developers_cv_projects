import { LkLocalizedProfile, LkProfile } from '@features/curriculumVitae/types'
import { Languages } from '@shared/types'

/**
 * Профиль владельца на языке интерфейса. У существующих пользователей `en` пока не заполнен
 * (миграция перенесла старые данные только в `ru`), поэтому недостающее поле берём из `ru`,
 * а не показываем пустым.
 */
export const getLocalizedProfile = (
  profile: LkProfile,
  language: Languages,
): LkLocalizedProfile => {
  const localized = language === Languages.En ? profile.en : profile.ru
  const fallback = profile.ru

  return {
    name: localized.name ?? fallback.name,
    surname: localized.surname ?? fallback.surname,
    patronymic: localized.patronymic ?? fallback.patronymic,
    position: localized.position ?? fallback.position,
    location: localized.location ?? fallback.location,
  }
}
