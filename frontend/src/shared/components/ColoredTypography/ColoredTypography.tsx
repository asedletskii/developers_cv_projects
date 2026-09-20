import { Stack, SxProps, Theme, Typography, TypographyVariant } from '@mui/material'

type Props = {
  firstPart: string
  variant?: TypographyVariant
  coloredPart: string
  color: string
  stackDirection?: 'row' | 'column'
  className?: string
  sx?: SxProps<Theme>
  coloredSx?: SxProps<Theme>
}

const toArray = (sx?: SxProps<Theme>) => (Array.isArray(sx) ? sx : [sx])

export const ColoredTypography = ({
  firstPart,
  coloredPart,
  variant,
  color,
  stackDirection = 'row',
  className,
  sx,
  coloredSx,
}: Props) => {
  return (
    <Stack
      className={className}
      direction={stackDirection}
      spacing={1}
      useFlexGap
      sx={{ flexWrap: 'wrap', justifyContent: 'center' }}
    >
      {firstPart && (
        <Typography variant={variant} sx={sx}>
          {firstPart}
        </Typography>
      )}

      {coloredPart && (
        <Typography variant={variant} sx={[{ color }, ...toArray(sx), ...toArray(coloredSx)]}>
          {coloredPart}
        </Typography>
      )}
    </Stack>
  )
}
