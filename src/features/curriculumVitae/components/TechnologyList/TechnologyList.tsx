import { Chip, Stack } from '@mui/material'
import { useTechnologyFilter } from '@features/curriculumVitae/components/TechnologyFilterProvider'
import { Technology } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'

type Props = {
  technologies: Technology[]
}

/** Строка компактных чипов со стеком. Клик по чипу включает фильтр по этой технологии. */
export const TechnologyList = ({ technologies }: Props) => {
  const translate = useTranslate('CurriculumVitae.Filter')
  const { isSelected, toggle } = useTechnologyFilter()

  return (
    <Stack direction="row" flexWrap="wrap" gap={0.75}>
      {technologies.map(({ name }) => {
        const selected = isSelected(name)

        return (
          <Chip
            key={name}
            size="small"
            label={name}
            color={selected ? 'primary' : 'default'}
            variant={selected ? 'filled' : 'outlined'}
            onClick={() => toggle(name)}
            aria-pressed={selected}
            aria-label={translate('select', { name })}
          />
        )
      })}
    </Stack>
  )
}
