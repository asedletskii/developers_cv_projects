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

  /** Что выделить в секции рамкой — редактируемый на этой вкладке участок (`data-preview` в разметке) */
  highlight?: string
}

/**
 * Живой предпросмотр секции: тот же рендер, что на странице резюме (`SECTION_DEFINITIONS`),
 * но данные — из формы «как сейчас набрано», без сохранения. Уменьшен через `zoom`, не кликается;
 * длинные секции прокручиваются внутри рамки.
 */
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
                // Свой стек отрисовки: без этого декоративные абсолютно спозиционированные
                // элементы внутри (например, кольцо вокруг аватара) могут прорисоваться поверх
                // диалогов и другого UI за пределами этого блока
                isolation: 'isolate',
                p: 2,
                ...glassCard,
                ...(highlight && {
                  // Всё, что на этой вкладке не редактируется, приглушаем — остаётся видно, где правка
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
