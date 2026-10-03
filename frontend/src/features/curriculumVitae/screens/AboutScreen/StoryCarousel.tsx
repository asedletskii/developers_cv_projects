import { MouseEvent } from 'react'
import { glassCard } from '@features/curriculumVitae/styles'
import { Box, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { STEP_DURATION_MS, STEP_FILL, STEPS } from './constants'
import { useStepCarousel } from './hooks/useStepCarousel'
import { useStepIconPop } from './hooks/useStepIconPop'

export const StoryCarousel = () => {
  const translate = useTranslate('AboutPage.Steps')
  const { index, next, prev, select, pause, resume } = useStepCarousel(STEPS.length)
  const iconRef = useStepIconPop(index)
  const { key, Icon } = STEPS[index]

  const handleFrameClick = (e: MouseEvent<HTMLDivElement>) => {
    const { left, width } = e.currentTarget.getBoundingClientRect()
    if (e.clientX - left < width / 2) prev()
    else next()
  }

  return (
    <Stack
      className="reveal"
      spacing={2.5}
      onMouseEnter={pause}
      onMouseLeave={resume}
      sx={{ ...glassCard, p: { xs: 2.5, sm: 4 }, maxWidth: 560, mx: 'auto' }}
    >
      <Stack direction="row" spacing={0.75}>
        {STEPS.map((step, i) => (
          <Box
            key={step.key}
            sx={{
              flex: 1,
              height: 4,
              borderRadius: 999,
              bgcolor: 'action.hover',
              overflow: 'hidden',
            }}
          >
            <Box
              key={i <= index ? `${index}-fill` : 'idle'}
              sx={{
                height: '100%',
                borderRadius: 999,
                bgcolor: 'primary.main',
                width: i < index ? '100%' : 0,
                animation:
                  i === index ? `${STEP_FILL} ${STEP_DURATION_MS}ms linear forwards` : 'none',
              }}
            />
          </Box>
        ))}
      </Stack>

      <Stack
        direction="row"
        spacing={2.5}
        alignItems="center"
        onClick={handleFrameClick}
        sx={{ minHeight: 160, cursor: 'pointer', userSelect: 'none' }}
      >
        <Stack
          ref={iconRef}
          alignItems="center"
          justifyContent="center"
          sx={{
            flexShrink: 0,
            width: 64,
            height: 64,
            borderRadius: '50%',
            bgcolor: 'action.hover',
            color: 'primary.main',
          }}
        >
          <Icon fontSize="large" />
        </Stack>

        <Stack spacing={0.5}>
          <Typography sx={{ color: 'gray.light', fontWeight: 600 }}>
            {translate('step', { number: index + 1 })}
          </Typography>
          <Typography variant="h6">{translate(`${key}.title`)}</Typography>
          <Typography sx={{ color: 'text.secondary' }}>{translate(`${key}.text`)}</Typography>
        </Stack>
      </Stack>

      <Stack direction="row" justifyContent="center" spacing={1}>
        {STEPS.map((step, i) => (
          <Box
            key={step.key}
            component="button"
            type="button"
            aria-label={translate('step', { number: i + 1 })}
            onClick={(e) => {
              e.stopPropagation()
              select(i)
            }}
            sx={{
              width: 8,
              height: 8,
              p: 0,
              border: 0,
              borderRadius: '50%',
              cursor: 'pointer',
              bgcolor: i === index ? 'primary.main' : 'action.hover',
            }}
          />
        ))}
      </Stack>
    </Stack>
  )
}
