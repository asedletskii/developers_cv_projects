import { ComponentType } from 'react'
import { ResumeEditorTab, SectionId } from '@features/curriculumVitae/types'
import { AboutEditor } from './editors/AboutEditor'
import { ActivitiesEditor } from './editors/ActivitiesEditor'
import { ContactsEditor } from './editors/ContactsEditor'
import { CoursesEditor } from './editors/CoursesEditor'
import { EducationEditor } from './editors/EducationEditor'
import { ExperienceEditor } from './editors/ExperienceEditor'
import { LanguagesEditor } from './editors/LanguagesEditor'
import { OverviewEditor } from './editors/OverviewEditor'
import { ProjectsEditor } from './editors/ProjectsEditor'
import { SkillsEditor } from './editors/SkillsEditor'
import { TechStackEditor } from './editors/TechStackEditor'

export type EditorTab = {
  id: ResumeEditorTab
  Editor: ComponentType
  preview?: SectionId
  hasStylePicker?: boolean
  highlight?: string
}

export const EDITOR_TABS: EditorTab[] = [
  { id: 'overview', Editor: OverviewEditor },
  { id: 'profile', Editor: AboutEditor, preview: SectionId.About, highlight: 'profile' },
  {
    id: 'languages',
    Editor: LanguagesEditor,
    preview: SectionId.About,
    highlight: 'languages',
    hasStylePicker: false,
  },
  { id: 'technologies', Editor: TechStackEditor, preview: SectionId.TechStack },
  { id: 'experience', Editor: ExperienceEditor, preview: SectionId.Experience },
  { id: 'skills', Editor: SkillsEditor, preview: SectionId.Skills },
  { id: 'education', Editor: EducationEditor, preview: SectionId.Education },
  { id: 'courses', Editor: CoursesEditor, preview: SectionId.Courses },
  { id: 'projects', Editor: ProjectsEditor, preview: SectionId.Projects },
  { id: 'activities', Editor: ActivitiesEditor, preview: SectionId.Activities },
  { id: 'contacts', Editor: ContactsEditor, preview: SectionId.Contacts },
]

export const TAB_SX = { textTransform: 'none', minWidth: 0, px: { xs: 1.5, sm: 2 } } as const
