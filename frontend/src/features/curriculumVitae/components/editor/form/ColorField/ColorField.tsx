import { glassCard } from '@features/curriculumVitae/styles'
import { Resume } from '@features/curriculumVitae/types'
import { Box, ClickAwayListener, Popper, Stack, TextField } from '@mui/material'
import { HexColorPicker } from 'react-colorful'
import { FieldPath } from 'react-hook-form'
import { useColorField } from './hooks'

type Props = {
  name: FieldPath<Resume>
  label: string
  defaultColor: string
}

export const ColorField = ({ name, label, defaultColor }: Props) => {
  const { color, anchorEl, openPicker, closePicker, onChange, onBlur } = useColorField(
    name,
    defaultColor,
  )

  return (
    <>
      <Stack direction="row" alignItems="center" spacing={1.5}>
        <Box
          role="button"
          aria-label={label}
          onClick={(e) => openPicker(e.currentTarget)}
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
          value={color}
          onChange={(e) => onChange(e.target.value)}
          onBlur={onBlur}
          sx={{ width: 140 }}
        />
      </Stack>

      <Popper open={!!anchorEl} anchorEl={anchorEl} placement="bottom-start" sx={{ zIndex: 1300 }}>
        <ClickAwayListener onClickAway={closePicker}>
          <Box sx={{ p: 1.5, mt: 1, ...glassCard }}>
            <HexColorPicker color={color} onChange={onChange} />
          </Box>
        </ClickAwayListener>
      </Popper>
    </>
  )
}
