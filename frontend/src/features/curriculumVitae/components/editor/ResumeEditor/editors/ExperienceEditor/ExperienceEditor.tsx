import { ArrayItemCard } from '@features/curriculumVitae/components/editor/form/ArrayItemCard'
import { DateField } from '@features/curriculumVitae/components/editor/form/DateField'
import { StringListField } from '@features/curriculumVitae/components/editor/form/StringListField'
import { TechnologyListField } from '@features/curriculumVitae/components/editor/form/TechnologyListField'
import { useFormList } from '@features/curriculumVitae/hooks'
import { blankExperience } from '@features/curriculumVitae/utils'
import { AddRounded } from '@mui/icons-material'
import { Button, Stack, TextField, Typography } from '@mui/material'
import { useTranslate } from '@shared/hooks'

export const ExperienceEditor = () => {
  const translate = useTranslate('Lk.Edit.Experience')
  const { register, fields, addItem, itemProps } = useFormList('experiences', blankExperience)

  return (
    <Stack spacing={1.5}>
      {fields.map((field, index) => (
        <ArrayItemCard key={field.id} {...itemProps(index)}>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <TextField
              label={translate('position')}
              {...register(`experiences.${index}.position`)}
              required
              fullWidth
            />
            <TextField
              label={translate('company')}
              {...register(`experiences.${index}.company`)}
              required
              fullWidth
            />
          </Stack>

          <TextField
            fullWidth
            label={translate('location')}
            {...register(`experiences.${index}.location`)}
          />

          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <DateField
              name={`experiences.${index}.startDate`}
              label={translate('startDate')}
              format="YYYY-MM"
              views={['year', 'month']}
              required
            />
            <DateField
              name={`experiences.${index}.endDate`}
              label={translate('endDate')}
              format="YYYY-MM"
              views={['year', 'month']}
            />
          </Stack>

          <TextField
            fullWidth
            label={translate('description')}
            {...register(`experiences.${index}.description`)}
            multiline
          />

          <Stack spacing={1}>
            <Typography variant="body2" sx={{ color: 'text.secondary' }}>
              {translate('achievements')}
            </Typography>
            <StringListField
              name={`experiences.${index}.achievements`}
              addLabel={translate('addAchievement')}
            />
          </Stack>

          <Stack spacing={1}>
            <Typography variant="body2" sx={{ color: 'text.secondary' }}>
              {translate('tasks')}
            </Typography>
            <StringListField name={`experiences.${index}.tasks`} addLabel={translate('addTask')} />
          </Stack>

          <TechnologyListField name={`experiences.${index}.stack`} />
        </ArrayItemCard>
      ))}

      <Button
        size="small"
        startIcon={<AddRounded />}
        onClick={addItem}
        sx={{ alignSelf: 'flex-start', textTransform: 'none' }}
      >
        {translate('addExperience')}
      </Button>
    </Stack>
  )
}
