import { ElementType } from 'react'
import { SectionStyle } from '@features/curriculumVitae/sections/sectionStyles'
import { glassCard, iconTile, richCard } from '@features/curriculumVitae/styles'
import { CodeOutlined, GroupsOutlined } from '@mui/icons-material'
import { Box, Chip, Stack, Typography } from '@mui/material'
import { alpha } from '@mui/material/styles'

type Props = {
  title: string
  color: 'primary' | 'secondary'
  skills: { id: string; label: string }[]
  variant?: SectionStyle
}

const ICONS: Record<Props['color'], ElementType> = {
  primary: CodeOutlined,
  secondary: GroupsOutlined,
}

export const Group = ({ title, color, skills, variant = 'first' }: Props) => {
  if (variant === 'third') {
    return (
      <Typography className="reveal-item" variant="body2">
        <Typography component="span" variant="inherit" sx={{ fontWeight: 700 }}>
          {title}:{' '}
        </Typography>

        <Typography component="span" variant="inherit" sx={{ color: 'text.secondary' }}>
          {skills.map(({ label }) => label).join(' · ')}
        </Typography>
      </Typography>
    )
  }

  if (variant === 'second') {
    const Icon = ICONS[color]

    return (
      <Stack spacing={2.5} sx={{ flex: 1, p: 3, ...richCard }}>
        <Stack direction="row" spacing={1.5} alignItems="center">
          <Stack
            alignItems="center"
            justifyContent="center"
            sx={{ width: 40, height: 40, borderRadius: 1.5, ...iconTile }}
          >
            <Icon fontSize="small" />
          </Stack>

          <Typography variant="h6" sx={{ fontWeight: 700, flex: 1 }}>
            {title}
          </Typography>

          <Typography variant="caption" sx={{ color: 'gray.light' }}>
            {skills.length}
          </Typography>
        </Stack>

        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, position: 'relative', zIndex: 1 }}>
          {skills.map(({ id, label }) => (
            <Chip
              key={id}
              className="reveal-item"
              label={label}
              sx={({ palette }) => ({
                color: 'text.primary',
                fontWeight: 600,
                bgcolor: alpha(palette[color].main, 0.16),
                border: 1,
                borderColor: alpha(palette[color].main, 0.4),
              })}
            />
          ))}
        </Box>
      </Stack>
    )
  }

  return (
    <Stack spacing={2} sx={{ flex: 1, p: 3, ...glassCard }}>
      <Typography variant="h6" sx={{ fontWeight: 700 }}>
        {title}
      </Typography>

      <Stack direction="row" flexWrap="wrap" gap={1}>
        {skills.map(({ id, label }) => (
          <Chip key={id} className="reveal-item" label={label} color={color} variant="outlined" />
        ))}
      </Stack>
    </Stack>
  )
}
