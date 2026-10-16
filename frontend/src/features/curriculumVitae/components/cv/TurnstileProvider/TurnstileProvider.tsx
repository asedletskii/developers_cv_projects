import { createContext, PropsWithChildren, useMemo } from 'react'
import { useTurnstile } from '@features/curriculumVitae/hooks/useTurnstile'
import { glassCard } from '@features/curriculumVitae/styles'
import { Box, Dialog, DialogContent, DialogTitle } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type TurnstileValue = {
  /** Токен проверки Cloudflare; `undefined`, если проверка выключена (нет site key) */
  getToken: () => Promise<string | undefined>
}

/** Без провайдера проверки нет: `getToken` отдаёт `undefined`, как при выключенном Turnstile */
export const TurnstileContext = createContext<TurnstileValue>({ getToken: async () => undefined })

export const TurnstileProvider = ({ children }: PropsWithChildren) => {
  const translate = useTranslate('CurriculumVitae.Contacts')
  const { setContainer, getToken, isChallengeVisible, cancel } = useTurnstile()
  const value = useMemo(() => ({ getToken }), [getToken])

  return (
    <TurnstileContext.Provider value={value}>
      {children}

      {/* Cloudflare показывает задачу, только если сочтёт её нужной, — тогда открывается модалка.
          keepMounted: виджет рисуется в контейнер заранее, а не в момент открытия. Закрыть можно
          по Esc или клику мимо — это отменяет ожидание токена, отдельная кнопка не нужна */}
      <Dialog
        open={isChallengeVisible}
        onClose={cancel}
        keepMounted
        fullWidth
        maxWidth="xs"
        slotProps={{ paper: { sx: { ...glassCard, backgroundImage: 'none' } } }}
      >
        <DialogTitle>{translate('Challenge.title')}</DialogTitle>

        <DialogContent>
          <Box ref={setContainer} sx={{ display: 'flex', justifyContent: 'center' }} />
        </DialogContent>
      </Dialog>
    </TurnstileContext.Provider>
  )
}
