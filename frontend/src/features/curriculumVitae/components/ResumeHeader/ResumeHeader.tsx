import { EditResumeTitleDialog } from '@features/curriculumVitae/components/EditResumeTitleDialog'
import { glassCard } from '@features/curriculumVitae/styles'
import { CurriculumVitaeError, ResumeView } from '@features/curriculumVitae/types'
import {
  ArrowBackRounded,
  DeleteOutlineRounded,
  EditRounded,
  MoreVertRounded,
  VisibilityOutlined,
} from '@mui/icons-material'
import {
  Button,
  Chip,
  IconButton,
  ListItemIcon,
  Menu,
  MenuItem,
  Stack,
  Typography,
} from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { generatePath, Link as RouterLink } from 'react-router'
import { useResumeHeader } from './hooks'

type Props = {
  userId: string
  view: ResumeView
  onReload: () => void
  onError: (error: CurriculumVitaeError) => void
}

/** Шапка страницы резюме в кабинете: назад, название, язык, предпросмотр и меню действий. */
export const ResumeHeader = ({ userId, view, onReload, onError }: Props) => {
  const translate = useTranslate('Lk.Edit')
  const translateResume = useTranslate('Lk.Resume')
  const header = useResumeHeader({ userId, resumeId: view.id, onRenamed: onReload, onError })

  return (
    <>
      <Stack direction="row" alignItems="center" spacing={1} flexWrap="wrap" useFlexGap>
        <IconButton
          component={RouterLink}
          to={generatePath(ROUTES.lkProfile, { userId })}
          aria-label={translateResume('back')}
          sx={{ ...glassCard, borderRadius: 999 }}
        >
          <ArrowBackRounded />
        </IconButton>

        <Typography variant="h4">{view.title}</Typography>

        <Chip label={view.language.toUpperCase()} size="small" />

        <Stack direction="row" spacing={1} sx={{ ml: 'auto' }}>
          <Button
            component={RouterLink}
            to={generatePath(ROUTES.lkResumePreview, { userId, cvId: view.id })}
            startIcon={<VisibilityOutlined />}
            sx={{ ...glassCard, borderRadius: 999, px: 2, textTransform: 'none' }}
          >
            {translate('preview')}
          </Button>

          <IconButton
            aria-label={translate('actions')}
            onClick={(e) => header.setMenuAnchor(e.currentTarget)}
            sx={{ ...glassCard, borderRadius: 999 }}
          >
            <MoreVertRounded />
          </IconButton>
        </Stack>
      </Stack>

      <Menu anchorEl={header.menuAnchor} open={!!header.menuAnchor} onClose={header.closeMenu}>
        <MenuItem onClick={header.startRename}>
          <ListItemIcon>
            <EditRounded fontSize="small" />
          </ListItemIcon>
          {translate('renameTitle')}
        </MenuItem>

        <MenuItem onClick={header.remove} sx={{ color: 'error.main' }}>
          <ListItemIcon sx={{ color: 'inherit' }}>
            <DeleteOutlineRounded fontSize="small" />
          </ListItemIcon>
          {translate('deleteResume')}
        </MenuItem>
      </Menu>

      {header.isRenaming && (
        <EditResumeTitleDialog userId={userId} resume={view} onClose={header.finishRename} />
      )}
    </>
  )
}
