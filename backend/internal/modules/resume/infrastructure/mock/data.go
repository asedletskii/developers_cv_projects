package mock

import (
	"github.com/google/uuid"

	"developers-cv/backend/internal/modules/resume/application/queries"
)

func str(value string) *string { return &value }
func intPtr(value int) *int    { return &value }
func technology(name string, grade string) queries.Technology {
	result := queries.Technology{Name: name}
	if grade != "" {
		result.Grade = str(grade)
	}
	return result
}
func stack(names ...string) []queries.Technology {
	result := make([]queries.Technology, 0, len(names))
	for _, name := range names {
		result = append(result, technology(name, ""))
	}
	return result
}

func newFirstPerson() queries.Person {
	return queries.Person{
		ID:             PersonOneID,
		Name:           "Антон",
		Surname:        "Седлецкий",
		Patronymic:     str("Александрович"),
		Position:       "Frontend",
		PhotoPath:      "/mock/avatar.svg",
		ResumeLink:     str("/mock/resume.pdf"),
		HiddenContacts: []string{"email", "phone"},
		Availability:   &queries.Availability{Status: "open", From: str("2026-10-01")},
		DateOfBirth:    "1998-04-12",
		Location:       "Новосибирск",
		Relocation:     true,
		WorkFormat:     []string{"fullTime", "remote", "hybrid"},
		Contacts: queries.Contacts{
			Phone:    str("+79999999999"),
			Telegram: str("a_sedletskii"),
			Email:    str("anton@example.com"),
		},
		Sections: []queries.SectionSettings{
			{ID: "about", Enabled: true},
			{ID: "techStack", Enabled: true},
			{ID: "experience", Enabled: true},
			{ID: "skills", Enabled: true},
			{ID: "education", Enabled: true},
			{ID: "courses", Enabled: true, Collapsed: true},
			{ID: "projects", Enabled: true},
			{ID: "activities", Enabled: true, Collapsed: true},
			{ID: "contacts", Enabled: true},
		},
		Greeting: queries.Greeting{
			First:           "Привет!",
			Second:          "Меня зовут",
			SecondHighlight: "Седлецкий Антон",
			Third:           "Я Frontend-разработчик из Сибири",
		},
		Languages: []queries.PersonLanguage{
			{Language: "ru", Grade: str("C2")},
			{Language: "en", Grade: str("B2")},
			{Language: "de", Grade: str("A2")},
		},
		Experience: []queries.Experience{
			{
				ID:          "exp-1",
				Company:     "Северный ветер",
				Position:    "Senior Frontend разработчик",
				Location:    str("Удалённо"),
				StartDate:   "2023-03",
				Description: str("Разрабатываю клиентские приложения для финтех-продуктов и веду фронтенд-команду."),
				Achievements: []string{
					"Перевёл монолитный SPA на модульную архитектуру, время сборки сократилось в 3 раза",
					"Внедрил дизайн-систему на MUI, которую используют 4 продуктовые команды",
					"Довёл Lighthouse-показатели ключевых страниц с 62 до 95",
				},
				Stack: stack("React", "TypeScript", "Redux Toolkit", "MUI", "Vite", "GSAP")},
			{
				ID:          "exp-2",
				Company:     "Tech Labs",
				Position:    "Frontend разработчик",
				Location:    str("Новосибирск"),
				StartDate:   "2021-02",
				EndDate:     str("2023-02"),
				Description: str("Делал личный кабинет и административную панель для B2B-платформы."),
				Achievements: []string{
					"Написал библиотеку компонентов и переиспользуемых хуков на 40+ элементов",
					"Настроил CI: линтинг, тесты и предпросмотр каждой ветки",
					"Покрыл критичные сценарии интеграционными тестами",
				},
				Stack: stack("React", "TypeScript", "Jest", "Storybook", "GitHub Actions")},
			{
				ID:           "exp-3",
				Company:      "Freelance",
				Position:     "Web-разработчик",
				Location:     str("Удалённо"),
				StartDate:    "2019-06",
				EndDate:      str("2021-01"),
				Description:  str("Лендинги и небольшие интернет-магазины для малого бизнеса."),
				Achievements: []string{"Выполнил более 20 проектов от вёрстки до деплоя"},
				Stack:        stack("JavaScript", "HTML", "SCSS", "Node.js")},
			{
				ID:           "exp-4",
				Company:      "Студия «Мозаика»",
				Position:     "Junior Frontend разработчик",
				Location:     str("Новосибирск"),
				StartDate:    "2018-07",
				EndDate:      str("2019-05"),
				Description:  str("Вёрстка и доработка корпоративных сайтов в небольшой веб-студии."),
				Achievements: []string{"Перевёл верстку студии на адаптивную сетку и препроцессор"},
				Stack:        stack("JavaScript", "jQuery", "SCSS")},
			{
				ID:           "exp-5",
				Company:      "ИП Кузнецов",
				Position:     "Стажёр-разработчик",
				Location:     str("Новосибирск"),
				StartDate:    "2017-09",
				EndDate:      str("2018-06"),
				Achievements: []string{"Собрал первый коммерческий сайт и настроил хостинг"},
				Stack:        stack("HTML", "CSS", "PHP")},
		},
		Education: []queries.Education{
			{
				ID:          "edu-1",
				Location:    "Новосибирск",
				University:  "НГТУ",
				Grade:       3,
				ProgramName: "Программная инженерия",
				ProfileName: str("Разработка веб-приложений"),
				Faculty:     "Автоматики и вычислительной техники",
				Course:      intPtr(2),
			},
			{
				ID:          "edu-2",
				Location:    "Москва",
				University:  "МГТУ им. Баумана",
				Grade:       1,
				ProgramName: "Информатика и вычислительная техника",
				ProfileName: str("Разработка ПО"),
				Faculty:     "Информатика и системы управления",
				EndDate:     str("2021"),
			},
		},
		Courses: []queries.Course{
			{
				ID:              "course-1",
				Name:            "Advanced React и производительность",
				Organization:    "Frontend Masters",
				EndDate:         "2024-05",
				Hours:           intPtr(40),
				CertificateLink: str("https://example.com/certificate/react"),
			},
			{
				ID:              "course-2",
				Name:            "Архитектура фронтенд-приложений",
				Organization:    "Яндекс Практикум",
				EndDate:         "2023-11",
				Hours:           intPtr(72),
				CertificateLink: str("https://example.com/certificate/architecture"),
			},
			{
				ID:           "course-3",
				Name:         "Node.js: разработка серверных приложений",
				Organization: "Stepik",
				EndDate:      "2022-03",
				Hours:        intPtr(36),
			},
		},
		Activities: []queries.Activity{
			{
				ID:          "act-1",
				Type:        "openSource",
				Title:       "Контрибьютор в open-source",
				Description: str("Исправления и улучшения в популярных UI-библиотеках, 15+ принятых pull request"),
				Date:        str("2024-01"),
				Link:        str("https://github.com/example"),
			},
			{
				ID:          "act-2",
				Type:        "hackathon",
				Title:       "Хакатон «Цифровой прорыв»",
				Description: str("Финалист в команде из 4 человек: прототип сервиса за 48 часов"),
				Date:        str("2023-09"),
			},
			{
				ID:          "act-3",
				Type:        "speaker",
				Title:       "Доклад «Анимации на GSAP без боли»",
				Description: str("Выступление на митапе фронтенд-разработчиков Новосибирска"),
				Date:        str("2024-03"),
				Link:        str("https://example.com/talk"),
			},
			{
				ID:          "act-4",
				Type:        "mentoring",
				Title:       "Менторство джунов",
				Description: str("Помогаю начинающим разработчикам с ревью кода и карьерным планом"),
			},
			{
				ID:          "act-5",
				Type:        "community",
				Title:       "Организатор митапов",
				Description: str("Помогаю проводить встречи локального IT-сообщества"),
				Date:        str("2022-10"),
			},
		},
		HardSkills: []string{
			"architecture", "testing", "cicd", "performance", "accessibility", "codeReview", "uiDesign", "api",
		},
		SoftSkills: []string{
			"communication", "teamwork", "problemSolving", "timeManagement", "mentoring", "adaptability", "ownership",
		},
		Technologies: []queries.Technology{
			technology("TypeScript", "senior"),
			technology("React", "senior"),
			technology("Redux Toolkit", "middle"),
			technology("MUI", "senior"),
			technology("GSAP", "middle"),
			technology("Vite", "middle"),
			technology("Node.js", "middle"),
			technology("PostgreSQL", "junior"),
			technology("Docker", "junior"),
		},
		Projects: []queries.ProjectSummary{
			{
				ID:          "p1",
				Name:        "CV Platform",
				Description: str("Платформа для отображения резюме разработчиков"),
				Stack:       stack("React", "TypeScript", "Vite", "GSAP"),
				Link:        str("https://example.com"),
			},
			{
				ID:          "p2",
				Name:        "Finance Tracker",
				Description: str("Приложение для учёта личных финансов с аналитикой по категориям"),
				Stack:       stack("React", "Redux", "Node.js", "PostgreSQL"),
				Link:        str("https://example.com/finance"),
			},
			{
				ID:          "p3",
				Name:        "Pulse Dashboard",
				Description: str("Дашборд метрик в реальном времени: графики, алерты и фильтры"),
				Stack:       stack("React", "TypeScript", "WebSocket", "D3"),
				Link:        str("https://example.com/pulse"),
			},
			{
				ID:          "p4",
				Name:        "Design Kit",
				Description: str("Библиотека UI-компонентов с темами и документацией в Storybook"),
				Stack:       stack("React", "MUI", "Storybook"),
			},
		},
	}
}

