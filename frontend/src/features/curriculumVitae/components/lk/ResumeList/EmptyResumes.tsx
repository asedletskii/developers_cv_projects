import { glassCard, iconTile } from '@features/curriculumVitae/styles'
import { AddRounded, DescriptionRounded } from '@mui/icons-material'
import { Box, Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

const STEPS = ['emptyStep1', 'emptyStep2', 'emptyStep3'] as const

type Props = {
  onCreate: () => void
}

/** Первый вход: вместо пустого списка — короткий путь «резюме → разделы → ссылка» с кнопкой старта */
export const EmptyResumes = ({ onCreate }: Props) => {
  const translate = useTranslate('Lk.Resumes')

  return (
    <Stack
      spacing={3}
      alignItems="center"
      textAlign="center"
      sx={{ p: { xs: 3, md: 5 }, ...glassCard }}
    >
      <Stack
        alignItems="center"
        justifyContent="center"
        sx={{ width: 64, height: 64, borderRadius: 3, ...iconTile }}
      >
        <DescriptionRounded sx={{ fontSize: 34 }} />
      </Stack>

      <Stack spacing={1}>
        <Typography variant="h5">{translate('emptyTitle')}</Typography>
        <Typography sx={{ color: 'text.secondary', maxWidth: 440 }}>
          {translate('emptyText')}
        </Typography>
      </Stack>

      <Stack
        component="ol"
        spacing={1.5}
        sx={{ m: 0, p: 0, listStyle: 'none', textAlign: 'left', width: '100%', maxWidth: 420 }}
      >
        {STEPS.map((step, index) => (
          <Stack key={step} component="li" direction="row" spacing={1.5} alignItems="center">
            <Box
              aria-hidden
              sx={{
                display: 'grid',
                placeItems: 'center',
                width: 28,
                height: 28,
                borderRadius: '50%',
                flexShrink: 0,
                fontSize: 14,
                fontWeight: 700,
                ...iconTile,
              }}
            >
              {index + 1}
            </Box>
            <Typography>{translate(step)}</Typography>
          </Stack>
        ))}
      </Stack>

      <Button
        variant="contained"
        size="large"
        startIcon={<AddRounded />}
        onClick={onCreate}
        sx={{ borderRadius: 999, px: 4, textTransform: 'none' }}
      >
        {translate('createFirst')}
      </Button>
    </Stack>
  )
}
