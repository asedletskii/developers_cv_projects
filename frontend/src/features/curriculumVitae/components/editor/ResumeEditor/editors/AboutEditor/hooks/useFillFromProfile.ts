import { MouseEvent, useState } from 'react'
import { selectLkProfile } from '@features/curriculumVitae/store/lk'
import { Resume } from '@features/curriculumVitae/types'
import { getLocalizedProfile } from '@features/curriculumVitae/utils'
import { useAppSelector } from '@shared/store'
import { Languages } from '@shared/types'
import { useFormContext } from 'react-hook-form'

export const useFillFromProfile = () => {
  const profile = useAppSelector(selectLkProfile)
  const { setValue } = useFormContext<Resume>()
  const [anchor, setAnchor] = useState<HTMLElement | null>(null)

  const fill = (language: Languages) => {
    if (!profile) return
    const localized = getLocalizedProfile(profile, language)
    const options = { shouldDirty: true, shouldTouch: true }

    setValue('profile.name', localized.name ?? '', options)
    setValue('profile.surname', localized.surname ?? '', options)
    setValue('profile.patronymic', localized.patronymic, options)
    setValue('profile.position', localized.position ?? '', options)
    setValue('profile.location', localized.location ?? '', options)
    if (profile.photoPath) setValue('profile.photoPath', profile.photoPath, options)

    setAnchor(null)
  }

  return {
    hasProfile: !!profile,
    anchor,
    openMenu: (e: MouseEvent<HTMLElement>) => setAnchor(e.currentTarget),
    closeMenu: () => setAnchor(null),
    fill,
  }
}
