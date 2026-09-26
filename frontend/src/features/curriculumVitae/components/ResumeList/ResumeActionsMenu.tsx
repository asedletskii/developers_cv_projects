import { DeleteOutlineRounded, EditRounded } from '@mui/icons-material'
import { ListItemIcon, Menu, MenuItem } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  anchor: HTMLElement | null
  onClose: () => void
  onRename: () => void
  onRemove: () => void
}

/** Меню действий над резюме в списке: переименовать и удалить. */
export const ResumeActionsMenu = ({ anchor, onClose, onRename, onRemove }: Props) => {
  const translate = useTranslate('Lk.Resumes')

  return (
    <Menu anchorEl={anchor} open={!!anchor} onClose={onClose}>
      <MenuItem onClick={onRename}>
        <ListItemIcon>
          <EditRounded fontSize="small" />
        </ListItemIcon>
        {translate('editTitle')}
      </MenuItem>

      <MenuItem onClick={onRemove} sx={{ color: 'error.main' }}>
        <ListItemIcon sx={{ color: 'inherit' }}>
          <DeleteOutlineRounded fontSize="small" />
        </ListItemIcon>
        {translate('delete')}
      </MenuItem>
    </Menu>
  )
}
