import { AddRounded } from '@mui/icons-material'
import { Button, Stack } from '@mui/material'
import { useTranslate } from '@shared/hooks'
import { useLanguagesEditor } from './hooks'
import { LanguageRow } from './LanguageRow'

/** Вкладка «Языки»: какими языками владеет человек и на каком уровне. */
export const LanguagesEditor = () => {
  const translate = useTranslate('Lk.Edit.About')
  const editor = useLanguagesEditor()

  return (
    <Stack spacing={1}>
      {editor.fields.map((field, index) => (
        <LanguageRow
          key={field.id}
          index={index}
          languageOptions={editor.getLanguageOptions(index)}
          canMoveUp={editor.canMoveUp(index)}
          canMoveDown={editor.canMoveDown(index)}
          onMoveUp={() => editor.moveUp(index)}
          onMoveDown={() => editor.moveDown(index)}
          onRemove={() => editor.remove(index)}
        />
      ))}

      <Button
        size="small"
        startIcon={<AddRounded />}
        disabled={!editor.canAdd}
        onClick={editor.addLanguage}
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addLanguage')}
      </Button>
    </Stack>
  )
}
