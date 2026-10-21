export const aboutPageDict = {
  ru: {
    AboutPage: {
      Hero: {
        kicker: 'О проекте',
        title: 'Как устроен Developers CV',
        subtitle:
          'Открытый pet-проект на двух разработчиков: Go на бэкенде, React на фронтенде — чтобы собрать резюме, которым не стыдно поделиться. Вот что происходит под капотом.',
        cta: 'Попробовать',
        Mock: {
          name: 'Иван Петров',
          role: 'Frontend-разработчик',
          available: 'Открыт к предложениям',
        },
      },

      Steps: {
        title: 'Как это работает',
        subtitle: 'От входа до ссылки, которую можно отправить рекрутеру — четыре шага',
        step: 'Шаг {{number}}',

        login: {
          title: 'Входишь через Telegram',
          text: 'Без паролей и лишних форм — один клик в Telegram, и на бэке уже готова защищённая сессия по cookie.',
        },
        build: {
          title: 'Собираешь резюме по разделам',
          text: 'Профиль, опыт, стек, образование, проекты, активности — каждый раздел редактируется и сохраняется отдельно, без риска потерять остальное.',
        },
        preview: {
          title: 'Настраиваешь тему и смотришь предпросмотр',
          text: 'Акцентный цвет, порядок секций, что свернуть по умолчанию — изменения видно рядом с формой, ещё до сохранения.',
        },
        share: {
          title: 'Создаёшь персональные ссылки',
          text: 'Каждому адресату — своя ссылка со своей аналитикой переходов. Настоящий id резюме в ней не светится.',
        },
      },

      Features: {
        title: 'Из чего это сделано',
        subtitle: 'Коротко о том, что уже умеет сервис',

        auth: { title: 'Вход через Telegram', text: 'OAuth-авторизация, никаких паролей и анкет' },
        sections: {
          title: 'Конструктор по разделам',
          text: 'Десяток независимых секций — заполняй в своём порядке',
        },
        livePreview: {
          title: 'Живой предпросмотр',
          text: 'Правки в форме сразу видно в итоговом виде резюме',
        },
        links: {
          title: 'Персональные ссылки',
          text: 'Своя ссылка на каждого адресата и аналитика по ней',
        },
        theme: {
          title: 'Гибкая тема',
          text: 'Свой акцентный цвет на каждое резюме, без привязки к общей палитре сайта',
        },
        print: {
          title: 'Версия для печати',
          text: 'Один клик — и резюме готово для распечатки или PDF',
        },
        analytics: {
          title: 'Аналитика ссылок',
          text: 'Просмотры, дочитывание, клики и источники — по каждой ссылке',
        },
        protection: {
          title: 'Защита контактов',
          text: 'Скрытые контакты открываются после проверки Cloudflare',
        },
        download: {
          title: 'Скачать как файл',
          text: 'Резюме собирается в HTML прямо в браузере и работает офлайн',
        },
      },

      Stats: {
        sections: 'секций в резюме',
        languages: 'языка резюме: RU и EN',
        contacts: 'способа связи под защитой',
        read: 'секунд — и визит засчитан прочтением',
      },

      Showcase: {
        title: 'Что внутри',
        subtitle: 'Четыре вещи, ради которых это затевалось',

        builder: {
          tab: 'Конструктор',
          title: 'Собирай в своём порядке',
          text: 'Девять независимых секций: включай нужные, меняй порядок и выбирай, какие свернуть при открытии.',
          p1: 'Каждая секция сохраняется отдельно',
          p2: 'Живой предпросмотр рядом с формой',
          p3: 'Свой акцентный цвет на каждое резюме',
        },
        analytics: {
          tab: 'Аналитика',
          title: 'Видно, что читают',
          text: 'Для каждой ссылки — просмотры, время на странице, до какой секции дочитали и куда кликали.',
          p1: 'Воронка по секциям и карта кликов',
          p2: 'Источники и устройства посетителей',
          p3: 'Анонимно: имена и почты не собираются',
        },
        protection: {
          tab: 'Защита',
          title: 'Контакты не светятся ботам',
          text: 'Телефон, почта и Telegram можно скрыть до клика: значение отдаётся только после проверки Cloudflare.',
          p1: 'Скрытые контакты не попадают в публичные данные',
          p2: 'Проверка «не робот» — только если Cloudflare сочтёт нужной',
          p3: 'Раскрытие и копирование контакта видны в аналитике',
        },
        download: {
          tab: 'Скачать',
          title: 'Резюме одним файлом',
          text: 'Страница собирается в самодостаточный HTML прямо в браузере: без скриптов и внешних запросов.',
          p1: 'Язык файла — язык резюме',
          p2: 'Работает офлайн и открывается где угодно',
          p3: 'Версия для печати и PDF — в один клик',
        },
      },

      Faq: {
        title: 'Частые вопросы',
        subtitle: 'Коротко о том, что обычно спрашивают',

        data: {
          q: 'Что собирает аналитика?',
          a: 'Время на странице, просмотренные разделы, клики и источник перехода. Имена и почты не собираются, IP-адрес в открытом виде не хранится.',
        },
        dnt: {
          q: 'Можно ли отключить сбор статистики?',
          a: 'Да: если в браузере включено «Не отслеживать», подробная статистика не собирается.',
        },
        hidden: {
          q: 'Как работает скрытие контактов?',
          a: 'Автор выбирает, какие контакты прятать. Посетитель нажимает «Показать», и после проверки Cloudflare значение приходит с сервера.',
        },
        languages: {
          q: 'На каких языках можно сделать резюме?',
          a: 'На русском и английском: у одного пользователя по одному резюме на каждый язык.',
        },
        links: {
          q: 'Что видит получатель ссылки?',
          a: 'Свою персональную ссылку без настоящего id резюме. Выключить её можно в любой момент в кабинете.',
        },
      },

      Stack: {
        title: 'На чём это работает',
      },

      Team: {
        title: 'Кто это сделал',
        subtitle: 'Два человека, один pet-проект',
      },

      Cta: {
        title: 'Готовы попробовать?',
        text: 'Соберите резюме за 10 минут и получите ссылку, которую не стыдно отправить рекрутеру.',
        button: 'Создать резюме',
      },
    },
  },

  en: {
    AboutPage: {
      Hero: {
        kicker: 'About',
        title: 'How Developers\u00A0CV works',
        subtitle:
          'An open pet project by two developers: Go on the backend, React on the frontend — built to make a résumé worth sharing. Here is what happens under the hood.',
        cta: 'Try it',
        Mock: {
          name: 'Ivan Petrov',
          role: 'Frontend developer',
          available: 'Open to offers',
        },
      },

      Steps: {
        title: 'How it works',
        subtitle: 'From sign-in to a link you can send a recruiter — four steps',
        step: 'Step {{number}}',

        login: {
          title: 'Sign in with Telegram',
          text: 'No passwords, no forms — one click in Telegram, and the backend already has a secure cookie session.',
        },
        build: {
          title: 'Build your résumé by section',
          text: 'Profile, experience, stack, education, projects, activities — every section is edited and saved on its own, with no risk to the rest.',
        },
        preview: {
          title: 'Tune the theme and watch the preview',
          text: 'Accent color, section order, what collapses by default — changes show up next to the form before you even save.',
        },
        share: {
          title: 'Create personal links',
          text: 'Every recipient gets their own link with its own view analytics. The real résumé id never shows up in it.',
        },
      },

      Features: {
        title: 'What it is built from',
        subtitle: 'A quick look at what the service can already do',

        auth: { title: 'Telegram sign-in', text: 'OAuth authentication, no passwords or forms' },
        sections: {
          title: 'Section-based builder',
          text: 'A dozen independent sections — fill them in any order',
        },
        livePreview: {
          title: 'Live preview',
          text: 'Form edits show up in the final résumé right away',
        },
        links: {
          title: 'Personal links',
          text: 'A dedicated link per recipient, each with its own analytics',
        },
        theme: {
          title: 'Flexible theme',
          text: 'A custom accent color per résumé, independent of the site palette',
        },
        print: {
          title: 'Print-ready version',
          text: 'One click and the résumé is ready to print or save as PDF',
        },
        analytics: {
          title: 'Link analytics',
          text: 'Views, read depth, clicks and sources — for every link',
        },
        protection: {
          title: 'Contact protection',
          text: 'Hidden contacts open after a Cloudflare check',
        },
        download: {
          title: 'Download as a file',
          text: 'The résumé is built into HTML right in the browser and works offline',
        },
      },

      Stats: {
        sections: 'résumé sections',
        languages: 'résumé languages: RU and EN',
        contacts: 'protected contact methods',
        read: 'seconds — and a visit counts as read',
      },

      Showcase: {
        title: 'What is inside',
        subtitle: 'Four things this was built for',

        builder: {
          tab: 'Builder',
          title: 'Build it your way',
          text: 'Nine independent sections: switch on the ones you need, reorder them and choose what collapses on open.',
          p1: 'Every section is saved on its own',
          p2: 'Live preview next to the form',
          p3: 'Your own accent color for each résumé',
        },
        analytics: {
          tab: 'Analytics',
          title: 'See what gets read',
          text: 'For every link — views, time on page, how far people read and where they clicked.',
          p1: 'Section funnel and a click map',
          p2: 'Visitor sources and devices',
          p3: 'Anonymous: names and emails are not collected',
        },
        protection: {
          tab: 'Protection',
          title: 'Contacts stay hidden from bots',
          text: 'Phone, email and Telegram can stay hidden until a click: the value is released only after a Cloudflare check.',
          p1: 'Hidden contacts never appear in public data',
          p2: 'The “not a robot” check shows up only when Cloudflare wants it',
          p3: 'Reveals and copies of a contact show up in analytics',
        },
        download: {
          tab: 'Download',
          title: 'Your résumé as one file',
          text: 'The page is built into a self-contained HTML file right in the browser: no scripts, no external requests.',
          p1: 'The file language is the résumé language',
          p2: 'Works offline and opens anywhere',
          p3: 'Print and PDF version in one click',
        },
      },

      Faq: {
        title: 'FAQ',
        subtitle: 'Quick answers to what people usually ask',

        data: {
          q: 'What does analytics collect?',
          a: 'Time on page, sections viewed, clicks and the referring source. Names and emails are not collected, and IP addresses are not stored in plain form.',
        },
        dnt: {
          q: 'Can I opt out of statistics?',
          a: 'Yes: if “Do Not Track” is on in your browser, detailed statistics are not collected.',
        },
        hidden: {
          q: 'How does hiding contacts work?',
          a: 'The author chooses which contacts to hide. A visitor presses “Show”, and after a Cloudflare check the value comes from the server.',
        },
        languages: {
          q: 'Which languages can a résumé be in?',
          a: 'Russian and English: one résumé per language for each user.',
        },
        links: {
          q: 'What does a link recipient see?',
          a: 'Their own personal link without the real résumé id. You can switch it off any time in your account.',
        },
      },

      Stack: {
        title: 'What it runs on',
      },

      Team: {
        title: 'Who built this',
        subtitle: 'Two people, one pet project',
      },

      Cta: {
        title: 'Ready to try it?',
        text: 'Put together a résumé in 10 minutes and get a link worth sending to a recruiter.',
        button: 'Create a résumé',
      },
    },
  },
}
