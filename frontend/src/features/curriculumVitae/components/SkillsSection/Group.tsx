import { glassCard } from '@features/curriculumVitae/styles'
import { Chip, Stack, Typography } from '@mui/material'

type Props = {
  title: string
  color: 'primary' | 'secondary'
  skills: { id: string; label: string }[]
}

/** Карточка с группой навыков: заголовок и чипы, которые появляются каскадом. */
export const Group = ({ title, color, skills }: Props) => (
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
