import { AccentTheme } from '@features/curriculumVitae/components/AccentTheme'
import { SectionCollapseProvider } from '@features/curriculumVitae/components/SectionCollapseProvider'
import { TechnologyFilterProvider } from '@features/curriculumVitae/components/TechnologyFilterProvider'
import { glassCard } from '@features/curriculumVitae/styles'
import { SectionId } from '@features/curriculumVitae/types'
import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useSectionPreview } from './hooks'

type Props = {
  section: SectionId
  highlight?: string
}

export const SectionPreview = ({ section, highlight }: Props) => {
  const translate = useTranslate('Lk.Edit')
  const { resume, definition, hasContent } = useSectionPreview(section)

  return (
    <AccentTheme colors={resume.overview.theme}>
      <SectionCollapseProvider initiallyCollapsed={[]}>
        <TechnologyFilterProvider projects={resume.projects ?? []}>
          <Stack spacing={1}>
            <Typography variant="caption" sx={{ color: 'text.secondary' }}>
              {translate('preview')}
            </Typography>

            <Box
              aria-hidden
              sx={{
                maxHeight: 360,
                overflow: 'auto',
                isolation: 'isolate',
                p: 2,
                ...glassCard,
                ...(highlight && {
                  [`& [data-preview]:not([data-preview="${highlight}"])`]: { opacity: 0.35 },
                }),
              }}
            >
              {hasContent ? (
                <Box sx={{ zoom: 0.6, pointerEvents: 'none' }}>{definition.render(resume)}</Box>
              ) : (
                <Typography sx={{ color: 'text.secondary' }}>
                  {translate('previewEmpty')}
                </Typography>
              )}
            </Box>
          </Stack>
        </TechnologyFilterProvider>
      </SectionCollapseProvider>
    </AccentTheme>
  )
}
