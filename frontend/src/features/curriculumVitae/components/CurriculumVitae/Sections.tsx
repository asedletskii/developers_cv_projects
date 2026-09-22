import { Fragment } from 'react'
import { SECTION_DEFINITIONS } from '@features/curriculumVitae/sections'
import { Person, SectionSettings } from '@features/curriculumVitae/types'
import { Container } from '@mui/material'

type Props = {
  person: Person
  sections: SectionSettings[]
}

/** Секции в порядке из настроек: обычные — в контейнере, полноширинные (лента стека) — без него. */
export const Sections = ({ person, sections }: Props) => (
  <>
    {sections.map(({ id }) => {
      const definition = SECTION_DEFINITIONS[id]

      return definition.isFullWidth ? (
        <Fragment key={id}>{definition.render(person)}</Fragment>
      ) : (
        <Container key={id} maxWidth="lg" sx={{ position: 'relative' }}>
          {definition.render(person)}
        </Container>
      )
    })}
  </>
)
