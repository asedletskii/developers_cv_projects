import { ComponentType } from 'react'
import { SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { Person, SectionId } from '@features/curriculumVitae/types'
import { Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFieldArray, useFormContext } from 'react-hook-form'
import { AboutEditor } from './AboutEditor'
import { ActivitiesEditor } from './ActivitiesEditor'
import { ContactsEditor } from './ContactsEditor'
import { CoursesEditor } from './CoursesEditor'
import { EducationEditor } from './EducationEditor'
import { ExperienceEditor } from './ExperienceEditor'
import { ProjectsEditor } from './ProjectsEditor'
import { SectionCard } from './SectionCard'
import { SkillsEditor } from './SkillsEditor'
import { TechStackEditor } from './TechStackEditor'

const editors: Record<SectionId, ComponentType> = {
  [SectionId.About]: AboutEditor,
  [SectionId.TechStack]: TechStackEditor,
  [SectionId.Experience]: ExperienceEditor,
  [SectionId.Skills]: SkillsEditor,
  [SectionId.Education]: EducationEditor,
  [SectionId.Courses]: CoursesEditor,
  [SectionId.Projects]: ProjectsEditor,
  [SectionId.Activities]: ActivitiesEditor,
  [SectionId.Contacts]: ContactsEditor,
}

/**
 * Настраиваемые секции резюме: порядок массива `sections` = порядок на странице, каждая
 * включается/выключается и переставляется независимо от остальных.
 *
 * `keyName: '_key'` — без этого RHF использует поле `id` самих данных (`SectionSettings.id`,
 * это `SectionId`, а не случайный ключ) как свой служебный React-key и затирает его случайным.
 */
export const SectionsEditor = () => {
  const translate = useTranslate('Lk.Edit.Sections')
  const { control } = useFormContext<Person>()
  const { fields, move, update } = useFieldArray({ control, name: 'sections', keyName: '_key' })

  return (
    <Stack spacing={2}>
      {fields.map(({ _key, ...section }, index) => {
        const Editor = editors[section.id]

        return (
          <SectionCard
            key={_key}
            title={translate(section.id)}
            enabled={section.enabled}
            onToggleEnabled={() => update(index, { ...section, enabled: !section.enabled })}
            collapsed={!!section.collapsed}
            onToggleCollapsed={() => update(index, { ...section, collapsed: !section.collapsed })}
            collapsible={SECTION_DEFINITIONS[section.id].isCollapsible}
            canMoveUp={index > 0}
            canMoveDown={index < fields.length - 1}
            onMoveUp={() => move(index, index - 1)}
            onMoveDown={() => move(index, index + 1)}
          >
            <Editor />
          </SectionCard>
        )
      })}
    </Stack>
  )
}
