import { useEffect, useMemo, useState } from 'react'
import { BackToResumes, LoadFailure, LoadPlaceholder } from '@features/curriculumVitae/components'
import { DEFAULT_SECTIONS } from '@features/curriculumVitae/sections'
import { getCurriculumVitaeInitialState } from '@features/curriculumVitae/store/curriculumVitae/curriculumVitae.state'
import { glassCard } from '@features/curriculumVitae/styles'
import { Person } from '@features/curriculumVitae/types'
import {
  clearCurriculumVitaeState,
  getPersonData,
  selectIsPersonDataLoading,
  selectPersonData,
  selectPersonDataError,
  updatePersonData,
} from '@features/curriculumVitae/store/curriculumVitae'
import { createLkResume } from '@features/curriculumVitae/store/lk'
import { VisibilityOutlined } from '@mui/icons-material'
import { Alert, Button, Container, Stack, Typography } from '@mui/material'
import { LocalizationProvider } from '@mui/x-date-pickers'
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs'
import { useLanguage, useTranslate } from '@shared/hooks'
import { useAppDispatch, useAppSelector } from '@shared/store'
import { FormProvider, useForm } from 'react-hook-form'
import { Link as RouterLink, useNavigate, useParams } from 'react-router'
import { HeroEditor } from './HeroEditor'
import { SectionsEditor } from './SectionsEditor'
import { ThemeEditor } from './ThemeEditor'

/** Пустое резюме для создания: те же секции, что и на публичной странице по умолчанию, все включены. */
const getBlankPerson = (): Person => ({
  ...getCurriculumVitaeInitialState().person,
  sections: [...DEFAULT_SECTIONS],
})

/**
 * Полное редактирование резюме: шапка (HeroEditor) + настраиваемые секции (SectionsEditor).
 * Без `cvId` в адресе (`/cv/lk/:userId/new`) — создание: форма пустая, «Сохранить» шлёт
 * `createLkResume` и переводит на обычный адрес редактирования уже с полученным id.
 */
export default function LkResumeEditScreen() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { userId, cvId } = useParams()
  const { language } = useLanguage()
  const translate = useTranslate('Lk.Edit')

  const isNew = !cvId

  const reduxIsLoading = useAppSelector(selectIsPersonDataLoading)
  const error = useAppSelector(selectPersonDataError)
  const reduxPerson = useAppSelector(selectPersonData)

  const [isSaving, setIsSaving] = useState(false)
  const [saveError, setSaveError] = useState(false)

  const blankPerson = useMemo(getBlankPerson, [])
  const isLoading = isNew ? false : reduxIsLoading
  const person = isNew ? blankPerson : reduxPerson

  const form = useForm<Person>({ values: person })

  useEffect(() => {
    if (!cvId || !userId) return

    const request = dispatch(getPersonData({ cvId, userId }))
    return () => {
      request.abort()
    }
  }, [dispatch, cvId, userId])

  useEffect(() => {
    return () => {
      dispatch(clearCurriculumVitaeState())
    }
  }, [dispatch])

  const onSubmit = form.handleSubmit(async (values) => {
    if (!userId) return

    setIsSaving(true)
    setSaveError(false)

    try {
      if (isNew) {
        const created = await dispatch(createLkResume({ userId, person: values })).unwrap()
        navigate(`/cv/lk/${encodeURIComponent(userId)}/${encodeURIComponent(created.id)}`, {
          replace: true,
        })
      } else if (cvId) {
        await dispatch(updatePersonData({ userId, cvId, person: values })).unwrap()
      }
    } catch {
      setSaveError(true)
    } finally {
      setIsSaving(false)
    }
  })

  if (isLoading) return <LoadPlaceholder />
  if (error) return <LoadFailure error={error} />

  return (
    <>
      {userId && <BackToResumes userId={userId} />}

      <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale={language}>
        <FormProvider {...form}>
          <Container maxWidth="md" sx={{ py: 6 }}>
            <Stack
              direction="row"
              alignItems="center"
              justifyContent="space-between"
              sx={{ mb: 3 }}
            >
              <Typography variant="h4">
                {person.position || person.name || translate('newResume')}
              </Typography>

              <Stack direction="row" spacing={1.5}>
                {!isNew && userId && cvId && (
                  <Button
                    component={RouterLink}
                    to={`/cv/lk/${encodeURIComponent(userId)}/${encodeURIComponent(cvId)}/preview`}
                    startIcon={<VisibilityOutlined />}
                    sx={{ ...glassCard, borderRadius: 999, px: 2, textTransform: 'none' }}
                  >
                    {translate('preview')}
                  </Button>
                )}

                <Button
                  variant="contained"
                  onClick={onSubmit}
                  disabled={isSaving}
                  sx={{ borderRadius: 999, px: 3, textTransform: 'none' }}
                >
                  {isSaving ? translate('saving') : translate('save')}
                </Button>
              </Stack>
            </Stack>

            {saveError && (
              <Alert severity="error" sx={{ mb: 3 }}>
                {translate('saveError')}
              </Alert>
            )}

            <Stack spacing={3}>
              <ThemeEditor />
              <HeroEditor />
              <SectionsEditor />
            </Stack>
          </Container>
        </FormProvider>
      </LocalizationProvider>
    </>
  )
}
