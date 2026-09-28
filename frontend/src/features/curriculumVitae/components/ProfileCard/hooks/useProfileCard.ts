import { LkProfile } from '@features/curriculumVitae/types'
import { getLocalizedProfile } from '@features/curriculumVitae/utils'
import { useLanguage, useScreenSizeValue } from '@shared/hooks'
import { getInitials, getPersonFIO } from '@shared/utils'

export const useProfileCard = (profile: LkProfile) => {
  const { language } = useLanguage()
  const localized = getLocalizedProfile(profile, language)
  const isCompact = useScreenSizeValue('sm', true, false)

  return {
    localized,
    fio: getPersonFIO({
      name: localized.name ?? '',
      surname: localized.surname ?? '',
      patronymic: localized.patronymic,
    }),
    initials: getInitials(`${localized.name ?? ''} ${localized.surname ?? ''}`),
    isCompact,
    avatarSize: isCompact ? 72 : 144,
  }
}
