import { useEffect, useRef } from 'react'
import { CloseRounded } from '@mui/icons-material'
import { IconButton, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useAppSelector } from '@shared/store'
import { selectToasts } from '@shared/store/toast'
import { closeToast } from '@shared/utils'
import { useSnackbar } from 'notistack'

/**
 * Показывает тосты из стора через notistack. Источник правды — стор: `showToast` добавляет тост
 * и сам убирает его по таймеру, а здесь мы лишь рисуем то, что в сторе есть, и закрываем то, что из него пропало.
 */
export const Host = () => {
  const translate = useTranslate('Toasts')
  const toasts = useAppSelector(selectToasts)
  const { enqueueSnackbar, closeSnackbar } = useSnackbar()
  const shownIdsRef = useRef(new Set<number>())

  useEffect(() => {
    const currentIds = new Set(toasts.map(({ id }) => id))

    toasts.forEach(({ id, type, messageKey, titleKey }) => {
      if (shownIdsRef.current.has(id)) return
      shownIdsRef.current.add(id)

      enqueueSnackbar(
        titleKey ? (
          <Stack>
            <Typography sx={{ fontWeight: 600, color: 'inherit' }}>{translate(titleKey)}</Typography>
            <Typography sx={{ color: 'inherit' }}>{translate(messageKey)}</Typography>
          </Stack>
        ) : (
          translate(messageKey)
        ),
        {
          key: id,
          variant: type,
          // Время жизни считает showToast (удаляет тост из стора), поэтому сам notistack не скрывает
          persist: true,
          action: (key) => (
            <IconButton size="small" color="inherit" onClick={() => closeToast(Number(key))}>
              <CloseRounded fontSize="small" />
            </IconButton>
          ),
        },
      )
    })

    shownIdsRef.current.forEach((id) => {
      if (!currentIds.has(id)) {
        closeSnackbar(id)
        shownIdsRef.current.delete(id)
      }
    })
  }, [toasts, translate, enqueueSnackbar, closeSnackbar])

  return null
}
