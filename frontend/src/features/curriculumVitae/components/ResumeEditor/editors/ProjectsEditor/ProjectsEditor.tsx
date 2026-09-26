import { ArrayItemCard } from '@features/curriculumVitae/components/form/ArrayItemCard'
import { TechnologyListField } from '@features/curriculumVitae/components/form/TechnologyListField'
import { useFormList } from '@features/curriculumVitae/hooks'
import { blankProject } from '@features/curriculumVitae/utils'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField } from '@mui/material'
import { useTranslate } from '@shared/hooks'

/** Секция «Проекты». */
export const ProjectsEditor = () => {
  const translate = useTranslate('Lk.Edit.Projects')
  const { register, fields, addItem, itemProps } = useFormList('projects', blankProject)

  return (
    <Stack spacing={1.5}>
      {fields.map((field, index) => (
        <ArrayItemCard key={field.id} {...itemProps(index)}>
          <TextField
            label={translate('name')}
            {...register(`projects.${index}.name`)}
            required
            fullWidth
          />
          <TextField
            fullWidth
            label={translate('description')}
            {...register(`projects.${index}.description`)}
            multiline
          />
          <TextField fullWidth label={translate('link')} {...register(`projects.${index}.link`)} />

          <TechnologyListField name={`projects.${index}.stack`} />
        </ArrayItemCard>
      ))}

      <Button
        size="small"
        startIcon={<AddRounded />}
        onClick={addItem}
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addProject')}
      </Button>
    </Stack>
  )
}
