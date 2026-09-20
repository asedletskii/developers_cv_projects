import { Section } from '@features/curriculumVitae/components/cv/Section'
import { Course as CourseData, SectionId } from '@features/curriculumVitae/types'
import { Box } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Course } from './Course'

type Props = {
  courses: CourseData[]
}

export const CoursesSection = ({ courses }: Props) => {
  const translate = useTranslate('CurriculumVitae.Courses')

  if (!courses.length) return null

  return (
    <Section id={SectionId.Courses} title={translate('title')}>
      <Box
        className="reveal-group"
        sx={{
          display: 'grid',
          gap: 2.5,
          gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' },
        }}
      >
        {courses.map((course) => (
          <Course key={course.id} course={course} />
        ))}
      </Box>
    </Section>
  )
}
