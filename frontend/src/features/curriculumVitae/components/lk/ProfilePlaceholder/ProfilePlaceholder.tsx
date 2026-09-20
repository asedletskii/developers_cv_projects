import { Container, Skeleton, Stack } from '@mui/material'

export const ProfilePlaceholder = () => (
  <Container maxWidth="lg" sx={{ py: 6 }}>
    <Stack direction={{ xs: 'column', md: 'row' }} spacing={4}>
      <Stack spacing={2} alignItems="center" sx={{ width: { xs: '100%', md: '25%' } }}>
        <Skeleton variant="circular" width={120} height={120} />

        <Skeleton variant="text" width="80%" height={40} />

        <Skeleton variant="text" width="60%" />

        <Skeleton variant="text" width="60%" />
      </Stack>

      <Stack spacing={2} sx={{ flex: 1 }}>
        <Skeleton variant="text" width={160} height={40} />

        <Skeleton variant="rounded" height={64} />

        <Skeleton variant="rounded" height={64} />

        <Skeleton variant="rounded" height={64} />
      </Stack>
    </Stack>
  </Container>
)
