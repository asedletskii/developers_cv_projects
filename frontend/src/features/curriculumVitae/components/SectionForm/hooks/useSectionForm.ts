import { useEffect, useState } from 'react'
import { Person, ResumeEditorTab, ResumeView } from '@features/curriculumVitae/types'
import { saveResumeSection } from '@features/curriculumVitae/utils'
import { useForm } from 'react-hook-form'

export type SaveStatus = 'idle' | 'saving' | 'saved' | 'error'

type Params = {
  userId: string
  view: ResumeView
  person: Person
  section: ResumeEditorTab

  /** Вызывается после успешного сохранения — экран перечитывает резюме с бэка */
  onSaved: () => Promise<void>
  onDirtyChange: (isDirty: boolean) => void
}

/** Форма одной вкладки: значения, сохранение своей секции и предупреждение о несохранённых правках. */
export const useSectionForm = ({
  userId,
  view,
  person,
  section,
  onSaved,
  onDirtyChange,
}: Params) => {
  const form = useForm<Person>({ values: person })
  const { isDirty } = form.formState

  const [status, setStatus] = useState<SaveStatus>('idle')

  // Несохранённые правки не должны теряться от случайной перезагрузки или закрытия вкладки
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
    reset: () => form.reset(),
    resetStatus: () => setStatus('idle'),
  }
}
