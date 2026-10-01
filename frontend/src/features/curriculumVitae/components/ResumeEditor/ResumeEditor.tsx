import { SectionForm } from '@features/curriculumVitae/components/SectionForm'
import { Resume } from '@features/curriculumVitae/types'
import { Stack, Tab, Tabs } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { EDITOR_TABS, TAB_SX } from './constants'
import { useResumeEditor } from './hooks'

type Props = {
  userId: string
  view: Resume
  onSaved: () => Promise<void>
  isDirtyRef: { current: boolean }
}

export const ResumeEditor = ({ userId, view, onSaved, isDirtyRef }: Props) => {
  const translateTabs = useTranslate('Lk.Edit.Tabs')
  const { values, activeTab, changeTab, setDirty } = useResumeEditor({ view, isDirtyRef })

  return (
    <Stack spacing={3}>
      <Tabs
        value={activeTab.id}
        onChange={(_, id) => changeTab(id)}
        variant="scrollable"
        scrollButtons="auto"
      >
        {EDITOR_TABS.map(({ id }) => (
          <Tab key={id} value={id} label={translateTabs(id)} sx={TAB_SX} />
        ))}
      </Tabs>

      <SectionForm
        key={activeTab.id}
        userId={userId}
        view={view}
        values={values}
        section={activeTab.id}
        previewSection={activeTab.preview}
        previewHighlight={activeTab.highlight}
        onSaved={onSaved}
        onDirtyChange={setDirty}
      >
        <activeTab.Editor />
      </SectionForm>
    </Stack>
  )
}
