import { TECHNOLOGY_SUGGESTIONS } from '@features/curriculumVitae/constants'
import { Autocomplete, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useTechnologyListField } from './hooks'

type Props = {
  name: 'technologies' | `experience.${number}.stack` | `projects.${number}.stack`
}

/** Технологии выпадающим списком с подсказками; если нужной нет — её можно вписать и нажать Enter.
 *  Переиспользуется в стэке резюме, опыте работы и проектах: везде одна и та же форма `Technology[]`. */
export const TechnologyListField = ({ name }: Props) => {
  const translate = useTranslate('Lk.Edit.TechStack')
  const { names, onChange, onBlur } = useTechnologyListField(name)

  return (
    <Autocomplete
      multiple
      freeSolo
      autoSelect
      filterSelectedOptions
      options={TECHNOLOGY_SUGGESTIONS}
      value={names}
      onChange={(_, next) => onChange(next)}
      onBlur={onBlur}
      renderInput={(params) => (
        <TextField {...params} variant="standard" placeholder={translate('addTechnology')} />
      )}
    />
  )
}