func newSecondPerson(first queries.Person) queries.Person {
	person := clonePerson(first)
	person.ID = PersonTwoID
	person.Name = "Дмитрий"
	person.Surname = "Костерин"
	person.Patronymic = nil
	person.Position = "Backend"
	person.Location = "Томск"
	person.Contacts = queries.Contacts{Telegram: str("KostDmitry"), Email: str("dmitry@example.com")}
	person.Greeting = queries.Greeting{
		First:           "Здравствуйте!",
		Second:          "Я",
		SecondHighlight: "Дмитрий Костерин",
		Third:           "Backend-разработчик: API, базы данных и всё, что работает под капотом",
	}
	person.Theme = &queries.ThemeSettings{Primary: "#2E9E6B", Secondary: "#7C5CFA"}
	person.HiddenContacts = []string{}
	person.Availability = &queries.Availability{Status: "considering"}
	person.Courses = []queries.Course{}
	person.Activities = []queries.Activity{}
	person.Sections = []queries.SectionSettings{
		{ID: "about", Enabled: true},
		{ID: "experience", Enabled: true},
		{ID: "techStack", Enabled: true},
		{ID: "skills", Enabled: false},
		{ID: "education", Enabled: true, Collapsed: true},
		{ID: "projects", Enabled: true},
		{ID: "contacts", Enabled: true},
	}
	return person
}

