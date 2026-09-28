import {
  DeleteOutlineRounded,
  EditRounded,
  PauseRounded,
  PlayArrowRounded,
} from '@mui/icons-material'
import { Divider, ListItemIcon, Menu, MenuItem } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  anchor: HTMLElement | null
  isActive: boolean
  onClose: () => void
  onRename: () => void
  onToggleStatus: () => void
  onRemove: () => void
}

export const LinkActionsMenu = ({
  anchor,
  isActive,
  onClose,
  onRename,
  onToggleStatus,
  onRemove,
}: Props) => {
  const translate = useTranslate('Lk.Links')

  return (
    <Menu anchorEl={anchor} open={!!anchor} onClose={onClose}>
      <MenuItem onClick={onRename}>
        <ListItemIcon>
          <EditRounded fontSize="small" />
        </ListItemIcon>
        {translate('rename')}
      </MenuItem>

      <MenuItem onClick={onToggleStatus}>
        <ListItemIcon>
          {isActive ? <PauseRounded fontSize="small" /> : <PlayArrowRounded fontSize="small" />}
        </ListItemIcon>
        {translate(isActive ? 'deactivate' : 'activate')}
      </MenuItem>

      <Divider />

      <MenuItem onClick={onRemove} sx={{ color: 'error.main' }}>
        <ListItemIcon sx={{ color: 'inherit' }}>
          <DeleteOutlineRounded fontSize="small" />
        </ListItemIcon>
        {translate('delete')}
      </MenuItem>
    </Menu>
  )
}
