import { useContext } from 'react'
import { SectionCollapseContext } from '@features/curriculumVitae/components/SectionCollapseProvider'

export const useSectionCollapse = () => useContext(SectionCollapseContext)
