import { useState } from 'react'
import { createLkResume } from '@features/curriculumVitae/store/lk'
import { ResumeLanguage } from '@features/curriculumVitae/types'
import { ROUTES } from '@shared/constants'
import { useAppDispatch } from '@shared/store'
import { generatePath, useNavigate } from 'react-router'

export const useCreateResume = (userId: string) => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [title, setTitle] = useState('')
  const [language, setLanguage] = useState<ResumeLanguage>('ru')
  const [isSaving, setIsSaving] = useState(false)
  const [hasError, setHasError] = useState(false)

  const create = async () => {
    setIsSaving(true)
    setHasError(false)

    try {
      const id = await dispatch(createLkResume({ userId, title: title.trim(), language })).unwrap()
      navigate(generatePath(ROUTES.lkResume, { userId, cvId: id }))
    } catch {
      setHasError(true)
      setIsSaving(false)
    }
  }

  return {
    title,
    setTitle,
    language,
    setLanguage,
    isSaving,
    hasError,
    canCreate: !isSaving && !!title.trim(),
    create,
  }
}
