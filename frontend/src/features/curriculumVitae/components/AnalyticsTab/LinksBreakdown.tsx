import { glassCard } from '@features/curriculumVitae/styles'
import { ResumeStatistics } from '@features/curriculumVitae/types'
import { formatVisitTime } from '@features/curriculumVitae/utils'
import { ButtonBase, Stack, Typography } from '@mui/material'
import { useLanguage, useTranslate } from '@shared/hooks'

type Props = {
  links: ResumeStatistics['links']
  onOpenLink: (linkId: string) => void
}

export const LinksBreakdown = ({ links, onOpenLink }: Props) => {
  const translate = useTranslate('Lk.Analytics.Links')
  const { language } = useLanguage()

  return (
    <Stack spacing={1} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Typography variant="h6">{translate('title')}</Typography>

      {links.map((link) => (
        <ButtonBase
          key={link.id}
          onClick={() => onOpenLink(link.id)}
          sx={{
            justifyContent: 'space-between',
            gap: 2,
            p: 1.5,
            borderRadius: 2,
            textAlign: 'left',
            '&:hover': { bgcolor: 'action.hover' },
          }}
        >
          <Stack sx={{ minWidth: 0 }}>
            <Typography noWrap sx={{ fontWeight: 600 }}>
              {link.name}
            </Typography>
            <Typography variant="caption" sx={{ color: 'text.secondary' }}>
              {link.lastViewAt
                ? translate('last', {
                    time: formatVisitTime(link.lastViewAt, language, { isCapitalized: false }),
                  })
                : translate('never')}
            </Typography>
          </Stack>

          <Typography sx={{ whiteSpace: 'nowrap' }}>
            {link.views} · {link.reads} {translate('reads')}
          </Typography>
        </ButtonBase>
      ))}
    </Stack>
  )
}
