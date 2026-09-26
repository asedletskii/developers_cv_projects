import { ReactNode } from 'react'
import { glassCard } from '@features/curriculumVitae/styles'
import { KeyboardArrowDownRounded, KeyboardArrowUpRounded } from '@mui/icons-material'
import { Checkbox, FormControlLabel, IconButton, Stack, Switch, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  title: string
  enabled: boolean
  onToggleEnabled: () => void
  collapsed: boolean
  onToggleCollapsed: () => void

  /** Можно ли вообще сворачивать секцию на публичной странице (см. SECTION_DEFINITIONS.isCollapsible) */
  collapsible?: boolean
  /** Закреплённая секция («Обо мне»): без переключателей и стрелок, всегда первая и включена */
  pinned?: boolean
  canMoveUp: boolean
  canMoveDown: boolean
  onMoveUp: () => void
  onMoveDown: () => void
  children?: ReactNode
}

/** Карточка одной настраиваемой секции резюме: заголовок, вкл/выкл, свёрнута ли по умолчанию,
 *  стрелки порядка, форма внутри. */
export const SectionCard = ({
  title,
  enabled,
  onToggleEnabled,
  collapsed,
  onToggleCollapsed,
  collapsible = true,
  pinned = false,
  canMoveUp,
  canMoveDown,
  onMoveUp,
  onMoveDown,
  children,
}: Props) => {
  const translate = useTranslate('Lk.Edit')

  return (
    <Stack spacing={2} sx={{ p: { xs: 2, sm: 3 }, ...glassCard, opacity: enabled ? 1 : 0.6 }}>
      <Stack direction="row" alignItems="center" spacing={1} flexWrap="wrap">
        {!pinned && (
          <Stack direction="row">
            <IconButton
              size="small"
              aria-label={translate('moveUp')}
              disabled={!canMoveUp}
              onClick={onMoveUp}
            >
              <KeyboardArrowUpRounded />
            </IconButton>

            <IconButton
              size="small"
              aria-label={translate('moveDown')}
              disabled={!canMoveDown}
              onClick={onMoveDown}
            >
              <KeyboardArrowDownRounded />
            </IconButton>
          </Stack>
        )}

        <Typography variant="h6" sx={{ flex: 1 }}>
          {title}
        </Typography>

        {collapsible && !pinned && (
          <FormControlLabel
            control={
              <Checkbox
                size="small"
                checked={collapsed}
                onChange={onToggleCollapsed}
                disabled={!enabled}
              />
            }
            label={translate('sectionCollapsed')}
            sx={{ mr: 0 }}
          />
        )}

        {!pinned && (
          <Switch
            checked={enabled}
            onChange={onToggleEnabled}
            inputProps={{ 'aria-label': translate('sectionEnabled') }}
          />
        )}
      </Stack>

      {enabled && children}
    </Stack>
  )
}
