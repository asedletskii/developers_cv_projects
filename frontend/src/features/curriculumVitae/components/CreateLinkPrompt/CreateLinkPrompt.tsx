import { useState } from 'react'
import { LinkNameDialog } from '@features/curriculumVitae/components/LinkNameDialog'
import { glassCard } from '@features/curriculumVitae/styles'
import { LinkRounded } from '@mui/icons-material'
import { Button, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  onCreate: (name: string) => void
}

/**
 * Призыв создать первую ссылку на резюме — показывается, пока ссылок ни одной нет (резюме
 * создаётся без них: делиться можно только по отдельной ссылке, свой id резюме напрямую не
 * светим). Имя по умолчанию — «Публичная ссылка», раз она пока единственная и для всех.
 */
export const CreateLinkPrompt = ({ onCreate }: Props) => {
  const translate = useTranslate('Lk.Links')
  const [isOpen, setIsOpen] = useState(false)

  return (
    <Stack spacing={1.5} alignItems="flex-start" sx={{ p: { xs: 2, sm: 3 }, ...glassCard }}>
      <Stack direction="row" spacing={1} alignItems="center">
        <LinkRounded color="primary" fontSize="small" />
        <Typography sx={{ fontWeight: 600 }}>{translate('noLinksTitle')}</Typography>
      </Stack>

      <Typography variant="body2" sx={{ color: 'text.secondary' }}>
        {translate('noLinksText')}
      </Typography>

      <Button variant="contained" onClick={() => setIsOpen(true)} sx={{ textTransform: 'none' }}>
        {translate('createFirst')}
      </Button>

      {isOpen && (
        <LinkNameDialog
          title={translate('add')}
          initialName={translate('publicLinkDefaultName')}
          onSave={onCreate}
          onClose={() => setIsOpen(false)}
        />
      )}
    </Stack>
  )
}
