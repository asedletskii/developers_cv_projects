import { Ref } from 'react'
import { ExpandMoreRounded } from '@mui/icons-material'
import { Button, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'

type Props = {
  ref?: Ref<HTMLButtonElement>
  isExpanded: boolean
  hiddenCount: number
  onClick(): void
}

export const MoreButton = ({ ref, isExpanded, hiddenCount, onClick }: Props) => {
  const translate = useTranslate('CurriculumVitae.Experience')

  return (
    <Stack alignItems="center" sx={{ mt: 4 }}>
      <Button
        ref={ref}
        variant="outlined"
        onClick={onClick}
        endIcon={
          <ExpandMoreRounded
            sx={{
              transform: isExpanded ? 'rotate(180deg)' : 'none',
              transition: 'transform 0.4s cubic-bezier(0.32, 0.72, 0, 1)',
            }}
          />
        }
        sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
      >
        {isExpanded ? translate('showLess') : translate('showMore', { count: hiddenCount })}
      </Button>
    </Stack>
  )
}
