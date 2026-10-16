import { useState } from 'react'
import { updateLkResumeTitle } from '@features/curriculumVitae/store/lk'
import { useAppDispatch } from '@shared/store'

type Params = {
  resume: { id: string; title: string; version?: number }
  onClose: () => void
}

export const useEditResumeTitle = ({ resume, onClose }: Params) => {
  const dispatch = useAppDispatch()

  const [title, setTitle] = useState(resume.title)
  const [isSaving, setIsSaving] = useState(false)
  const [hasError, setHasError] = useState(false)

  const save = async () => {
    setIsSaving(true)
    setHasError(false)

    try {
      await dispatch(
        updateLkResumeTitle({ resumeId: resume.id, version: resume.version ?? 0, title }),
      ).unwrap()
      onClose()
    } catch {
      setHasError(true)
      setIsSaving(false)
    }
  }

  return { title, setTitle, isSaving, hasError, canSave: !isSaving && !!title.trim(), save }
}
