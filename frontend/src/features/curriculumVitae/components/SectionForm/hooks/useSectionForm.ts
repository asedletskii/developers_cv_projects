import { useEffect, useState } from 'react'
import { Resume, ResumeEditorTab } from '@features/curriculumVitae/types'
import { saveResumeSection } from '@features/curriculumVitae/utils'
import { useForm } from 'react-hook-form'

export type SaveStatus = 'idle' | 'saving' | 'saved' | 'error'

type Params = {
  userId: string
  view: Resume
  values: Resume
  section: ResumeEditorTab
  onSaved: () => Promise<void>
  onDirtyChange: (isDirty: boolean) => void
}

export const useSectionForm = ({
  userId,
  view,
  values,
  section,
  onSaved,
  onDirtyChange,
}: Params) => {
  const form = useForm<Resume>({ values })
  const { dirtyFields } = form.formState
  const isDirty = Object.keys(dirtyFields).length > 0

  const [status, setStatus] = useState<SaveStatus>('idle')

  useEffect(() => {
    if (!isDirty) return

    const warn = (e: BeforeUnloadEvent) => e.preventDefault()
    window.addEventListener('beforeunload', warn)
    return () => window.removeEventListener('beforeunload', warn)
  }, [isDirty])

  useEffect(() => {
    onDirtyChange(isDirty)
  }, [isDirty, onDirtyChange])

  const submit = form.handleSubmit(async (values) => {
    setStatus('saving')

    try {
      await saveResumeSection(userId, view, section, values)
      await onSaved()
      setStatus('saved')
    } catch {
      setStatus('error')
    }
  })

  return {
    form,
    isDirty,
    status,
    submit,
    reset: () => form.reset(values),
    resetStatus: () => setStatus('idle'),
  }
}
