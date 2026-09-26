import { Person } from '@features/curriculumVitae/types'
import { DatePicker, DatePickerProps } from '@mui/x-date-pickers/DatePicker'
import dayjs from 'dayjs'
import { Controller, FieldPath, useFormContext } from 'react-hook-form'

type Props = {
  name: FieldPath<Person>
  label: string

  /** Формат строки, как она хранится в `Person` (YYYY-MM-DD / YYYY-MM / YYYY) */
  format: string

  views?: DatePickerProps['views']
  required?: boolean
}

/** Дата в форме — обычная строка (`YYYY-MM[-DD]`), а не объект даты — MUI X работает с dayjs,
 *  поэтому конвертируем на границе через Controller (как SelectField, но для дат). */
export const DateField = ({ name, label, format, views, required }: Props) => {
  const { control } = useFormContext<Person>()

  return (
    <Controller
      name={name}
      control={control}
      render={({ field }) => (
        <DatePicker
          label={label}
          views={views}
          value={field.value ? dayjs(field.value as string) : null}
          onChange={(value) => field.onChange(value?.isValid() ? value.format(format) : '')}
          slotProps={{
            textField: { variant: 'standard', fullWidth: true, required, onBlur: field.onBlur },
          }}
        />
      )}
    />
  )
}
