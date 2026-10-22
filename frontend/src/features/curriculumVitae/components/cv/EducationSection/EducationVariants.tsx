import { Section } from '@features/curriculumVitae/components/cv/Section'
import { iconTile, richCard } from '@features/curriculumVitae/styles'
import { Education, SectionId } from '@features/curriculumVitae/types'
import { SchoolOutlined } from '@mui/icons-material'
import { Chip, Divider, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  education: Education[]
}

const useEducationText = () => {
  const translate = useTranslate('CurriculumVitae.Education')

  return {
    title: translate('title'),
    degree: (item: Education) =>
      [translate(`Levels.level${item.grade}`), item.programName, item.profileName]
        .filter(Boolean)
        .join(' · '),
    place: (item: Education) => [item.faculty, item.location].filter(Boolean).join(' · '),
    when: (item: Education) =>
      item.endDate
        ? translate('graduated', { year: item.endDate })
        : item.course
          ? translate('course', { course: item.course })
          : '',
  }
}

export const EducationRich = ({ education }: Props) => {
  const text = useEducationText()

  if (!education.length) return null

  return (
    <Section id={SectionId.Education} title={text.title}>
      <Stack className="reveal-group" spacing={2.5}>
        {education.map((item) => (
          <Stack
            key={item.id}
            className="reveal-item"
            direction="row"
            spacing={2.5}
            alignItems="center"
            sx={{ p: { xs: 2.5, md: 3.5 }, ...richCard }}
          >
            <Stack
              alignItems="center"
              justifyContent="center"
              sx={{ width: 56, height: 56, borderRadius: 2, flexShrink: 0, ...iconTile }}
            >
              <SchoolOutlined />
            </Stack>

            <Stack spacing={0.5} flex={1} minWidth={0}>
              <Typography variant="h5" sx={{ fontSize: { xs: '1.125rem', md: '1.375rem' } }}>
                {item.university}
              </Typography>

              <Typography sx={{ color: 'primary.main', fontWeight: 600 }}>
                {text.degree(item)}
              </Typography>

              <Typography variant="caption" sx={{ color: 'gray.light' }}>
                {text.place(item)}
              </Typography>

              {text.when(item) && (
                <Chip
                  size="small"
                  variant="outlined"
                  label={text.when(item)}
                  sx={{ alignSelf: 'flex-start', mt: 0.5 }}
                />
              )}
            </Stack>
          </Stack>
        ))}
      </Stack>
    </Section>
  )
}

export const EducationCompact = ({ education }: Props) => {
  const text = useEducationText()

  if (!education.length) return null

  return (
    <Section id={SectionId.Education} title={text.title}>
      <Stack className="reveal-group" divider={<Divider flexItem />} spacing={1.25}>
        {education.map((item) => (
          <Stack
            key={item.id}
            className="reveal-item"
            direction="row"
            spacing={2}
            justifyContent="space-between"
            alignItems="baseline"
          >
            <Typography sx={{ fontWeight: 700, minWidth: 0 }}>
              {item.university}
              <Typography component="span" sx={{ color: 'text.secondary', fontWeight: 400 }}>
                {' · '}
                {text.degree(item)}
              </Typography>
            </Typography>

            <Typography
              variant="caption"
              sx={{ color: 'gray.light', flexShrink: 0, whiteSpace: 'nowrap' }}
            >
              {text.when(item)}
            </Typography>
          </Stack>
        ))}
      </Stack>
    </Section>
  )
}
