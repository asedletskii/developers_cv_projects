import { useContext } from 'react'
import { SectionCollapseContext } from '@features/curriculumVitae/components/cv/SectionCollapseProvider'

export const useSectionCollapse = () => useContext(SectionCollapseContext)
