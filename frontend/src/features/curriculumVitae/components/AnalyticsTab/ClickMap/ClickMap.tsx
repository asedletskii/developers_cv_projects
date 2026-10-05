import { useRef, useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { useAnalyticsRequest } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { ClickBreakpoint, SectionId } from '@features/curriculumVitae/types'
import {
  Box,
  Button,
  Chip,
  Stack,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useElementWidth, useTranslate } from '@shared/hooks'
import { createPortal } from 'react-dom'
import { generatePath } from 'react-router'
import { BREAKPOINT_WIDTH, BREAKPOINTS, VIEWPORT_HEIGHT } from './constants'
import { HEAT_LEGEND_GRADIENT } from './heatColor'
import { HeatLayer } from './HeatLayer'
import { useEmbeddedSections } from './useEmbeddedSections'

type Props = {
  linkId: string
  resumeId: string
}

export const ClickMap = ({ linkId, resumeId }: Props) => {
  const translate = useTranslate('Lk.Analytics.ClickMap')
  const translateSection = useTranslate('Lk.Edit.Sections')
  const [breakpoint, setBreakpoint] = useState<ClickBreakpoint>('lg')
  const [includeStale, setIncludeStale] = useState(false)
  const hostRef = useRef<HTMLDivElement>(null)
  const hostWidth = useElementWidth(hostRef as React.RefObject<HTMLElement>)
  const { iframeRef, onLoad, sections } = useEmbeddedSections()

  const width = BREAKPOINT_WIDTH[breakpoint]
  // Предпросмотр верстается в натуральную ширину брейкпоинта и уменьшается под доступное место
  const scale = hostWidth ? Math.min(1, hostWidth / width) : 1
  const sectionIds = sections.map(({ id }) => id)

  const { data: grids } = useAnalyticsRequest(
    (signal) =>
      Promise.all(
        sectionIds.map((id) =>
          curriculumVitaeApi.getLinkClicks(linkId, id, breakpoint, includeStale, signal),
        ),
      ),
    [linkId, breakpoint, includeStale, sectionIds.join(',')],
  )

  const gridBySection = new Map((grids ?? []).map((grid) => [grid.section, grid]))
  const max = Math.max(1, ...(grids ?? []).flatMap(({ cells }) => cells.map(({ count }) => count)))
  const staleTotal = (grids ?? []).reduce((sum, { staleTotal: stale }) => sum + stale, 0)
  const hasClicks = (grids ?? []).some(({ total }) => total > 0)
  // Где кликали больше всего — из тех же сеток, отдельная ручка не нужна
  const topSections = [...(grids ?? [])]
    .filter(({ total }) => total > 0)
    .sort((a, b) => b.total - a.total)
    .slice(0, 3)

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack direction="row" justifyContent="space-between" alignItems="center" flexWrap="wrap">
        <Typography variant="h6">{translate('title')}</Typography>

        <ToggleButtonGroup
          exclusive
          size="small"
          value={breakpoint}
          onChange={(_, value: ClickBreakpoint | null) => value && setBreakpoint(value)}
        >
          {BREAKPOINTS.map((key) => (
            <ToggleButton key={key} value={key} sx={{ textTransform: 'none', px: 1.5 }}>
              {translate(`breakpoint.${key}`)}
            </ToggleButton>
          ))}
        </ToggleButtonGroup>
      </Stack>

      <Typography variant="caption" sx={{ color: 'text.secondary' }}>
        {translate('hint')}
      </Typography>

      {topSections.length > 0 && (
        <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap" useFlexGap>
          <Typography variant="caption" sx={{ color: 'text.secondary' }}>
            {translate('top')}
          </Typography>

          {topSections.map(({ section, total }) => (
            <Chip
              key={section}
              size="small"
              label={translate('topItem', {
                section: translateSection(section as SectionId),
                count: total,
              })}
            />
          ))}
        </Stack>
      )}

      <Stack direction="row" spacing={1} alignItems="center">
        <Typography variant="caption">{translate('legendLow')}</Typography>
        <Box sx={{ width: 120, height: 8, borderRadius: 999, background: HEAT_LEGEND_GRADIENT }} />
        <Typography variant="caption">{translate('legendHigh')}</Typography>
      </Stack>

      {(staleTotal > 0 || includeStale) && (
        <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap" useFlexGap>
          <Typography variant="caption" sx={{ color: 'text.secondary' }}>
            {includeStale
              ? translate('staleIncluded')
              : translate('staleHidden', { count: staleTotal })}
          </Typography>

          <Button
            size="small"
            onClick={() => setIncludeStale((current) => !current)}
            sx={{ textTransform: 'none' }}
          >
            {translate(includeStale ? 'staleHide' : 'staleShow')}
          </Button>
        </Stack>
      )}

      {sections.length > 0 && grids && !hasClicks && (
        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
          {translate('empty')}
        </Typography>
      )}

      <Box
        ref={hostRef}
        sx={{
          position: 'relative',
          width: '100%',
          height: VIEWPORT_HEIGHT * scale,
          overflow: 'hidden',
          borderRadius: 2,
          border: 1,
          borderColor: 'divider',
        }}
      >
        <Box
          sx={{
            position: 'relative',
            width: width * scale,
            height: VIEWPORT_HEIGHT * scale,
            mx: 'auto',
          }}
        >
          <iframe
            ref={iframeRef}
            onLoad={onLoad}
            title={translate('title')}
            src={`${generatePath(ROUTES.lkResumePreview, { cvId: resumeId })}?embed=1`}
            style={{
              width,
              height: VIEWPORT_HEIGHT,
              border: 0,
              display: 'block',
              transform: `scale(${scale})`,
              transformOrigin: 'top left',
            }}
          />

          {sections.map(({ id, host }) => {
            const grid = gridBySection.get(id)

            return (
              grid &&
              createPortal(
                <HeatLayer
                  title={translateSection(id as SectionId)}
                  grid={grid}
                  max={max}
                  scale={scale}
                />,
                host,
                id,
              )
            )
          })}
        </Box>
      </Box>
    </Stack>
  )
}
