import { useState } from 'react'
import { glassCard } from '@features/curriculumVitae/styles'
import { Person } from '@features/curriculumVitae/types'
import { Box, ClickAwayListener, Popper, Stack, TextField } from '@mui/material'
import { HexColorPicker } from 'react-colorful'
import { Controller, FieldPath, useFormContext } from 'react-hook-form'

type Props = {
  name: FieldPath<Person>
  label: string

  /** Цвет текущей темы сайта — показываем его, пока пользователь не задал свой (поле не пустое,
   *  но само по себе оно ничего не переопределяет: значение в форме остаётся пустым). */
  defaultColor: string
}

/** Цвет темы: попап с react-colorful (проще и предсказуемее нативного `<input type="color">`,
 *  чья OS-палитра нестабильно рендерится в некоторых окружениях) + поле с hex для точного ввода. */
export const ColorField = ({ name, label, defaultColor }: Props) => {
  const { control } = useFormContext<Person>()
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null)

  return (
    <Controller
      name={name}
      control={control}
      render={({ field }) => {
        const color = (field.value as string) || defaultColor

        return (
          <>
            <Stack direction="row" alignItems="center" spacing={1.5}>
              <Box
                role="button"
                aria-label={label}
                onClick={(e) => setAnchorEl(e.currentTarget)}
                sx={{
                  width: 36,
                  height: 36,
                  borderRadius: '50%',
                  cursor: 'pointer',
                  border: 1,
                  borderColor: 'divider',
                  bgcolor: color,
                  flexShrink: 0,
                }}
              />

              <TextField
                label={label}
                value={field.value || color}
                onChange={(e) => field.onChange(e.target.value)}
                onBlur={field.onBlur}
                sx={{ width: 140 }}
              />
            </Stack>

            <Popper
              open={!!anchorEl}
              anchorEl={anchorEl}
              placement="bottom-start"
              sx={{ zIndex: 1300 }}
            >
              <ClickAwayListener onClickAway={() => setAnchorEl(null)}>
                <Box sx={{ p: 1.5, mt: 1, ...glassCard }}>
                  <HexColorPicker color={color} onChange={(next) => field.onChange(next)} />
                </Box>
              </ClickAwayListener>
            </Popper>
          </>
        )
      }}
    />
  )
}
