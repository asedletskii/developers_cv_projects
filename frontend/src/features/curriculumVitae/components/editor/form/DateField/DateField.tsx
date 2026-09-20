import { Resume } from '@features/curriculumVitae/types'
import { DatePicker, DatePickerProps } from '@mui/x-date-pickers/DatePicker'
import dayjs from 'dayjs'
import { Controller, FieldPath, useFormContext } from 'react-hook-form'

type Props = {
  name: FieldPath<Resume>
  label: string
  format: string
  views?: DatePickerProps['views']
  required?: boolean
}

export const DateField = ({ name, label, format, views, required }: Props) => {
  const { control } = useFormContext<Resume>()

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
