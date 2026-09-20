import { PropsWithChildren, Ref } from 'react'
import { Box, Stack } from '@mui/material'

type Props = PropsWithChildren<{
  ref?: Ref<HTMLDivElement>
}>

/** Вертикальный таймлайн: линия слева «дорисовывается» вместе со скроллом (см. useTimelineAnimation). */
export const Timeline = ({ ref, children }: Props) => (
  <Stack ref={ref} spacing={4} sx={{ position: 'relative', pl: { xs: 4, sm: 6 } }}>
    <Box
      aria-hidden
      sx={{
        position: 'absolute',
        top: 8,
        bottom: 8,
        left: { xs: 9, sm: 13 },
        width: 2,
        bgcolor: 'divider',
        borderRadius: 1,
      }}
    >
      <Box
        className="timeline-fill"
        sx={{
          height: '100%',
          transformOrigin: 'top center',
          borderRadius: 1,
          background: ({ palette }) =>
            `linear-gradient(${palette.primary.main}, ${palette.secondary.main})`,
        }}
      />
    </Box>

    {children}
  </Stack>
)
