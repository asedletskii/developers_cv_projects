import { useState } from 'react'
import { updateLkResumeTitle } from '@features/curriculumVitae/store/lk'
import { LkResumeSummary } from '@features/curriculumVitae/types'
import { useAppDispatch } from '@shared/store'

type Params = {
  userId: string
  resume: LkResumeSummary
  onClose: () => void
}

export const useEditResumeTitle = ({ userId, resume, onClose }: Params) => {
  const dispatch = useAppDispatch()

  const [title, setTitle] = useState(resume.title)
  const [isSaving, setIsSaving] = useState(false)
  const [hasError, setHasError] = useState(false)

  const save = async () => {
    setIsSaving(true)
    setHasError(false)

    try {
      await dispatch(updateLkResumeTitle({ userId, resumeId: resume.id, title })).unwrap()
      onClose()
    } catch {
      setHasError(true)
      setIsSaving(false)
    }
  }

  return { title, setTitle, isSaving, hasError, canSave: !isSaving && !!title.trim(), save }
}
