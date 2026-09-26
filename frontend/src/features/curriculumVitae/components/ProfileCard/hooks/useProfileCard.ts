import { LkProfile } from '@features/curriculumVitae/types'
import { useScreenSizeValue } from '@shared/hooks'
import { getInitials, getPersonFIO } from '@shared/utils'

export const useProfileCard = (profile: LkProfile) => {
  const isCompact = useScreenSizeValue('sm', true, false)

  return {
    fio: getPersonFIO(profile),
    initials: getInitials(`${profile.name} ${profile.surname}`),
    isCompact,
    avatarSize: isCompact ? 72 : 144,
  }
}
