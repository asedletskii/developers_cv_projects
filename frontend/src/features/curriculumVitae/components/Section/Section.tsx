import { PropsWithChildren } from 'react'
import { SectionId } from '@features/curriculumVitae/types'
import { Box, Stack } from '@mui/material'
import { useSectionCollapse } from '../SectionCollapseProvider'
import { contentSx } from './constants'
import { Header } from './Header'
import { Panel } from './Panel'

type Props = PropsWithChildren<{
  id: SectionId
  title: string

  /** Можно ли сворачивать секцию (по умолчанию да) */
  collapsible?: boolean
}>

/** Секция страницы. Сворачиваемая работает как аккордеон: заголовок скрывает и показывает содержимое. */
export const Section = ({ id, title, collapsible = true, children }: Props) => {
  const { isCollapsed, toggle } = useSectionCollapse()
  const collapsed = collapsible && isCollapsed(id)

  return (
    <Stack id={id} component="section" className="reveal" sx={{ py: 2, scrollMarginTop: 72 }}>
      <Header
        id={id}
        title={title}
        isCollapsible={collapsible}
        isCollapsed={collapsed}
        onToggle={() => toggle(id)}
      />

      {collapsible ? (
        <Panel id={id} isCollapsed={collapsed}>
          {children}
        </Panel>
      ) : (
        <Box sx={contentSx}>{children}</Box>
      )}
    </Stack>
  )
}
