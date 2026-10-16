import { Section } from '@features/curriculumVitae/components/cv/Section'
import { Activity as ActivityData, SectionId } from '@features/curriculumVitae/types'
import { Box } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Activity } from './Activity'

type Props = {
  activities: ActivityData[]
}

export const ActivitiesSection = ({ activities }: Props) => {
  const translate = useTranslate('CurriculumVitae.Activities')

  if (!activities.length) return null

  return (
    <Section id={SectionId.Activities} title={translate('title')}>
      <Box
        className="reveal-group"
        sx={{
          display: 'grid',
          gap: 2.5,
          gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
        }}
      >
        {activities.map((activity) => (
          <Activity key={activity.id} activity={activity} />
        ))}
      </Box>
    </Section>
  )
}
