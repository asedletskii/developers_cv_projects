import { CreateLinkPrompt } from '@features/curriculumVitae/components/lk/CreateLinkPrompt'
import { LinkNameDialog } from '@features/curriculumVitae/components/lk/LinkNameDialog'
import { useResumeLinks } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { ResumeStatistics } from '@features/curriculumVitae/types'
import { AddRounded } from '@mui/icons-material'
import { Button, LinearProgress, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useLinksTab } from './hooks'
import { LinkActionsMenu } from './LinkActionsMenu'
import { LinkRow } from './LinkRow'
import { LinksLoadError } from './LinksLoadError'

type Props = {
  linksState: ReturnType<typeof useResumeLinks>
  statistics?: ResumeStatistics
}

export const LinksTab = ({ linksState, statistics }: Props) => {
  const translate = useTranslate('Lk.Links')
  const tab = useLinksTab(linksState)

  if (tab.isLoading) return <LinearProgress />

  if (tab.hasLoadError) return <LinksLoadError onRetry={linksState.reload} />

  if (tab.links.length === 0) return <CreateLinkPrompt onCreate={linksState.create} />

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <Typography variant="h6">{translate('title')}</Typography>

        <Button startIcon={<AddRounded />} onClick={tab.openCreate} sx={{ textTransform: 'none' }}>
          {translate('add')}
        </Button>
      </Stack>

      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {translate('hint')}
      </Typography>

      <Stack spacing={1}>
        {tab.links.map((link) => (
          <LinkRow
            key={link.id}
            name={link.name}
            url={tab.getLinkUrl(link)}
            isActive={link.status === 'active'}
            isCopied={tab.copiedKey === link.id}
            stats={statistics?.links.find(({ id }) => id === link.id)}
            onCopy={() => tab.copy(link.id, tab.getLinkUrl(link))}
            onOpenMenu={tab.openMenu(link)}
          />
        ))}
      </Stack>

      <LinkActionsMenu
        anchor={tab.menu?.anchor ?? null}
        isActive={tab.menu?.link.status !== 'inactive'}
        onClose={tab.closeMenu}
        onRename={tab.renameFromMenu}
        onToggleStatus={tab.toggleStatusFromMenu}
        onRemove={tab.removeFromMenu}
      />

      {tab.dialog && (
        <LinkNameDialog
          title={translate(tab.dialog.mode === 'create' ? 'add' : 'rename')}
          initialName={tab.dialog.mode === 'rename' ? tab.dialog.link.name : ''}
          onSave={tab.saveName}
          onClose={tab.closeDialog}
        />
      )}
    </Stack>
  )
}
