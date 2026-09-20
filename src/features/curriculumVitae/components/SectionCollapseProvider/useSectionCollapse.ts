import { useContext } from 'react'
import { SectionCollapseContext } from './SectionCollapseProvider'

export const useSectionCollapse = () => useContext(SectionCollapseContext)
