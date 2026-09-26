import { ExternalLink } from '@features/curriculumVitae/components/ExternalLink'
import { IconTile } from '@features/curriculumVitae/components/IconTile'
import { ACTIVITY_ICONS } from '@features/curriculumVitae/constants'
import { glassCard } from '@features/curriculumVitae/styles'
import { Activity as ActivityData, ActivityType } from '@features/curriculumVitae/types'
import { Chip, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { formatMonth } from '@shared/utils'

type Props = {
  activity: ActivityData
}

export const Activity = ({ activity }: Props) => {
  const translate = useTranslate('CurriculumVitae.Activities')

  const Icon = ACTIVITY_ICONS[activity.type] ?? ACTIVITY_ICONS[ActivityType.OpenSource]

  return (
    <Stack
      className="reveal-item"
      spacing={1.5}
      sx={{
        p: 3,
        ...glassCard,
        transition: 'box-shadow .25s, border-color .25s',
        '&:hover': { boxShadow: 2, borderColor: 'primary.light' },
      }}
    >
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <IconTile size={44}>
          <Icon />
        </IconTile>

        <Chip size="small" label={translate(`Types.${activity.type}`)} />
      </Stack>

      <Stack spacing={0.5} flex={1}>
        <Typography variant="h6" sx={{ fontWeight: 700 }}>
          {activity.title}
        </Typography>

        {activity.date && (
          <Typography variant="caption" sx={{ color: 'gray.light' }}>
            {formatMonth(activity.date)}
          </Typography>
        )}

        {activity.description && (
          <Typography sx={{ color: 'text.secondary' }}>{activity.description}</Typography>
        )}
      </Stack>

      {activity.link && <ExternalLink href={activity.link}>{translate('open')}</ExternalLink>}
    </Stack>
  )
}
