import { LoadFailure } from '@features/curriculumVitae/components'
import { AnalyticsTab } from '@features/curriculumVitae/components/AnalyticsTab'
import { LinksTab } from '@features/curriculumVitae/components/LinksTab'
import { MainTab } from '@features/curriculumVitae/components/MainTab'
import { ResumeEditor } from '@features/curriculumVitae/components/ResumeEditor'
import { ResumeHeader } from '@features/curriculumVitae/components/ResumeHeader'
import { useOwnerResume } from '@features/curriculumVitae/hooks'
import { LinearProgress, Stack, Tab, Tabs } from '@mui/material'
import { LocalizationProvider } from '@mui/x-date-pickers'
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import { useLanguage, useTranslate } from '@shared/hooks'
import { useParams } from 'react-router'
import { TAB_SX, TOP_SECTIONS } from './constants'
import { useResumeSections } from './hooks'

/**
 * Правая часть личного кабинета в режиме редактирования (`/cv/lk/:userId/:cvId`): слева остаётся
 * карточка пользователя из layout-а, здесь — шапка резюме и разделы «Главная», «Редактирование»,
 * «Ссылки», «Аналитика». Активный раздел живёт в `?section=`, чтобы на него можно было сослаться.
 */
const LkResumeEditScreen = () => {
  const { userId, cvId } = useParams()
  const { language } = useLanguage()
  const translateSections = useTranslate('Lk.Sections')
  const { view, error, setError, reload } = useOwnerResume(userId, cvId)
  const { section, changeSection, openEditTab, isDirtyRef } = useResumeSections()

  if (error) return <LoadFailure error={error} />
  if (!view || !userId) return <LinearProgress />

  return (
    <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale={language}>
      <Stack spacing={3} sx={{ width: '100%', minWidth: 0 }}>
        <ResumeHeader userId={userId} view={view} onReload={reload} onError={setError} />

        <Tabs
          value={section}
          onChange={(_, id) => changeSection(id)}
          variant="scrollable"
          scrollButtons="auto"
        >
          {TOP_SECTIONS.map((id) => (
            <Tab key={id} value={id} label={translateSections(id)} sx={TAB_SX} />
          ))}
        </Tabs>

        {section === 'main' && <MainTab view={view} onOpenTab={openEditTab} />}

        {section === 'edit' && (
          <ResumeEditor userId={userId} view={view} onSaved={reload} isDirtyRef={isDirtyRef} />
        )}

        {section === 'links' && <LinksTab resumeId={view.id} />}

        {section === 'analytics' && <AnalyticsTab />}
      </Stack>
    </LocalizationProvider>
  )
}

export default LkResumeEditScreen
