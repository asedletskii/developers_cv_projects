import { useId, useState } from 'react'
import { LinkVisit, SectionId, VisitEvent } from '@features/curriculumVitae/types'
import { formatDuration, formatVisitTime, isVisibleAction } from '@features/curriculumVitae/utils'
import { ExpandMoreRounded } from '@mui/icons-material'
import { Box, ButtonBase, Chip, Collapse, Stack, Typography } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'

const BRAND_NAMES: Record<string, string> = { ios: 'iOS', macos: 'macOS' }

const prettify = (name: string) => BRAND_NAMES[name] ?? name.charAt(0).toUpperCase() + name.slice(1)

type Props = {
  visit: LinkVisit
  visitorVisits?: number
}

export const Visit = ({ visit, visitorVisits = 1 }: Props) => {
  const translate = useTranslate('Lk.Analytics.Feed')
  const translateDetails = useTranslate('Lk.Analytics.Feed.Details')
  const translateDevice = useTranslate('Lk.Analytics.Device')
  const translateTraffic = useTranslate('Lk.Analytics.Feed.Traffic')
  const translateSection = useTranslate('Lk.Edit.Sections')
  const { language } = useLanguage()
  const [isOpen, setIsOpen] = useState(false)
  const panelId = useId()

  const sectionName = (id: string) => translateSection(id as SectionId)
  const browserAndOs = [visit.browser, visit.os].filter(Boolean).map(prettify).join(' ')
  const source = visit.referrerHost
    ? translate('referrer', { host: visit.referrerHost })
    : translate('direct')
  const activeTime = visit.isRead ? formatDuration(visit.activeSeconds, language) : ''

  const describeEvent = ({ type, contact, technology, target }: VisitEvent) => {
    if (type === 'contact_view') return translate('contact', { contact })
    if (type === 'tech_filter') return translate('techFilter', { technology })
    if (type === 'download') return translate('download')
    if (type === 'print') return translate('print')
    if (type === 'external_click') return translate('externalClick', { target })
    if (type === 'contact_copy') return translate('contactCopy', { contact })
    if (type === 'contact_click') return translate('contactClick', { contact })
    return ''
  }

  const getOffset = (at: string) => {
    const seconds = Math.round(
      (new Date(at).getTime() - new Date(visit.startedAt).getTime()) / 1000,
    )
    return Number.isFinite(seconds)
      ? translateDetails('offset', { time: formatDuration(Math.max(seconds, 0), language) })
      : ''
  }

  const actions = visit.events
    .filter(isVisibleAction)
    .map((event) => ({ event, label: describeEvent(event) }))

  const isHuman = visit.traffic === 'human'

  const fields = [
    {
      label: translateDetails('startedAt'),
      value: new Intl.DateTimeFormat(language, { dateStyle: 'medium', timeStyle: 'medium' }).format(
        new Date(visit.startedAt),
      ),
    },
    { label: translateDetails('device'), value: translateDevice(visit.device) },
    { label: translateDetails('browser'), value: browserAndOs || '—' },
    { label: translateDetails('source'), value: source },
    {
      label: translateDetails('activeTime'),
      value: visit.isRead ? activeTime : translate('bounced'),
    },
    {
      label: translateDetails('visitor'),
      value:
        translate('visitor', { n: visit.visitor.ordinal }) +
        (visit.visitor.isReturning ? ` (${translate('returning')})` : ''),
    },
    ...(visitorVisits > 1
      ? [{ label: translateDetails('visitorVisits'), value: String(visitorVisits) }]
      : []),
    ...(isHuman
      ? []
      : [{ label: translateDetails('traffic'), value: translateTraffic(visit.traffic) }]),
    {
      label: translateDetails('reached'),
      value: visit.maxSection ? sectionName(visit.maxSection) : '—',
    },
    {
      label: translateDetails('exit'),
      value: visit.exitSection ? sectionName(visit.exitSection) : '—',
    },
  ]

  return (
    <Box
      sx={{ borderRadius: 2, bgcolor: 'action.hover', opacity: visit.isRead || isOpen ? 1 : 0.65 }}
    >
      <ButtonBase
        onClick={() => setIsOpen((current) => !current)}
        aria-expanded={isOpen}
        aria-controls={panelId}
        sx={{
          display: 'flex',
          width: '100%',
          justifyContent: 'space-between',
          alignItems: 'flex-start',
          gap: 1,
          p: 1.5,
          borderRadius: 2,
          textAlign: 'left',
        }}
      >
        <Stack spacing={0.25} sx={{ minWidth: 0 }}>
          <Typography sx={{ fontWeight: 600 }}>
            {!isHuman && (
              <Chip
                size="small"
                label={translateTraffic(visit.traffic)}
                sx={{ mr: 1, height: 20 }}
              />
            )}
            {formatVisitTime(visit.startedAt, language)}
            <Typography component="span" sx={{ color: 'text.secondary', fontWeight: 400 }}>
              {` · ${[translateDevice(visit.device), browserAndOs, source, activeTime]
                .filter(Boolean)
                .join(' · ')}`}
            </Typography>
          </Typography>

          {!!actions.length && (
            <Typography variant="body2" sx={{ color: 'text.secondary' }}>
              {translateDetails('actionsCount', { count: actions.length })}
            </Typography>
          )}
        </Stack>

        <ExpandMoreRounded
          sx={{
            flexShrink: 0,
            color: 'text.secondary',
            transition: 'transform .2s',
            transform: isOpen ? 'rotate(180deg)' : 'none',
          }}
        />
      </ButtonBase>

      <Collapse in={isOpen} unmountOnExit id={panelId}>
        <Stack spacing={2} sx={{ px: 1.5, pb: 1.5 }}>
          <Box
            sx={{
              display: 'grid',
              gap: 1.5,
              gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)' },
            }}
          >
            {fields.map(({ label, value }) => (
              <Stack key={label}>
                <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                  {label}
                </Typography>
                <Typography variant="body2">{value}</Typography>
              </Stack>
            ))}
          </Box>

          <Stack spacing={0.75}>
            <Typography variant="caption" sx={{ color: 'text.secondary' }}>
              {translateDetails('actions')}
            </Typography>

            {actions.length === 0 ? (
              <Typography variant="body2">{translateDetails('noActions')}</Typography>
            ) : (
              actions.map(({ event, label }) => (
                <Stack key={event.at + event.type + label} direction="row" spacing={1.5}>
                  <Typography
                    variant="body2"
                    sx={{
                      minWidth: 136,
                      color: 'text.secondary',
                      fontVariantNumeric: 'tabular-nums',
                    }}
                  >
                    {getOffset(event.at)}
                  </Typography>
                  <Typography variant="body2">{label}</Typography>
                </Stack>
              ))
            )}
          </Stack>
        </Stack>
      </Collapse>
    </Box>
  )
}
