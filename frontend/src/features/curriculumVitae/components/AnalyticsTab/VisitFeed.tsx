import { useMemo, useState } from 'react'
import { useLinkVisits } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { LinkVisit } from '@features/curriculumVitae/types'
import { isVisibleAction } from '@features/curriculumVitae/utils'
import {
  Button,
  FormControlLabel,
  LinearProgress,
  Stack,
  Switch,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { ReportError } from './ReportError'
import { Visit } from './Visit'

type Filter = 'all' | 'reads' | 'actions'

const FILTERS: Record<Filter, (visit: LinkVisit) => boolean> = {
  all: () => true,
  reads: ({ isRead }) => isRead,
  actions: ({ events }) => events.some(isVisibleAction),
}

type Props = {
  linkId: string
}

export const VisitFeed = ({ linkId }: Props) => {
  const translate = useTranslate('Lk.Analytics.Feed')
  const [filter, setFilter] = useState<Filter>('all')
  const [includeFiltered, setIncludeFiltered] = useState(false)
  const feed = useLinkVisits(linkId, includeFiltered)

  const visible = useMemo(() => feed.visits.filter(FILTERS[filter]), [feed.visits, filter])

  const visitsByVisitor = useMemo(() => {
    const counts = new Map<number, number>()
    feed.visits.forEach(({ visitor }) =>
      counts.set(visitor.ordinal, (counts.get(visitor.ordinal) ?? 0) + 1),
    )
    return counts
  }, [feed.visits])

  if (feed.hasError && feed.visits.length === 0) return <ReportError onRetry={feed.reload} />

  return (
    <Stack spacing={1.5} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack
        direction="row"
        alignItems="center"
        justifyContent="space-between"
        flexWrap="wrap"
        useFlexGap
        spacing={1}
      >
        <Typography variant="h6">{translate('title')}</Typography>

        <FormControlLabel
          sx={{ mr: 0 }}
          control={
            <Switch
              size="small"
              checked={includeFiltered}
              onChange={(_, checked) => setIncludeFiltered(checked)}
            />
          }
          label={<Typography variant="body2">{translate('showFiltered')}</Typography>}
        />
      </Stack>

      <ToggleButtonGroup
        exclusive
        size="small"
        value={filter}
        onChange={(_, value: Filter | null) => value && setFilter(value)}
        sx={{ alignSelf: 'flex-start' }}
      >
        {(Object.keys(FILTERS) as Filter[]).map((key) => (
          <ToggleButton key={key} value={key} sx={{ textTransform: 'none', px: 1.5 }}>
            {translate(`Filters.${key}`)}
          </ToggleButton>
        ))}
      </ToggleButtonGroup>

      {visible.map((visit) => (
        <Visit
          key={visit.id}
          visit={visit}
          visitorVisits={visitsByVisitor.get(visit.visitor.ordinal) ?? 1}
        />
      ))}

      {feed.isLoading && <LinearProgress />}

      {!feed.isLoading && feed.visits.length === 0 && (
        <Typography sx={{ color: 'text.secondary' }}>{translate('empty')}</Typography>
      )}

      {!feed.isLoading && feed.visits.length > 0 && visible.length === 0 && (
        <Typography sx={{ color: 'text.secondary' }}>{translate('noMatches')}</Typography>
      )}

      {feed.hasMore && !feed.isLoading && (
        <Button
          variant="outlined"
          color="inherit"
          onClick={feed.loadMore}
          sx={{ alignSelf: 'center', borderRadius: 999, px: 3, textTransform: 'none' }}
        >
          {translate('loadMore')}
        </Button>
      )}
    </Stack>
  )
}
