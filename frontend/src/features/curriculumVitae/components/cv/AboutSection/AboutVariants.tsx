import { Section } from '@features/curriculumVitae/components/cv/Section'
import { iconTile, richCard } from '@features/curriculumVitae/styles'
import { ResumeLanguageItem, ResumeProfile, SectionId } from '@features/curriculumVitae/types'
import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { getPersonFIO } from '@shared/utils'
import { useFacts } from './Facts'
import { Portrait } from './Portrait'

type Props = {
  profile: ResumeProfile
  languages: ResumeLanguageItem[]
}

const useAbout = ({ profile, languages }: Props) => {
  const { name, surname, patronymic, position, photoPath, ...facts } = profile

  return {
    title: useTranslate('CurriculumVitae.About')('title'),
    fio: getPersonFIO({ name, surname, patronymic }),
    initials: `${name[0] ?? ''}${surname[0] ?? ''}`,
    position,
    photoPath,
    facts: useFacts({ ...facts, languages }),
  }
}

export const AboutRich = (props: Props) => {
  const about = useAbout(props)

  return (
    <Section id={SectionId.About} title={about.title} collapsible={false}>
      <Stack
        direction={{ xs: 'column', md: 'row' }}
        spacing={{ xs: 3, md: 5 }}
        alignItems="center"
        sx={{ p: { xs: 3, md: 5 }, ...richCard }}
      >
        <Box sx={{ p: 2 }}>
          <Portrait src={about.photoPath} alt={about.fio} initials={about.initials} size={176} />
        </Box>

        <Stack spacing={3} flex={1} minWidth={0} sx={{ position: 'relative', zIndex: 1 }}>
          <Stack data-preview="profile" spacing={0.5} textAlign={{ xs: 'center', md: 'left' }}>
            <Typography variant="h3" sx={{ fontSize: { xs: '1.75rem', md: '2.5rem' } }}>
              {about.fio}
            </Typography>

            <Typography variant="h6" sx={{ color: 'primary.main' }}>
              {about.position}
            </Typography>
          </Stack>

          <Box
            className="reveal-group"
            sx={{
              display: 'grid',
              gap: 1.5,
              gridTemplateColumns: 'repeat(auto-fill, minmax(190px, 1fr))',
            }}
          >
            {about.facts.map(({ key, icon, label, preview }) => (
              <Stack
                key={key}
                className="reveal-item"
                data-preview={preview}
                direction="row"
                spacing={1.5}
                alignItems="center"
              >
                <Stack
                  alignItems="center"
                  justifyContent="center"
                  sx={{ width: 36, height: 36, borderRadius: 1.5, flexShrink: 0, ...iconTile }}
                >
                  {icon}
                </Stack>

                <Typography variant="body2" sx={{ fontWeight: 600 }}>
                  {label}
                </Typography>
              </Stack>
            ))}
          </Box>
        </Stack>
      </Stack>
    </Section>
  )
}

export const AboutCompact = (props: Props) => {
  const about = useAbout(props)

  return (
    <Section id={SectionId.About} title={about.title} collapsible={false}>
      <Stack direction="row" spacing={2} alignItems="center">
        <Portrait src={about.photoPath} alt={about.fio} initials={about.initials} size={64} />

        <Stack spacing={0.5} minWidth={0}>
          <Typography data-preview="profile" variant="h5" sx={{ fontSize: '1.25rem' }}>
            {about.fio}
            <Typography component="span" sx={{ color: 'primary.main', fontWeight: 600 }}>
              {' · '}
              {about.position}
            </Typography>
          </Typography>

          <Box
            component="ul"
            sx={{
              m: 0,
              p: 0,
              listStyle: 'none',
              display: 'flex',
              flexWrap: 'wrap',
              color: 'text.secondary',
              typography: 'body2',
              '& li + li::before': { content: '"·"', mx: 0.75 },
            }}
          >
            {about.facts.map(({ key, label, preview }) => (
              <li key={key} data-preview={preview}>
                {label}
              </li>
            ))}
          </Box>
        </Stack>
      </Stack>
    </Section>
  )
}
