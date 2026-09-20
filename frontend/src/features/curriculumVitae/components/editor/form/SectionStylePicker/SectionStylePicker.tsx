import { useState } from 'react'
import { DEFAULT_SECTIONS } from '@features/curriculumVitae/sections'
import {
  resolveSectionStyle,
  SECTION_STYLES,
  SectionStyle,
} from '@features/curriculumVitae/sections/sectionStyles'
import { Resume, SectionId } from '@features/curriculumVitae/types'
import { saveResumeSection } from '@features/curriculumVitae/utils'
import { withPinnedAbout } from '@features/curriculumVitae/utils/withPinnedAbout'
import { Stack, ToggleButton, ToggleButtonGroup, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useFormContext, useWatch } from 'react-hook-form'

type Props = {
  view: Resume
  section: SectionId
  isDirty: boolean
  onSaved: () => Promise<void>
}

export const SectionStylePicker = ({ view, section, isDirty, onSaved }: Props) => {
  const translate = useTranslate('Lk.Edit.Style')
  const { control, setValue } = useFormContext<Resume>()
  const sections = useWatch({ control, name: 'overview.sections' }) ?? []
  const [status, setStatus] = useState<'idle' | 'saving' | 'error'>('idle')

  const current = resolveSectionStyle(sections.find(({ id }) => id === section)?.style)

  const choose = async (style: SectionStyle) => {
    const next = withPinnedAbout(sections.length ? sections : DEFAULT_SECTIONS).map((item) =>
      item.id === section ? { ...item, style } : item,
    )

    setValue('overview.sections', next)
    setStatus('saving')

    try {
      await saveResumeSection(view, 'overview', {
        ...view,
        overview: { ...view.overview, sections: next },
      })
      await onSaved()
      setStatus('idle')
    } catch {
      setValue('overview.sections', sections)
      setStatus('error')
    }
  }

  return (
    <Stack spacing={0.75}>
      <Stack direction="row" spacing={2} alignItems="center" flexWrap="wrap" useFlexGap>
        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
          {translate('title')}
        </Typography>

        <ToggleButtonGroup
          exclusive
          size="small"
          value={current}
          disabled={isDirty || status === 'saving'}
          onChange={(_, next: SectionStyle | null) => next && choose(next)}
        >
          {SECTION_STYLES.map((style) => (
            <ToggleButton key={style} value={style} sx={{ textTransform: 'none', px: 2 }}>
              {translate(style)}
            </ToggleButton>
          ))}
        </ToggleButtonGroup>
      </Stack>

      {(isDirty || status === 'error') && (
        <Typography
          variant="caption"
          sx={{ color: status === 'error' ? 'error.main' : 'text.secondary' }}
        >
          {translate(status === 'error' ? 'error' : 'saveFirst')}
        </Typography>
      )}
    </Stack>
  )
}
