import { Fragment } from 'react'
import { SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { Resume, SectionSettings } from '@features/curriculumVitae/types'
import { Container } from '@mui/material'

type Props = {
  resume: Resume
  sections: SectionSettings[]
}

export const Sections = ({ resume, sections }: Props) => (
  <>
    {sections.map(({ id }) => {
      const definition = SECTION_DEFINITIONS[id]

      return definition.isFullWidth ? (
        <Fragment key={id}>{definition.render(resume)}</Fragment>
      ) : (
        <Container key={id} maxWidth="lg" sx={{ position: 'relative' }}>
          {definition.render(resume)}
        </Container>
      )
    })}
  </>
)
