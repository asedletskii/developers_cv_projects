import { SectionId } from '@features/curriculumVitae/types'
import { ButtonBase, Stack, Typography } from '@mui/material'
import { getContentId, getTitleId } from '@features/curriculumVitae/utils/getSectionElementId'
import { Chevron } from './Chevron'
import { Line } from './Line'

type Props = {
  id: SectionId
  title: string
  isCollapsible: boolean
  isCollapsed: boolean
  onToggle(): void
}

/** Заголовок секции: у сворачиваемой секции это кнопка аккордеона. */
export const Header = ({ id, title, isCollapsible, isCollapsed, onToggle }: Props) => (
  <Typography
    id={getTitleId(id)}
    component="h2"
    variant="h4"
    sx={{ fontSize: { xs: '1.5rem', md: '1.75rem' } }}
  >
    {isCollapsible ? (
      <ButtonBase
        aria-expanded={!isCollapsed}
        aria-controls={getContentId(id)}
        onClick={onToggle}
        sx={{
          width: '100%',
          gap: 2,
          justifyContent: 'flex-start',
          textAlign: 'left',
          font: 'inherit',
          letterSpacing: 'inherit',
          color: 'inherit',
          borderRadius: 2,
          '&:focus-visible': {
            outline: '2px solid',
            outlineColor: 'primary.main',
            outlineOffset: 4,
          },
        }}
      >
        <span>{title}</span>

        <Line />

        <Chevron isCollapsed={isCollapsed} />
      </ButtonBase>
    ) : (
      <Stack direction="row" alignItems="center" spacing={2} component="span">
        <span>{title}</span>

        <Line />
      </Stack>
    )}
  </Typography>
)
