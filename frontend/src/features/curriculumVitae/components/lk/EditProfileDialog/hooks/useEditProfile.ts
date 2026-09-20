import { useState } from 'react'
import { updateLkProfile } from '@features/curriculumVitae/store/lk'
import { LkLocalizedProfile, LkProfile } from '@features/curriculumVitae/types'
import { toOptionalString } from '@features/curriculumVitae/utils'
import { useAppDispatch } from '@shared/store'
import { Languages } from '@shared/types'

type Params = {
  userId: string
  profile: LkProfile
  onClose: () => void
}

const toBody = (localized: LkLocalizedProfile): LkLocalizedProfile => ({
  name: toOptionalString(localized.name),
  surname: toOptionalString(localized.surname),
  patronymic: toOptionalString(localized.patronymic),
  position: toOptionalString(localized.position),
  location: toOptionalString(localized.location),
})

export const useEditProfile = ({ userId, profile, onClose }: Params) => {
  const dispatch = useAppDispatch()

  const [language, setLanguage] = useState(Languages.Ru)
  const [ru, setRu] = useState(profile.ru)
  const [en, setEn] = useState(profile.en)
  const [photoPath, setPhotoPath] = useState(profile.photoPath ?? '')
  const [isSaving, setIsSaving] = useState(false)
  const [hasError, setHasError] = useState(false)

  const current = language === Languages.En ? en : ru
  const setCurrent = language === Languages.En ? setEn : setRu

  const updateField = (field: keyof LkLocalizedProfile, value: string) =>
    setCurrent((prev) => ({ ...prev, [field]: value }))

  const save = async () => {
    setIsSaving(true)
    setHasError(false)

    try {
      await dispatch(
        updateLkProfile({
          userId,
          ru: toBody(ru),
          en: toBody(en),
          photoPath: toOptionalString(photoPath),
        }),
      ).unwrap()
      onClose()
    } catch {
      setHasError(true)
      setIsSaving(false)
    }
  }

  return {
    language,
    setLanguage,
    current,
    updateField,
    photoPath,
    setPhotoPath,
    isSaving,
    hasError,
    save,
  }
}
