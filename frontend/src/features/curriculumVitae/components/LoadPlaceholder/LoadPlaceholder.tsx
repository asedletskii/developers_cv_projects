import { Container, Skeleton, Stack } from '@mui/material'

/** Заглушка на время первой загрузки резюме. */
export const LoadPlaceholder = () => (
  <Container maxWidth="lg">
    <Stack spacing={3} alignItems="center" width="100%" py={12}>
      <Skeleton variant="text" width={220} height={56} />
      <Skeleton variant="text" width={340} height={56} />
      <Skeleton variant="text" width={280} height={32} />
      <Skeleton variant="rounded" width="100%" height={220} sx={{ mt: 6 }} />
    </Stack>
  </Container>
)
