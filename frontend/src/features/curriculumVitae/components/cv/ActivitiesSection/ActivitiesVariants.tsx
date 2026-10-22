import { ExternalLink } from '@features/curriculumVitae/components/cv/ExternalLink'
import { Section } from '@features/curriculumVitae/components/cv/Section'
import { ACTIVITY_ICONS } from '@features/curriculumVitae/constants'
import { iconTile, richCard } from '@features/curriculumVitae/styles'
import { Activity, ActivityType, SectionId } from '@features/curriculumVitae/types'
import { Box, Divider, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { formatMonth } from '@shared/utils'

type Props = {
  activities: Activity[]
}

const useActivityText = () => {
  const translate = useTranslate('CurriculumVitae.Activities')

  return {
    title: translate('title'),
    open: translate('open'),
    type: (activity: Activity) => translate(`Types.${activity.type}`),
    date: (activity: Activity) => (activity.date ? formatMonth(activity.date) : ''),
  }
}

export const ActivitiesRich = ({ activities }: Props) => {
  const text = useActivityText()

  if (!activities.length) return null

  return (
    <Section id={SectionId.Activities} title={text.title}>
      <Box
        className="reveal-group"
        sx={{
          display: 'grid',
          gap: 2.5,
          gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' },
        }}
      >
        {activities.map((activity) => {
          const Icon = ACTIVITY_ICONS[activity.type] ?? ACTIVITY_ICONS[ActivityType.OpenSource]

          return (
            <Stack
              key={activity.id}
              className="reveal-item"
              spacing={1.5}
              sx={{ p: 3, ...richCard }}
            >
              <Stack direction="row" spacing={1.5} alignItems="center">
                <Stack
                  alignItems="center"
                  justifyContent="center"
                  sx={{ width: 44, height: 44, borderRadius: 2, flexShrink: 0, ...iconTile }}
                >
                  <Icon fontSize="small" />
                </Stack>

                <Typography variant="caption" sx={{ color: 'gray.light' }}>
                  {[text.type(activity), text.date(activity)].filter(Boolean).join(' · ')}
                </Typography>
              </Stack>

              <Typography variant="h6" sx={{ fontSize: '1.125rem' }}>
                {activity.title}
              </Typography>

              {activity.description && (
                <Typography sx={{ color: 'text.secondary' }}>{activity.description}</Typography>
              )}

              {activity.link && (
                <Box sx={{ mt: 'auto' }}>
                  <ExternalLink href={activity.link}>{text.open}</ExternalLink>
                </Box>
              )}
            </Stack>
          )
        })}
      </Box>
    </Section>
  )
}

export const ActivitiesCompact = ({ activities }: Props) => {
  const text = useActivityText()

  if (!activities.length) return null

  return (
    <Section id={SectionId.Activities} title={text.title}>
      <Stack className="reveal-group" divider={<Divider flexItem />} spacing={1}>
        {activities.map((activity) => (
          <Stack
            key={activity.id}
            className="reveal-item"
            direction="row"
            spacing={2}
            justifyContent="space-between"
            alignItems="baseline"
          >
            <Typography sx={{ minWidth: 0 }}>
              <Typography component="span" sx={{ fontWeight: 700 }}>
                {activity.title}
              </Typography>

              <Typography component="span" variant="caption" sx={{ color: 'gray.light' }}>
                {' · '}
                {[text.type(activity), text.date(activity)].filter(Boolean).join(' · ')}
              </Typography>
            </Typography>

            {activity.link && <ExternalLink href={activity.link}>{text.open}</ExternalLink>}
          </Stack>
        ))}
      </Stack>
    </Section>
  )
}
