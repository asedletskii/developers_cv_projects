import { SectionCard } from '@features/curriculumVitae/components/editor/form/SectionCard'
import { SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useSectionsEditor } from './hooks'

export const SectionsEditor = () => {
  const translate = useTranslate('Lk.Edit.Sections')
  const editor = useSectionsEditor()

  return (
    <Stack spacing={2}>
      {editor.fields.map(({ _key, ...section }, index) => (
        <SectionCard
          key={_key}
          title={translate(section.id)}
          enabled={section.enabled}
          onToggleEnabled={() => editor.toggleEnabled(index)}
          collapsed={!!section.collapsed}
          onToggleCollapsed={() => editor.toggleCollapsed(index)}
          collapsible={SECTION_DEFINITIONS[section.id].isCollapsible}
          pinned={editor.isPinned(section.id)}
          canMoveUp={editor.canMoveUp(index)}
          canMoveDown={editor.canMoveDown(index)}
          onMoveUp={() => editor.moveUp(index)}
          onMoveDown={() => editor.moveDown(index)}
        />
      ))}
    </Stack>
  )
}
