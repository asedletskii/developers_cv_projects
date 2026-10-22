import { Section } from '@features/curriculumVitae/components/cv/Section'
import { useTechnologyFilter } from '@features/curriculumVitae/hooks'
import { iconTile, richCard } from '@features/curriculumVitae/styles'
import { SectionId, Technology } from '@features/curriculumVitae/types'
import { Box, ButtonBase, Container, Typography } from '@mui/material'
import { alpha } from '@mui/material/styles'
import { useTranslate } from '@shared/hooks'

type Props = {
  technologies: Technology[]
}

const useTechList = () => {
  const filter = useTechnologyFilter()
  const filterLabel = useTranslate('CurriculumVitae.Filter')

  return {
    title: useTranslate('CurriculumVitae.TechStack')('title'),
    filter,
    select: (name: string) => filterLabel('select', { name }),
  }
}

export const TechStackRich = ({ technologies }: Props) => {
  const { title, filter, select } = useTechList()

  if (!technologies.length) return null

  return (
    <Container maxWidth="lg" sx={{ position: 'relative' }}>
      <Section id={SectionId.TechStack} title={title} collapsible={false}>
        <Box
          className="reveal-group"
          sx={{
            display: 'grid',
            gap: 1.5,
            gridTemplateColumns: {
              xs: 'repeat(2, 1fr)',
              sm: 'repeat(3, 1fr)',
              md: 'repeat(5, 1fr)',
            },
          }}
        >
          {technologies.map(({ name }) => {
            const selected = filter.isSelected(name)

            return (
              <ButtonBase
                key={name}
                className="reveal-item"
                onClick={() => filter.toggle(name)}
                aria-label={select(name)}
                aria-pressed={selected}
                sx={{
                  ...richCard,
                  p: 1.5,
                  gap: 1.5,
                  justifyContent: 'flex-start',
                  textAlign: 'left',
                  borderColor: selected ? 'primary.main' : 'divider',
                  boxShadow: ({ palette }) =>
                    selected ? `0 0 0 3px ${alpha(palette.primary.main, 0.3)}` : undefined,
                }}
              >
                <Box
                  sx={{
                    display: 'grid',
                    placeItems: 'center',
                    width: 40,
                    height: 40,
                    borderRadius: 2,
                    flexShrink: 0,
                    fontWeight: 800,
                    textTransform: 'uppercase',
                    ...iconTile,
                  }}
                >
                  {name.slice(0, 2)}
                </Box>

                <Typography sx={{ fontWeight: 600, minWidth: 0, overflowWrap: 'anywhere' }}>
                  {name}
                </Typography>
              </ButtonBase>
            )
          })}
        </Box>
      </Section>
    </Container>
  )
}

export const TechStackCompact = ({ technologies }: Props) => {
  const { title, filter, select } = useTechList()

  if (!technologies.length) return null

  return (
    <Container maxWidth="lg" sx={{ position: 'relative' }}>
      <Section id={SectionId.TechStack} title={title} collapsible={false}>
        <Box
          component="ul"
          sx={{
            m: 0,
            p: 0,
            listStyle: 'none',
            display: 'flex',
            flexWrap: 'wrap',
            color: 'text.secondary',
            '& li + li::before': { content: '"·"', mx: 1 },
          }}
        >
          {technologies.map(({ name }) => {
            const selected = filter.isSelected(name)

            return (
              <li key={name}>
                <ButtonBase
                  onClick={() => filter.toggle(name)}
                  aria-label={select(name)}
                  aria-pressed={selected}
                  sx={{
                    typography: 'body2',
                    borderRadius: 0.5,
                    fontWeight: selected ? 700 : 400,
                    color: selected ? 'primary.main' : 'inherit',
                    textDecoration: selected ? 'underline' : 'none',
                    '&:hover': { color: 'text.primary' },
                  }}
                >
                  {name}
                </ButtonBase>
              </li>
            )
          })}
        </Box>
      </Section>
    </Container>
  )
}
