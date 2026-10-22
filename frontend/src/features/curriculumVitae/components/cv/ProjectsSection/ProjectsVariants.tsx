import { ExternalLink } from '@features/curriculumVitae/components/cv/ExternalLink'
import { Section } from '@features/curriculumVitae/components/cv/Section'
import { TechnologyList } from '@features/curriculumVitae/components/cv/TechnologyList'
import { useTechnologyFilter } from '@features/curriculumVitae/hooks'
import { dimmed, richCard } from '@features/curriculumVitae/styles'
import { Resume, SectionId } from '@features/curriculumVitae/types'
import { safeExternalUrl } from '@features/curriculumVitae/utils/safeExternalUrl'
import { Box, Chip, Divider, Stack, Typography } from '@mui/material'
import { alpha } from '@mui/material/styles'
import { useTranslate } from '@shared/hooks'

type Props = {
  projects: Resume['projects']
}

export const ProjectsRich = ({ projects }: Props) => {
  const translate = useTranslate('CurriculumVitae.Projects')
  const { isDimmed } = useTechnologyFilter()

  if (!projects.length) return null

  return (
    <Section id={SectionId.Projects} title={translate('title')}>
      <Box
        className="reveal-group"
        sx={{
          display: 'grid',
          gap: 2.5,
          gridTemplateColumns: { xs: '1fr', md: 'repeat(2, 1fr)' },
        }}
      >
        {projects.map((project) => {
          const cover = safeExternalUrl(project.screenshots?.[0])

          return (
            <Stack
              key={project.id}
              className="reveal-item"
              sx={{ ...richCard, ...dimmed(isDimmed(project.stack)) }}
            >
              {cover ? (
                <Box
                  component="img"
                  src={cover}
                  alt={project.name}
                  loading="lazy"
                  sx={{ width: '100%', aspectRatio: '16 / 9', objectFit: 'cover' }}
                />
              ) : (
                <Box
                  aria-hidden
                  sx={{
                    display: 'grid',
                    placeItems: 'center',
                    width: '100%',
                    height: 96,
                    fontSize: '2.5rem',
                    fontWeight: 800,
                    color: 'primary.main',
                    background: ({ palette }) =>
                      `linear-gradient(135deg, ${alpha(palette.primary.main, 0.28)}, ${alpha(palette.primary.dark, 0.08)})`,
                  }}
                >
                  {project.name[0]}
                </Box>
              )}

              <Stack spacing={1.5} flex={1} sx={{ p: 3, position: 'relative', zIndex: 1 }}>
                <Typography variant="h5" sx={{ fontSize: '1.25rem' }}>
                  {project.name}
                </Typography>

                {project.role && (
                  <Chip
                    size="small"
                    variant="outlined"
                    label={project.role}
                    sx={{ alignSelf: 'flex-start' }}
                  />
                )}

                {project.description && (
                  <Typography sx={{ color: 'text.secondary' }}>{project.description}</Typography>
                )}

                <Box sx={{ mt: 'auto', pt: 1 }}>
                  <TechnologyList technologies={project.stack} />
                </Box>

                <Stack direction="row" spacing={2} flexWrap="wrap">
                  {project.link && (
                    <ExternalLink href={project.link} kind="project">
                      {translate('open')}
                    </ExternalLink>
                  )}

                  {project.gitLink && (
                    <ExternalLink href={project.gitLink} kind="project">
                      {translate('code')}
                    </ExternalLink>
                  )}
                </Stack>
              </Stack>
            </Stack>
          )
        })}
      </Box>
    </Section>
  )
}

/** Компактный вид: название и описание в одну строку, ссылка справа, стек текстом */
export const ProjectsCompact = ({ projects }: Props) => {
  const translate = useTranslate('CurriculumVitae.Projects')
  const { isDimmed } = useTechnologyFilter()

  if (!projects.length) return null

  return (
    <Section id={SectionId.Projects} title={translate('title')}>
      <Stack className="reveal-group" divider={<Divider flexItem />} spacing={1.25}>
        {projects.map((project) => (
          <Box key={project.id} className="reveal-item" sx={dimmed(isDimmed(project.stack))}>
            <Stack direction="row" spacing={2} justifyContent="space-between" alignItems="baseline">
              <Typography noWrap sx={{ minWidth: 0 }}>
                <Typography component="span" sx={{ fontWeight: 700 }}>
                  {project.name}
                </Typography>
                {project.description && (
                  <Typography component="span" sx={{ color: 'text.secondary' }}>
                    {' — '}
                    {project.description}
                  </Typography>
                )}
              </Typography>

              {project.link && (
                <ExternalLink href={project.link} kind="project">
                  {translate('open')}
                </ExternalLink>
              )}
            </Stack>

            {!!project.stack.length && (
              <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                {project.stack.map(({ name }) => name).join(', ')}
              </Typography>
            )}
          </Box>
        ))}
      </Stack>
    </Section>
  )
}
