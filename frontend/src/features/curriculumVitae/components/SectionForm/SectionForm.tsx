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
  previewSection?: SectionId
  previewHighlight?: string
  onSaved: () => Promise<void>
  onDirtyChange: (isDirty: boolean) => void
  children: ReactNode
}

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

        {(isDirty || status !== 'idle') && (
          <SectionFormActions status={status} isDirty={isDirty} onReset={reset} />
        )}
      </Stack>
    </FormProvider>
  )
}
