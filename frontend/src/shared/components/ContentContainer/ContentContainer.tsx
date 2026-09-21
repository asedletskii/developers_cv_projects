import { ComponentType, JSX, PropsWithChildren } from 'react'
import { Container, SkeletonProps, Stack, SxProps } from '@mui/material'
import { BackButton } from './BackButton'
import { Title } from './Title'

type Props = PropsWithChildren<{
  title?: string
  Skeleton?: ComponentType<SkeletonProps>
  isLoading?: boolean
  backButton?: boolean
  containerSx?: SxProps
  InfoBlock?: JSX.Element
}>

export const ContentContainer = ({
  title,
  containerSx,
  Skeleton,
  children,
  isLoading = false,
  backButton = false,
}: Props) => (
  <Container
    maxWidth={'lg'}
    sx={{
      display: 'flex',
      position: 'relative',
      flexDirection: 'column',
      flex: 1,
      overflow: 'visible',
      ...containerSx,
    }}
  >
    {isLoading ? (
      <Stack width="100%" alignItems={'center'}>
        {Skeleton && <Skeleton />}
      </Stack>
    ) : (
      <>
        {backButton && <BackButton />}

        {title && <Title title={title} hasBackButton={backButton} />}

        {children}
      </>
    )}
  </Container>
)
