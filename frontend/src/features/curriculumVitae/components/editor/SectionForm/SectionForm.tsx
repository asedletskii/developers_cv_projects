import { ReactNode } from 'react'
import { SectionStylePicker } from '@features/curriculumVitae/components/editor/form/SectionStylePicker'
import { SectionPreview } from '@features/curriculumVitae/components/editor/SectionPreview'
import { glassCard } from '@features/curriculumVitae/styles'
import { Resume, ResumeEditorTab, SectionId } from '@features/curriculumVitae/types'
import { Stack } from '@mui/material'
import { FormProvider } from 'react-hook-form'
import { useSectionForm } from './hooks'
import { SectionFormActions } from './SectionFormActions'

type Props = {
  view: Resume
  values: Resume
  section: ResumeEditorTab
  previewSection?: SectionId
  previewHighlight?: string
  hasStylePicker?: boolean
  onSaved: () => Promise<void>
  onDirtyChange: (isDirty: boolean) => void
  children: ReactNode
}

export const SectionForm = ({
  view,
  values,
  section,
  previewSection,
  previewHighlight,
  hasStylePicker = true,
  onSaved,
  onDirtyChange,
  children,
}: Props) => {
  const { form, isDirty, status, submit, reset, resetStatus } = useSectionForm({
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

        {previewSection && hasStylePicker && (
          <SectionStylePicker
            view={view}
            section={previewSection}
            isDirty={isDirty}
            onSaved={onSaved}
          />
        )}

        {children}

        {(isDirty || status !== 'idle') && (
          <SectionFormActions status={status} isDirty={isDirty} onReset={reset} />
        )}
      </Stack>
    </FormProvider>
  )
}
