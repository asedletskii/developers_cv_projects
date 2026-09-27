import { ReactNode } from 'react'
import { SectionPreview } from '@features/curriculumVitae/components/SectionPreview'
import { glassCard } from '@features/curriculumVitae/styles'
import { Resume, ResumeEditorTab, SectionId } from '@features/curriculumVitae/types'
import { Stack } from '@mui/material'
import { FormProvider } from 'react-hook-form'
import { useSectionForm } from './hooks'
import { SectionFormActions } from './SectionFormActions'

type Props = {
  userId: string
  view: Resume
  values: Resume
  section: ResumeEditorTab

  /** Секция страницы, которую вкладка правит — по ней рисуется предпросмотр (нет — без него) */
  previewSection?: SectionId

  /** Участок секции, выделяемый в предпросмотре */
  previewHighlight?: string

  /** Вызывается после успешного сохранения — экран перечитывает резюме с бэка */
  onSaved: () => Promise<void>

  /** Есть ли несохранённые правки — экран спрашивает подтверждение при смене вкладки */
  onDirtyChange: (isDirty: boolean) => void
  children: ReactNode
}

/** Форма одной вкладки: своё состояние, свои кнопки «Сохранить»/«Сбросить», свой запрос на бэк. */
export const SectionForm = ({
  userId,
  view,
  values,
  section,
  previewSection,
  previewHighlight,
  onSaved,
  onDirtyChange,
  children,
}: Props) => {
  const { form, isDirty, status, submit, reset, resetStatus } = useSectionForm({
    userId,
    view,
    values,
    section,
    onSaved,
    onDirtyChange,
  })

  return (
    <FormProvider {...form}>
      <Stack
        component="form"
        spacing={3}
        onSubmit={submit}
        onChange={resetStatus}
        sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}
      >
        {previewSection && <SectionPreview section={previewSection} highlight={previewHighlight} />}

        {children}

        <SectionFormActions status={status} isDirty={isDirty} onReset={reset} />
      </Stack>
    </FormProvider>
  )
}
