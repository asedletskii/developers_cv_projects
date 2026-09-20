import { Resume } from '@features/curriculumVitae/types'
import { Autocomplete, TextField } from '@mui/material'
import { FieldPath } from 'react-hook-form'
import { useSkillsField } from './hooks'

type Props = {
  name: Extract<FieldPath<Resume>, 'skills.hardSkills' | 'skills.softSkills'>
  label: string
  helperText: string
  options: string[]
  getOptionLabel: (skill: string) => string
}

export const SkillsField = ({ name, label, helperText, options, getOptionLabel }: Props) => {
  const { selected, onChange, onBlur } = useSkillsField(name)

  return (
    <Autocomplete
      multiple
      freeSolo
      autoSelect
      filterSelectedOptions
      options={options}
      value={selected}
      getOptionLabel={getOptionLabel}
      onChange={(_, next) => onChange(next)}
      onBlur={onBlur}
      renderInput={(params) => (
        <TextField {...params} variant="standard" label={label} helperText={helperText} />
      )}
    />
  )
}
