import { Person } from '@features/curriculumVitae/types'
import { MenuItem, TextField, TextFieldProps } from '@mui/material'
import { Controller, FieldPath, useFormContext } from 'react-hook-form'

type Option = { value: string | number; label: string }

type Props = {
  name: FieldPath<Person>
  label: string
  multiple?: boolean
  options: Option[]
} & Pick<TextFieldProps, 'required' | 'sx'>

/**
 * Единый select для enum-полей формы: MUI Select требует контролируемое значение, поэтому
 * оборачиваем в RHF Controller (просто register тут не отражает изменения при reset формы).
 * Значение типизировано широко (any/never) — конкретные типы у каждого enum-поля разные,
 * а этот компонент общий для всех них.
 */
export const SelectField = ({ name, label, multiple, options, ...rest }: Props) => {
  const { control } = useFormContext<Person>()

  return (
    <Controller
      name={name}
      control={control}
      render={({ field }) => (
        <TextField
          select
          label={label}
          value={field.value ?? (multiple ? [] : '')}
          onChange={field.onChange}
          onBlur={field.onBlur}
          inputRef={field.ref}
          slotProps={{ select: { multiple } }}
          {...rest}
        >
          {!multiple && <MenuItem value="">—</MenuItem>}
          {options.map((opt) => (
            <MenuItem key={opt.value} value={opt.value}>
              {opt.label}
            </MenuItem>
          ))}
        </TextField>
      )}
    />
  )
}
