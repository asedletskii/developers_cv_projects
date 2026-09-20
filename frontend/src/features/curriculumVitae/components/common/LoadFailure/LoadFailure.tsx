import { useState } from 'react'
import { AuthShell } from '@features/auth/components/AuthShell'
import { PULSE_RING, RISE, SHAKE, STILL } from '@features/auth/components/AuthShell/constants'
import { glassCard, secondaryButton } from '@features/curriculumVitae/styles'
import { CurriculumVitaeError } from '@features/curriculumVitae/types'
import {
  CheckRounded,
  CloudOffRounded,
  ContentCopyRounded,
  HomeRounded,
  RefreshRounded,
  SearchOffRounded,
} from '@mui/icons-material'
import { Box, Button, Chip, IconButton, Stack, Tooltip, Typography } from '@mui/material'
import { ROUTES } from '@shared/constants'
import { useTranslate } from '@shared/hooks'
import { Link as RouterLink } from 'react-router'
import { FLICKER, SWEEP } from './constants'

type Props = {
  error: CurriculumVitaeError
  subject?: 'resume' | 'profile'
}

export const LoadFailure = ({ error, subject = 'resume' }: Props) => {
  const translate = useTranslate('CurriculumVitae.Error')
  const translateProfile = useTranslate('Lk.Error')
  const translateMessage = subject === 'profile' ? translateProfile : translate
  const [isCopied, setIsCopied] = useState(false)

  const isNotFound = error.kind === 'notFound'
  const Icon = isNotFound ? SearchOffRounded : CloudOffRounded

  const diagnostics = [
    error.status && translate('errorCode', { status: error.status }),
    error.requestId && translate('requestId', { requestId: error.requestId }),
  ]
    .filter(Boolean)
    .join(' · ')

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(diagnostics)
      setIsCopied(true)
      setTimeout(() => setIsCopied(false), 1800)
    } catch {}
  }

  return (
    <AuthShell>
      <Stack
        role="alert"
        spacing={2.5}
        alignItems="center"
        textAlign="center"
        sx={{
          maxWidth: 480,
          mx: 'auto',
          p: { xs: 3, md: 5 },
          ...glassCard,
          animation: `${RISE} 0.5s ease-out both`,
          ...STILL,
        }}
      >
        <Box sx={{ position: 'relative', width: 88, height: 88 }}>
          <Box
            aria-hidden
            sx={{
              position: 'absolute',
              inset: 0,
              borderRadius: '50%',
              border: 2,
              borderColor: isNotFound ? 'primary.main' : 'error.main',
              animation: `${PULSE_RING} 2.6s ease-out infinite`,
              ...STILL,
            }}
          />
          <Stack
            alignItems="center"
            justifyContent="center"
            sx={{
              position: 'absolute',
              inset: 0,
              borderRadius: '50%',
              bgcolor: 'action.hover',
              color: isNotFound ? 'primary.main' : 'error.main',
              animation: isNotFound ? undefined : `${SHAKE} 0.6s ease-in-out 0.3s both`,
              ...STILL,
            }}
          >
            <Icon
              aria-hidden
              sx={{
                fontSize: 44,
                animation: isNotFound
                  ? `${SWEEP} 3s ease-in-out infinite`
                  : `${FLICKER} 3.5s linear infinite`,
                ...STILL,
              }}
            />
          </Stack>
        </Box>

        <Typography variant="h4" component="h1">
          {translateMessage(isNotFound ? 'notFoundTitle' : 'failedTitle')}
        </Typography>

        <Typography sx={{ color: 'text.secondary' }}>
          {translateMessage(isNotFound ? 'notFoundText' : 'failedText')}
        </Typography>

        {diagnostics && (
          <Chip
            variant="outlined"
            label={diagnostics}
            onDelete={handleCopy}
            deleteIcon={
              <Tooltip title={translate(isCopied ? 'copied' : 'copyDetails')}>
                <IconButton
                  size="small"
                  aria-label={translate('copyDetails')}
                  sx={{ color: isCopied ? 'success.main' : 'inherit' }}
                >
                  {isCopied ? (
                    <CheckRounded fontSize="small" />
                  ) : (
                    <ContentCopyRounded fontSize="small" />
                  )}
                </IconButton>
              </Tooltip>
            }
            sx={{
              maxWidth: '100%',
              height: 'auto',
              py: 1,
              pl: 0.75,
              pr: 0.5,
              borderRadius: 3,
              fontFamily: 'monospace',
              '& .MuiChip-label': { whiteSpace: 'normal', px: 2, py: 0.5 },
            }}
          />
        )}

        <Stack spacing={1.5} sx={{ width: '100%' }}>
          {!isNotFound && (
            <Button
              variant="contained"
              size="large"
              startIcon={<RefreshRounded />}
              onClick={() => window.location.reload()}
              sx={{ borderRadius: 999, textTransform: 'none' }}
            >
              {translate('retry')}
            </Button>
          )}

          <Button
            component={RouterLink}
            to={ROUTES.home}
            variant={isNotFound ? 'contained' : 'outlined'}
            size="large"
            startIcon={<HomeRounded />}
            sx={{
              borderRadius: 999,
              textTransform: 'none',
              ...(isNotFound ? {} : secondaryButton),
            }}
          >
            {translate('home')}
          </Button>
        </Stack>
      </Stack>
    </AuthShell>
  )
}