func localizedGreeting(id uuid.UUID, lang queries.Language, fallback queries.Greeting) queries.Greeting {
	if lang != queries.LanguageEN {
		return fallback
	}
	switch id.String() {
	case PersonOneID:
		return queries.Greeting{
			First:           "Hello!",
			Second:          "My name is",
			SecondHighlight: "Anton Sedletskii",
			Third:           "I'm a Frontend developer from Siberia",
		}
	case PersonTwoID:
		return queries.Greeting{
			First:           "Hello!",
			Second:          "I'm",
			SecondHighlight: "Dmitry Kosterin",
			Third:           "Backend developer: APIs, databases and everything under the hood",
		}
	default:
		return fallback
	}
}

func localizedAvailabilityNote(id uuid.UUID, lang queries.Language) *string {
	if lang != queries.LanguageEN || id.String() != PersonOneID {
		return nil
	}
	return str("Looking for a Senior Frontend role, remote or hybrid")
}

func localizedContactsText(id uuid.UUID, lang queries.Language) *queries.ContactsText {
	if id.String() != PersonOneID {
		return nil
	}
	if lang == queries.LanguageEN {
		return &queries.ContactsText{
			Title:    str("Tell me about your project"),
			Subtitle: str("Let's discuss the task, timeline and how I can help"),
		}
	}
	return &queries.ContactsText{
		Title:    str("Расскажите о своём проекте"),
		Subtitle: str("Обсудим задачу, сроки и то, чем я могу быть полезен"),
	}
}
