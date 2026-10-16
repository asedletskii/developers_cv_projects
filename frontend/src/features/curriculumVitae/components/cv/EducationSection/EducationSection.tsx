import { Section } from '@features/curriculumVitae/components/cv/Section'
import { Education, SectionId } from '@features/curriculumVitae/types'
import { Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { Entry } from './Entry'

type Props = {
  education: Education[]
}

export const EducationSection = ({ education }: Props) => {
  const translate = useTranslate('CurriculumVitae.Education')

  if (!education.length) return null

  return (
    <Section id={SectionId.Education} title={translate('title')}>
      <Stack className="reveal-group" spacing={2}>
        {education.map((item) => (
          <Entry key={item.id} education={item} />
        ))}
      </Stack>
    </Section>
  )
}
