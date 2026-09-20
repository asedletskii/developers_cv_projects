export const sharedDict = {
  ru: {
    Header: {
      logout: 'Выход',

      Tabs: {
        main: 'Главная',
        about: 'О проекте',
      },
    },

    Footer: {
      info: 'Pet-проект, который агрегирует все web-проекты двух разработчиков',
      developers: 'Разработчики',

      Developers: {
        sedletskii: 'Седлецкий Антон',
        kosterin: 'Костерин Дмитрий',

        Roles: {
          frontend: 'Frontend разработчик',
          backend: 'Backend разработчик',
        },
      },
    },

    ContentContainer: {
      backButton: 'Назад',
    },

    Toasts: {
      copied: 'Скопировано',
      copyFailed: 'Не удалось скопировать',
      cannotGetPersonData: 'Не удалось загрузить резюме',
      cannotGetProject: 'Не удалось загрузить проект',
    },
  },
  en: {
    Header: {
      logout: 'Log out',

      Tabs: {
        main: 'Main',
        about: 'About',
      },
    },

    Footer: {
      info: 'A Pet project that aggregates all the web projects of two developers',
      developers: 'Developers',

      Developers: {
        sedletskii: 'Sedletskii Anton',
        kosterin: 'Kosterin Dmitry',

        Roles: {
          frontend: 'Frontend developer',
          backend: 'Backend developer',
        },
      },
    },

    ContentContainer: {
      backButton: '',
    },

    Toasts: {
      copied: 'Copied',
      copyFailed: "Couldn't copy",
      cannotGetPersonData: "Couldn't load the résumé",
      cannotGetProject: "Couldn't load the project",
    },
  },
} as const
