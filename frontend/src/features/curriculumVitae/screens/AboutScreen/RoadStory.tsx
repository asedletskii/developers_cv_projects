import { Box, Container, Stack, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { ROAD_NODES, ROAD_PATH_D, ROAD_VIEWBOX, STEPS } from './constants'
import { useRoadAnimation } from './hooks/useRoadAnimation'

export const RoadStory = () => {
  const translate = useTranslate('AboutPage.Steps')
  const stageRef = useRoadAnimation(STEPS.length)

  return (
    <Box ref={stageRef} sx={{ position: 'relative', height: `${STEPS.length * 80}vh` }}>
      <Box
        sx={{
          position: 'sticky',
          top: 0,
          height: '100vh',
          overflow: 'hidden',
          display: 'flex',
          alignItems: 'center',
        }}
      >
        <Container maxWidth="lg">
          <Stack direction="row" spacing={{ md: 6, lg: 10 }} alignItems="center">
            <Box sx={{ position: 'relative', width: 240, flexShrink: 0 }}>
              <Box
                component="svg"
                viewBox={ROAD_VIEWBOX}
                sx={{ width: '100%', height: 'auto', overflow: 'visible' }}
              >
                <Box
                  component="path"
                  d={ROAD_PATH_D}
                  fill="none"
                  strokeWidth={3}
                  strokeLinecap="round"
                  stroke="currentColor"
                  sx={{ color: 'divider' }}
                />
                <Box
                  component="path"
                  className="road-path"
                  d={ROAD_PATH_D}
                  fill="none"
                  strokeWidth={4}
                  strokeLinecap="round"
                  stroke="currentColor"
                  sx={{ color: 'primary.main', filter: 'drop-shadow(0 0 6px currentColor)' }}
                />

                {STEPS.map((step, i) => (
                  <Box
                    key={step.key}
                    component="circle"
                    className="road-node"
                    cx={ROAD_NODES[i].x}
                    cy={ROAD_NODES[i].y}
                    r={14}
                    fill="currentColor"
                    sx={{ color: 'primary.main' }}
                  />
                ))}
              </Box>
            </Box>

            <Box sx={{ position: 'relative', flex: 1, minHeight: 280 }}>
              {STEPS.map(({ key, Icon }, index) => (
                <Stack key={key} className="road-frame" spacing={2}>
                  <Stack direction="row" spacing={2} alignItems="center">
                    <Stack
                      alignItems="center"
                      justifyContent="center"
                      sx={{
                        width: 72,
                        height: 72,
                        borderRadius: '50%',
                        bgcolor: 'action.hover',
                        color: 'primary.main',
                      }}
                    >
                      <Icon sx={{ fontSize: 36 }} />
                    </Stack>
                    <Typography sx={{ color: 'gray.light', fontWeight: 600 }}>
                      {translate('step', { number: index + 1 })}
                    </Typography>
                  </Stack>

                  <Typography variant="h2" sx={{ fontSize: 'clamp(1.75rem, 3.5vw, 3rem)' }}>
                    {translate(`${key}.title`)}
                  </Typography>
                  <Typography sx={{ color: 'text.secondary', fontSize: '1.1rem', maxWidth: 480 }}>
                    {translate(`${key}.text`)}
                  </Typography>
                </Stack>
              ))}
            </Box>
          </Stack>
        </Container>
      </Box>
    </Box>
  )
}
