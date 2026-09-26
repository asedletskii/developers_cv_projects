import { glassCard } from '@features/curriculumVitae/styles'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useLinksTab } from './hooks'
import { LinkActionsMenu } from './LinkActionsMenu'
import { LinkNameDialog } from './LinkNameDialog'
import { LinkRow } from './LinkRow'

type Props = {
  resumeId: string
}

/** «Ссылки»: публичный адрес резюме и именованные ссылки на него — по одной на каждого адресата. */
export const LinksTab = ({ resumeId }: Props) => {
  const translate = useTranslate('Lk.Links')
  const tab = useLinksTab(resumeId)

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
        <LinkRow
          isPublic
          name={translate('publicLink')}
          url={tab.publicUrl}
          isCopied={tab.copiedKey === 'public'}
          onCopy={() => tab.copy('public', tab.publicUrl)}
        />

        {tab.links.map((link) => (
          <LinkRow
            key={link.id}
            name={link.name}
            url={tab.getLinkUrl(link)}
            isCopied={tab.copiedKey === link.id}
            onCopy={() => tab.copy(link.id, tab.getLinkUrl(link))}
            onOpenMenu={tab.openMenu(link)}
          />
        ))}
      </Stack>

      {tab.links.length === 0 && (
        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
          {translate('empty')}
        </Typography>
      )}

      <LinkActionsMenu
        anchor={tab.menu?.anchor ?? null}
        onClose={tab.closeMenu}
        onRename={tab.renameFromMenu}
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
