import { Section } from '@features/curriculumVitae/components/cv/Section'
import { glassCard } from '@features/curriculumVitae/styles'
import { ResumeLanguageItem, ResumeProfile, SectionId } from '@features/curriculumVitae/types'
import { Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { getPersonFIO } from '@shared/utils'
import { Facts } from './Facts'
import { Portrait } from './Portrait'

type Props = {
  profile: ResumeProfile
  languages: ResumeLanguageItem[]
}

export const AboutSection = ({ profile, languages }: Props) => {
  const { name, surname, patronymic, position, photoPath, ...facts } = profile
  const translate = useTranslate('CurriculumVitae.About')

  const fio = getPersonFIO({ name, surname, patronymic })

  return (
    <Section id={SectionId.About} title={translate('title')} collapsible={false}>
      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        spacing={{ xs: 3, sm: 4 }}
        alignItems={{ xs: 'center', sm: 'flex-start' }}
        sx={{ p: { xs: 3, sm: 4 }, ...glassCard }}
      >
        <Portrait src={photoPath} alt={fio} initials={`${name[0] ?? ''}${surname[0] ?? ''}`} />

        <Stack spacing={2} sx={{ minWidth: 0 }}>
          <Stack
            data-preview="profile"
            spacing={0.5}
            alignItems={{ xs: 'center', sm: 'flex-start' }}
          >
            <Typography variant="h4" sx={{ fontSize: { xs: '1.5rem', sm: '1.75rem' } }}>
              {fio}
            </Typography>

            <Typography variant="h6" sx={{ color: 'primary.main' }}>
              {position}
            </Typography>
          </Stack>

          <Facts {...facts} languages={languages} />
        </Stack>
      </Stack>
    </Section>
  )
}
