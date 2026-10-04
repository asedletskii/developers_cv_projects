import { useLinkVisits } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { Button, LinearProgress, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { ReportError } from './ReportError'
import { Visit } from './Visit'

type Props = {
  linkId: string
}

export const VisitFeed = ({ linkId }: Props) => {
  const translate = useTranslate('Lk.Analytics.Feed')
  const feed = useLinkVisits(linkId)

  if (feed.hasError && feed.visits.length === 0) return <ReportError onRetry={feed.reload} />

  return (
    <Stack spacing={1.5} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Typography variant="h6">{translate('title')}</Typography>

      {feed.visits.map((visit) => (
        <Visit key={visit.id} visit={visit} />
      ))}

      {feed.isLoading && <LinearProgress />}

      {!feed.isLoading && feed.visits.length === 0 && (
        <Typography sx={{ color: 'text.secondary' }}>{translate('empty')}</Typography>
      )}

      {feed.hasMore && !feed.isLoading && (
        <Button onClick={feed.loadMore} sx={{ alignSelf: 'center', textTransform: 'none' }}>
          {translate('loadMore')}
        </Button>
      )}
    </Stack>
  )
}
