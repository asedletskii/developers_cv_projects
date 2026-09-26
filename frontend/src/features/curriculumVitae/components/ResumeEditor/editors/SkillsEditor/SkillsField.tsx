import { Person } from '@features/curriculumVitae/types'
import { Autocomplete, TextField } from '@mui/material'
import { FieldPath } from 'react-hook-form'
import { useSkillsField } from './hooks'

type Props = {
  name: Extract<FieldPath<Person>, 'hardSkills' | 'softSkills'>
  label: string
  options: { value: string; label: string }[]
}

/** Выпадающий список с поиском; только значения из справочника — бэк отклоняет остальные. */
export const SkillsField = ({ name, label, options }: Props) => {
  const { selected, onChange, onBlur } = useSkillsField(name, options)

  return (
    <Autocomplete
      multiple
      filterSelectedOptions
      options={options}
      value={selected}
      onChange={(_, next) => onChange(next)}
      onBlur={onBlur}
      isOptionEqualToValue={(option, value) => option.value === value.value}
      renderInput={(params) => <TextField {...params} variant="standard" label={label} />}
    />
  )
}
