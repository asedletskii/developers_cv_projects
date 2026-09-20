import { ResumeStatus } from '@features/curriculumVitae/types'
import {
  DeleteOutlineRounded,
  EditRounded,
  PublicOffRounded,
  PublicRounded,
} from '@mui/icons-material'
import { Divider, ListItemIcon, Menu, MenuItem } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  anchor: HTMLElement | null
  status: ResumeStatus
  onClose: () => void
  onRename: () => void
  onPublish: () => void
  onUnpublish: () => void
  onRemove: () => void
}

export const ResumeActionsMenu = ({
  anchor,
  status,
  onClose,
  onRename,
  onPublish,
  onUnpublish,
  onRemove,
}: Props) => {
  const translate = useTranslate('Lk.Resumes')

  return (
    <Menu anchorEl={anchor} open={!!anchor} onClose={onClose}>
      <MenuItem onClick={onRename}>
        <ListItemIcon>
          <EditRounded fontSize="small" />
        </ListItemIcon>
        {translate('editTitle')}
      </MenuItem>

      {status === 'opened' ? (
        <MenuItem onClick={onUnpublish}>
          <ListItemIcon>
            <PublicOffRounded fontSize="small" />
          </ListItemIcon>
          {translate('unpublish')}
        </MenuItem>
      ) : (
        <MenuItem onClick={onPublish}>
          <ListItemIcon>
            <PublicRounded fontSize="small" />
          </ListItemIcon>
          {translate('publish')}
        </MenuItem>
      )}

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
