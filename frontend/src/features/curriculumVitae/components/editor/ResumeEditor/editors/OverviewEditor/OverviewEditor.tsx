import { HeroEditor } from '@features/curriculumVitae/components/editor/ResumeEditor/editors/HeroEditor'
import { SectionsEditor } from '@features/curriculumVitae/components/editor/ResumeEditor/editors/SectionsEditor'
import { ThemeEditor } from '@features/curriculumVitae/components/editor/ResumeEditor/editors/ThemeEditor'
import { Stack } from '@mui/material'

/** Вкладка «Обзор»: цвета, приветствие, доступность и порядок/видимость секций. */
export const OverviewEditor = () => (
  <Stack spacing={3}>
    <ThemeEditor />
    <HeroEditor />
    <SectionsEditor />
  </Stack>
)
