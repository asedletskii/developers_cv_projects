import { ExternalLink } from '@features/curriculumVitae/components/cv/ExternalLink'
import { Section } from '@features/curriculumVitae/components/cv/Section'
import { iconTile, richCard } from '@features/curriculumVitae/styles'
import { Course, SectionId } from '@features/curriculumVitae/types'
import { WorkspacePremiumOutlined } from '@mui/icons-material'
import { Box, Divider, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { formatMonth } from '@shared/utils'

type Props = {
  courses: Course[]
}

const useCourseText = () => {
  const translate = useTranslate('CurriculumVitae.Courses')

  return {
    title: translate('title'),
    certificate: translate('certificate'),
    meta: (course: Course) =>
      [
        course.organization,
        formatMonth(course.endDate),
        course.hours && translate('hours', { hours: course.hours }),
      ]
        .filter(Boolean)
        .join(' · '),
  }
}

export const CoursesRich = ({ courses }: Props) => {
  const text = useCourseText()

  if (!courses.length) return null

  return (
    <Section id={SectionId.Courses} title={text.title}>
      <Box
        className="reveal-group"
        sx={{
          display: 'grid',
          gap: 2.5,
          gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' },
        }}
      >
        {courses.map((course) => (
          <Stack key={course.id} className="reveal-item" spacing={2} sx={{ p: 3, ...richCard }}>
            <Stack
              alignItems="center"
              justifyContent="center"
              sx={{ width: 48, height: 48, borderRadius: 2, ...iconTile }}
            >
              <WorkspacePremiumOutlined />
            </Stack>

            <Stack spacing={0.5} flex={1}>
              <Typography variant="h6" sx={{ fontSize: '1.125rem' }}>
                {course.name}
              </Typography>

              <Typography variant="caption" sx={{ color: 'gray.light' }}>
                {text.meta(course)}
              </Typography>
            </Stack>

            {course.certificateLink && (
              <ExternalLink href={course.certificateLink}>{text.certificate}</ExternalLink>
            )}
          </Stack>
        ))}
      </Box>
    </Section>
  )
}

export const CoursesCompact = ({ courses }: Props) => {
  const text = useCourseText()

  if (!courses.length) return null

  return (
    <Section id={SectionId.Courses} title={text.title}>
      <Stack className="reveal-group" divider={<Divider flexItem />} spacing={1}>
        {courses.map((course) => (
          <Stack
            key={course.id}
            className="reveal-item"
            direction="row"
            spacing={2}
            justifyContent="space-between"
            alignItems="baseline"
          >
            <Typography sx={{ minWidth: 0 }}>
              <Typography component="span" sx={{ fontWeight: 700 }}>
                {course.name}
              </Typography>

              <Typography component="span" variant="caption" sx={{ color: 'gray.light' }}>
                {' · '}
                {text.meta(course)}
              </Typography>
            </Typography>

            {course.certificateLink && (
              <ExternalLink href={course.certificateLink}>{text.certificate}</ExternalLink>
            )}
          </Stack>
        ))}
      </Stack>
    </Section>
  )
}
