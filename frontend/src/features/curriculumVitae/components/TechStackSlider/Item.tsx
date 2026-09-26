import { useTechnologyFilter } from '@features/curriculumVitae/hooks'
import { Box, ButtonBase } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  name: string
  isClone: boolean
}

export const Item = ({ name, isClone }: Props) => {
  const translate = useTranslate('CurriculumVitae.Filter')
  const { selected, isSelected, toggle } = useTechnologyFilter()

  const isActive = isSelected(name)
  const isMuted = selected !== null && !isActive

  return (
    <ButtonBase
      onClick={() => toggle(name)}
      tabIndex={isClone ? -1 : 0}
      aria-hidden={isClone}
      aria-pressed={isActive}
      aria-label={translate('select', { name })}
      sx={{
        typography: 'h3',
        whiteSpace: 'nowrap',
        color: 'primary.contrastText',
        fontSize: { xs: '1.5rem', md: '2rem' },
        '&::after': { content: '"✦"', mx: { xs: 3, md: 5 }, opacity: 0.5 },
      }}
    >
      <Box
        component="span"
        sx={{
          opacity: isMuted ? 0.45 : 1,
          textDecoration: isActive ? 'underline' : 'none',
          textDecorationThickness: 3,
          textUnderlineOffset: 8,
          transition: 'opacity .25s',
        }}
      >
        {name}
      </Box>
    </ButtonBase>
  )
}
