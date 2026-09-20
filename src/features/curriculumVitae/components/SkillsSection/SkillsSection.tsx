import { Stack } from '@mui/material'
import { Section } from '@features/curriculumVitae/components/Section'
import { HardSkill, SoftSkill, SectionId } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'
import { Group } from './Group'

type Props = {
  hardSkills: HardSkill[]
  softSkills: SoftSkill[]
}

export const SkillsSection = ({ hardSkills, softSkills }: Props) => {
  const translate = useTranslate('CurriculumVitae.Skills')

  const groups = [
    {
      key: 'hard',
      title: translate('hard'),
      color: 'primary' as const,
      skills: hardSkills.map((skill) => ({ id: skill, label: translate(`HardSkills.${skill}`) })),
    },
    {
      key: 'soft',
      title: translate('soft'),
      color: 'secondary' as const,
      skills: softSkills.map((skill) => ({ id: skill, label: translate(`SoftSkills.${skill}`) })),
    },
  ].filter(({ skills }) => skills.length)

  if (!groups.length) return null

  return (
    <Section id={SectionId.Skills} title={translate('title')}>
      <Stack
        className="reveal-group"
        direction={{ xs: 'column', md: 'row' }}
        spacing={2.5}
        useFlexGap
      >
        {groups.map(({ key, ...group }) => (
          <Group key={key} {...group} />
        ))}
      </Stack>
    </Section>
  )
}
