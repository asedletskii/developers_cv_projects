import { Box, Stack, Typography } from '@mui/material'

type Row = {
  key: string
  label: string
  value: number
  /** Подпись справа от числа */
  caption?: string
}

type Props = {
  rows: Row[]
  /** Значение, соответствующее полной ширине полосы. По умолчанию — максимум из строк */
  max?: number
}

/** Горизонтальные полосы с подписями — общий вид воронки, стека, устройств и источников */
export const BarList = ({ rows, max = Math.max(...rows.map((row) => row.value), 1) }: Props) => (
  <Stack spacing={1.25}>
    {rows.map(({ key, label, value, caption }) => (
      <Stack key={key} spacing={0.5}>
        <Stack direction="row" justifyContent="space-between" spacing={2}>
          <Typography variant="body2" noWrap>
            {label}
          </Typography>

          <Typography variant="body2" sx={{ whiteSpace: 'nowrap', color: 'text.secondary' }}>
            <b style={{ color: 'inherit', fontWeight: 600 }}>{value}</b>
            {caption ? ` · ${caption}` : ''}
          </Typography>
        </Stack>

        <Box sx={{ height: 8, borderRadius: 999, bgcolor: 'action.hover', overflow: 'hidden' }}>
          <Box
            sx={{
              height: '100%',
              width: `${(value / max) * 100}%`,
              borderRadius: 999,
              bgcolor: 'primary.main',
              transition: 'width .4s ease',
            }}
          />
        </Box>
      </Stack>
    ))}
  </Stack>
)
