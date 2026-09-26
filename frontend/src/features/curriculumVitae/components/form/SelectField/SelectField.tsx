import { Resume } from '@features/curriculumVitae/types'
import { MenuItem, TextField, TextFieldProps } from '@mui/material'
import { Controller, FieldPath, useFormContext } from 'react-hook-form'

type Option = { value: string | number; label: string }

type Props = {
  name: FieldPath<Resume>
  label: string
  multiple?: boolean

  /** Есть ли пункт «—» (сбросить значение). Для обязательных полей его нет */
  clearable?: boolean
  options: Option[]
} & Pick<TextFieldProps, 'required' | 'sx'>

/**
 * Единый select для enum-полей формы: MUI Select требует контролируемое значение, поэтому
 * оборачиваем в RHF Controller (просто register тут не отражает изменения при reset формы).
 * Значение типизировано широко (any/never) — конкретные типы у каждого enum-поля разные,
 * а этот компонент общий для всех них.
 */
export const SelectField = ({
  name,
  label,
  multiple,
  clearable = true,
  options,
  ...rest
}: Props) => {
  const { control } = useFormContext<Resume>()

  return (
    <Controller
      name={name}
      control={control}
      render={({ field }) => (
        <TextField
          select
          fullWidth
          label={label}
          value={field.value ?? (multiple ? [] : '')}
          onChange={field.onChange}
          onBlur={field.onBlur}
          inputRef={field.ref}
          slotProps={{ select: { multiple } }}
          {...rest}
        >
          {!multiple && clearable && <MenuItem value="">—</MenuItem>}
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
