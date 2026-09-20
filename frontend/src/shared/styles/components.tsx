import { ThemeOptions } from '@mui/material'
import { alpha, Theme } from '@mui/material/styles'
import SF_Pro_Bold from '@shared/assets/fonts/SFProText-Heavy.woff2'
import SF_Pro_Regular from '@shared/assets/fonts/SFProText-Regular.woff2'
import SF_Pro_Light from '@shared/assets/fonts/SFProText-Thin.woff2'

/**
 * Всплывающая панель (меню, выпадающие списки): тот же полупрозрачный «стеклянный» материал, что у
 * карточек страницы. При reduced transparency становится плотной.
 */
const glassPaper = (theme: Theme) => {
  const isLight = theme.palette.mode === 'light'

  return {
    backgroundImage: 'none',
    borderRadius: 16,
    border: `1px solid ${theme.palette.divider}`,
    backgroundColor: alpha(theme.palette.surface.light, isLight ? 0.82 : 0.72),
    backdropFilter: 'blur(20px) saturate(180%)',
    WebkitBackdropFilter: 'blur(20px) saturate(180%)',
    boxShadow: `inset 0 1px 0 ${alpha(theme.palette.common.white, isLight ? 0.9 : 0.1)}, 0 12px 40px ${alpha(theme.palette.common.black, isLight ? 0.08 : 0.4)}`,

    '@media (prefers-reduced-transparency: reduce)': {
      backgroundColor: theme.palette.surface.light,
      backdropFilter: 'none',
      WebkitBackdropFilter: 'none',
    },
  }
}

/** Пункт меню или списка: скруглённый, с мягкой подсветкой при наведении и выборе. */
const glassOption = (theme: Theme) => ({
  borderRadius: 10,
  minHeight: 40,
  transition: 'background-color .15s',
  '&:hover, &.Mui-focusVisible': { backgroundColor: theme.palette.action.hover },
  '&.Mui-selected, &[aria-selected="true"]': {
    backgroundColor: alpha(theme.palette.primary.main, 0.16),
  },
  '&.Mui-selected:hover': { backgroundColor: alpha(theme.palette.primary.main, 0.24) },
})

export const components: ThemeOptions['components'] = {
  MuiMenu: {
    styleOverrides: {
      paper: ({ theme }) => ({ ...glassPaper(theme), marginTop: theme.spacing(0.5) }),
      list: ({ theme }) => ({
        padding: theme.spacing(0.75),
        '& .MuiDivider-root': { margin: theme.spacing(0.5, 0.75) },
      }),
    },
  },

  MuiMenuItem: {
    styleOverrides: {
      root: ({ theme }) => glassOption(theme),
    },
  },

  MuiAutocomplete: {
    styleOverrides: {
      paper: ({ theme }) => glassPaper(theme),
      listbox: ({ theme }) => ({ padding: theme.spacing(0.75) }),
      option: ({ theme }) => glassOption(theme),
    },
  },

  MuiCssBaseline: {
    styleOverrides: (theme) => ({
      '[data-color-scheme="dark"]': {
        colorScheme: 'dark',
      },

      '[data-color-scheme="light"]': {
        colorScheme: 'light',
      },

      '#root': {
        backgroundColor: theme.palette.surface.main,
      },

      '@font-face': [
        {
          fontFamily: 'SF Pro Light',
          src: `url(${SF_Pro_Light}) format('woff2')`,
        },
        {
          fontFamily: 'SF Pro Regular',
          src: `url(${SF_Pro_Regular}) format('woff2')`,
        },
        {
          fontFamily: 'SF Pro Bold',
          src: `url(${SF_Pro_Bold}) format('woff2')`,
        },
      ],
    }),
  },

  // Реакция на нажатие мгновенная (pointerdown). Используем свойство scale, а не transform,
  // чтобы не конфликтовать с transform от GSAP.
  MuiButtonBase: {
    styleOverrides: {
      root: {
        transition: 'scale 120ms ease-out',
        '&:active': { scale: '0.97' },
      },
    },
  },

  MuiButton: {
    styleOverrides: {
      root: ({ theme }) => ({
        transition: [
          theme.transitions.create(['background-color', 'box-shadow', 'border-color', 'color'], {
            duration: theme.transitions.duration.shortest,
          }),
          'scale 120ms ease-out',
        ].join(', '),
      }),

      textPrimary: ({ theme }) =>
        theme.palette.mode === 'light'
          ? {
              color: theme.palette.link.dark,
              '&:hover': { backgroundColor: alpha(theme.palette.link.main, 0.1) },
            }
          : {},
      outlinedPrimary: ({ theme }) =>
        theme.palette.mode === 'light'
          ? {
              color: theme.palette.link.dark,
              borderColor: alpha(theme.palette.link.main, 0.5),
              '&:hover': {
                borderColor: theme.palette.link.main,
                backgroundColor: alpha(theme.palette.link.main, 0.1),
              },
            }
          : {},
    },
  },

  MuiTypography: {
    defaultProps: {
      variant: 'body1',
    },
    styleOverrides: {
      root({ theme }) {
        return { color: theme.palette.text.primary }
      },
    },
  },

  MuiTextField: {
    defaultProps: {
      variant: 'standard',
      size: 'small',
      slotProps: { inputLabel: { shrink: true } },
    },
  },

  MuiTab: {
    styleOverrides: {
      root: ({ theme }) => ({
        textTransform: 'none',
        minHeight: 48,
        ...theme.typography.body1,
        color: theme.palette.text.primary,

        paddingInline: theme.spacing(0),
        [theme.breakpoints.up('sm')]: {
          paddingInline: theme.spacing(2),
        },

        '&.Mui-selected': {
          color: theme.palette.text.primary,
        },
      }),
    },
  },
}
