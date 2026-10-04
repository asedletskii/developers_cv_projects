import { LinkVisit, SectionId } from '@features/curriculumVitae/types'
import { formatDuration, formatVisitTime } from '@features/curriculumVitae/utils'
import { Stack, Typography } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'

const BRAND_NAMES: Record<string, string> = { ios: 'iOS', macos: 'macOS' }

const prettify = (name: string) => BRAND_NAMES[name] ?? name.charAt(0).toUpperCase() + name.slice(1)

type Props = {
  visit: LinkVisit
}

export const Visit = ({ visit }: Props) => {
  const translate = useTranslate('Lk.Analytics.Feed')
  const translateDevice = useTranslate('Lk.Analytics.Device')
  const translateSection = useTranslate('Lk.Edit.Sections')
  const { language } = useLanguage()

  const sectionName = (id: string) => translateSection(id as SectionId)

  const header = [
    translateDevice(visit.device),
    [visit.browser, visit.os].filter(Boolean).map(prettify).join(' '),
    visit.referrerHost ? translate('referrer', { host: visit.referrerHost }) : translate('direct'),
    visit.isRead ? formatDuration(visit.activeSeconds, language) : '',
  ].filter(Boolean)

  const path = [
    visit.maxSection && translate('reached', { section: sectionName(visit.maxSection) }),
    visit.exitSection && translate('exit', { section: sectionName(visit.exitSection) }),
  ].filter(Boolean)

  const events = visit.events.map((event) => {
    if (event.type === 'contact_view') return translate('contact', { contact: event.contact })
    if (event.type === 'tech_filter')
      return translate('techFilter', { technology: event.technology })
    if (event.type === 'download') return translate('download')
    return ''
  })

  const details = [
    translate('visitor', { n: visit.visitor.ordinal }) +
      (visit.visitor.isReturning ? ` (${translate('returning')})` : ''),
    ...(visit.isRead ? path : [translate('bounced')]),
    ...events,
  ].filter(Boolean)

  return (
    <Stack
      spacing={0.25}
      sx={{ p: 1.5, borderRadius: 2, bgcolor: 'action.hover', opacity: visit.isRead ? 1 : 0.65 }}
    >
      <Typography sx={{ fontWeight: 600 }}>
        {formatVisitTime(visit.startedAt, language)}
        <Typography component="span" sx={{ color: 'text.secondary', fontWeight: 400 }}>
          {` · ${header.join(' · ')}`}
        </Typography>
      </Typography>

      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {details.join(' · ')}
      </Typography>
    </Stack>
  )
}
