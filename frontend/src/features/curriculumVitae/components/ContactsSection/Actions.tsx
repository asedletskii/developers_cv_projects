import { useTurnstile } from '@features/curriculumVitae/hooks'
import { glassCard } from '@features/curriculumVitae/styles'
import { ContactLink } from '@features/curriculumVitae/utils'
import { Box, Dialog, DialogContent, DialogTitle, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { ContactButton } from './ContactButton'

type Props = {
  resumeId: string
  links: ContactLink[]
}

export const Actions = ({ resumeId, links }: Props) => {
  const translate = useTranslate('CurriculumVitae.Contacts')
  const { setContainer, getToken, isChallengeVisible, cancel } = useTurnstile()

  return (
    <>
      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        spacing={1.5}
        width={{ xs: '100%', sm: 'auto' }}
        sx={{ position: 'relative' }}
      >
        {links.map((link) => (
          <ContactButton key={link.key} resumeId={resumeId} link={link} getToken={getToken} />
        ))}
      </Stack>

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
    </>
  )
}
