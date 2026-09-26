import { SectionForm } from '@features/curriculumVitae/components/SectionForm'
import { ResumeView } from '@features/curriculumVitae/types'
import { Stack, Tab, Tabs } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useResumeEditor } from './hooks'
import { EDITOR_TABS } from './constants'

type Props = {
  userId: string
  view: ResumeView
  onSaved: () => Promise<void>

  /** Есть ли несохранённые правки — экран спрашивает подтверждение при уходе с вкладки */
  isDirtyRef: { current: boolean }
}

/** Раздел «Редактирование»: вкладка на каждую секцию резюме, у каждой своя форма и своё сохранение.
 *  Активная вкладка живёт в `?tab=`, чтобы на неё можно было сослаться. */
export const ResumeEditor = ({ userId, view, onSaved, isDirtyRef }: Props) => {
  const translateTabs = useTranslate('Lk.Edit.Tabs')
  const { person, activeTab, changeTab, setDirty } = useResumeEditor({ view, isDirtyRef })

  return (
    <Stack spacing={3}>
      <Tabs
        value={activeTab.id}
        onChange={(_, id) => changeTab(id)}
        variant="scrollable"
        scrollButtons="auto"
      >
        {EDITOR_TABS.map(({ id }) => (
          <Tab key={id} value={id} label={translateTabs(id)} sx={{ textTransform: 'none' }} />
        ))}
      </Tabs>

      {/* key: при смене вкладки форма создаётся заново, а не тащит значения прошлой */}
      <SectionForm
        key={activeTab.id}
        userId={userId}
        view={view}
        person={person}
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
