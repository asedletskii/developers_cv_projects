-- Демо-резюме: те же два человека, что раньше отдавал мок фронта и API
-- (см. frontend/src/features/curriculumVitae/api/curriculumVitae.mock.ts и
-- backend/internal/modules/resume/infrastructure/mock/data.go), только теперь в Postgres.

INSERT INTO people (id, name, surname, patronymic, position, resume_link, photo_path, date_of_birth, location, relocation) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Антон', 'Седлецкий', 'Александрович', 'Frontend', '/mock/resume.pdf', '/mock/avatar.svg', '1998-04-12', 'Новосибирск', TRUE),
    ('550e8400-e29b-41d4-a716-446655440002', 'Дмитрий', 'Костерин', NULL, 'Backend', '/mock/resume.pdf', '/mock/avatar.svg', '1998-04-12', 'Томск', TRUE);

INSERT INTO person_greetings (person_id, first_ru, first_en, second_ru, second_en, second_highlight_ru, second_highlight_en, third_ru, third_en) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Привет!', 'Hello!', 'Меня зовут', 'My name is', 'Седлецкий Антон', 'Anton Sedletskii', 'Я Frontend-разработчик из Сибири', 'I''m a Frontend developer from Siberia'),
    ('550e8400-e29b-41d4-a716-446655440002', 'Здравствуйте!', 'Hello!', 'Я', 'I''m', 'Дмитрий Костерин', 'Dmitry Kosterin', 'Backend-разработчик: API, базы данных и всё, что работает под капотом', 'Backend developer: APIs, databases and everything under the hood');

INSERT INTO person_availability (person_id, status, available_from, note_ru, note_en) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'open', '2026-10-01', NULL, 'Looking for a Senior Frontend role, remote or hybrid'),
    ('550e8400-e29b-41d4-a716-446655440002', 'considering', NULL, NULL, NULL);

INSERT INTO person_themes (person_id, primary_color, secondary_color) VALUES
    ('550e8400-e29b-41d4-a716-446655440002', '#2E9E6B', '#7C5CFA');

INSERT INTO person_contacts_text (person_id, title_ru, title_en, subtitle_ru, subtitle_en) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Расскажите о своём проекте', 'Tell me about your project', 'Обсудим задачу, сроки и то, чем я могу быть полезен', 'Let''s discuss the task, timeline and how I can help');

INSERT INTO person_contacts (person_id, phone, telegram, email) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', '+79999999999', 'a_sedletskii', 'anton@example.com'),
    ('550e8400-e29b-41d4-a716-446655440002', NULL, 'KostDmitry', 'dmitry@example.com');

INSERT INTO person_sections (person_id, section_id, position, enabled, collapsed) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'about', 0, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440001', 'techStack', 1, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440001', 'experience', 2, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440001', 'skills', 3, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440001', 'education', 4, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440001', 'courses', 5, TRUE, TRUE),
    ('550e8400-e29b-41d4-a716-446655440001', 'projects', 6, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440001', 'activities', 7, TRUE, TRUE),
    ('550e8400-e29b-41d4-a716-446655440001', 'contacts', 8, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440002', 'about', 0, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440002', 'experience', 1, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440002', 'techStack', 2, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440002', 'skills', 3, FALSE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440002', 'education', 4, TRUE, TRUE),
    ('550e8400-e29b-41d4-a716-446655440002', 'projects', 5, TRUE, FALSE),
    ('550e8400-e29b-41d4-a716-446655440002', 'contacts', 6, TRUE, FALSE);

INSERT INTO person_work_formats (person_id, value, position) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'fullTime', 0),
    ('550e8400-e29b-41d4-a716-446655440001', 'remote', 1),
    ('550e8400-e29b-41d4-a716-446655440001', 'hybrid', 2),
    ('550e8400-e29b-41d4-a716-446655440002', 'fullTime', 0),
    ('550e8400-e29b-41d4-a716-446655440002', 'remote', 1),
    ('550e8400-e29b-41d4-a716-446655440002', 'hybrid', 2);

INSERT INTO person_hidden_contacts (person_id, value, position) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'email', 0),
    ('550e8400-e29b-41d4-a716-446655440001', 'phone', 1);

INSERT INTO person_hard_skills (person_id, value, position) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'architecture', 0),
    ('550e8400-e29b-41d4-a716-446655440001', 'testing', 1),
    ('550e8400-e29b-41d4-a716-446655440001', 'cicd', 2),
    ('550e8400-e29b-41d4-a716-446655440001', 'performance', 3),
    ('550e8400-e29b-41d4-a716-446655440001', 'accessibility', 4),
    ('550e8400-e29b-41d4-a716-446655440001', 'codeReview', 5),
    ('550e8400-e29b-41d4-a716-446655440001', 'uiDesign', 6),
    ('550e8400-e29b-41d4-a716-446655440001', 'api', 7),
    ('550e8400-e29b-41d4-a716-446655440002', 'architecture', 0),
    ('550e8400-e29b-41d4-a716-446655440002', 'testing', 1),
    ('550e8400-e29b-41d4-a716-446655440002', 'cicd', 2),
    ('550e8400-e29b-41d4-a716-446655440002', 'performance', 3),
    ('550e8400-e29b-41d4-a716-446655440002', 'accessibility', 4),
    ('550e8400-e29b-41d4-a716-446655440002', 'codeReview', 5),
    ('550e8400-e29b-41d4-a716-446655440002', 'uiDesign', 6),
    ('550e8400-e29b-41d4-a716-446655440002', 'api', 7);

INSERT INTO person_soft_skills (person_id, value, position) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'communication', 0),
    ('550e8400-e29b-41d4-a716-446655440001', 'teamwork', 1),
    ('550e8400-e29b-41d4-a716-446655440001', 'problemSolving', 2),
    ('550e8400-e29b-41d4-a716-446655440001', 'timeManagement', 3),
    ('550e8400-e29b-41d4-a716-446655440001', 'mentoring', 4),
    ('550e8400-e29b-41d4-a716-446655440001', 'adaptability', 5),
    ('550e8400-e29b-41d4-a716-446655440001', 'ownership', 6),
    ('550e8400-e29b-41d4-a716-446655440002', 'communication', 0),
    ('550e8400-e29b-41d4-a716-446655440002', 'teamwork', 1),
    ('550e8400-e29b-41d4-a716-446655440002', 'problemSolving', 2),
    ('550e8400-e29b-41d4-a716-446655440002', 'timeManagement', 3),
    ('550e8400-e29b-41d4-a716-446655440002', 'mentoring', 4),
    ('550e8400-e29b-41d4-a716-446655440002', 'adaptability', 5),
    ('550e8400-e29b-41d4-a716-446655440002', 'ownership', 6);

INSERT INTO person_languages (person_id, language, grade, position) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'ru', 'C2', 0),
    ('550e8400-e29b-41d4-a716-446655440001', 'en', 'B2', 1),
    ('550e8400-e29b-41d4-a716-446655440001', 'de', 'A2', 2),
    ('550e8400-e29b-41d4-a716-446655440002', 'ru', 'C2', 0),
    ('550e8400-e29b-41d4-a716-446655440002', 'en', 'B2', 1),
    ('550e8400-e29b-41d4-a716-446655440002', 'de', 'A2', 2);

INSERT INTO person_technologies (person_id, name, grade, position) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'TypeScript', 'senior', 0),
    ('550e8400-e29b-41d4-a716-446655440001', 'React', 'senior', 1),
    ('550e8400-e29b-41d4-a716-446655440001', 'Redux Toolkit', 'middle', 2),
    ('550e8400-e29b-41d4-a716-446655440001', 'MUI', 'senior', 3),
    ('550e8400-e29b-41d4-a716-446655440001', 'GSAP', 'middle', 4),
    ('550e8400-e29b-41d4-a716-446655440001', 'Vite', 'middle', 5),
    ('550e8400-e29b-41d4-a716-446655440001', 'Node.js', 'middle', 6),
    ('550e8400-e29b-41d4-a716-446655440001', 'PostgreSQL', 'junior', 7),
    ('550e8400-e29b-41d4-a716-446655440001', 'Docker', 'junior', 8),
    ('550e8400-e29b-41d4-a716-446655440002', 'TypeScript', 'senior', 0),
    ('550e8400-e29b-41d4-a716-446655440002', 'React', 'senior', 1),
    ('550e8400-e29b-41d4-a716-446655440002', 'Redux Toolkit', 'middle', 2),
    ('550e8400-e29b-41d4-a716-446655440002', 'MUI', 'senior', 3),
    ('550e8400-e29b-41d4-a716-446655440002', 'GSAP', 'middle', 4),
    ('550e8400-e29b-41d4-a716-446655440002', 'Vite', 'middle', 5),
    ('550e8400-e29b-41d4-a716-446655440002', 'Node.js', 'middle', 6),
    ('550e8400-e29b-41d4-a716-446655440002', 'PostgreSQL', 'junior', 7),
    ('550e8400-e29b-41d4-a716-446655440002', 'Docker', 'junior', 8);

INSERT INTO experiences (id, person_id, position, company, job_position, location, start_date, end_date, description) VALUES
    ('4c6c35da-9165-407e-b44d-7db5843188a2', '550e8400-e29b-41d4-a716-446655440001', 0, 'Северный ветер', 'Senior Frontend разработчик', 'Удалённо', '2023-03-01', NULL, 'Разрабатываю клиентские приложения для финтех-продуктов и веду фронтенд-команду.'),
    ('5cacc2c0-a10e-4fc5-a045-34a7bffa796b', '550e8400-e29b-41d4-a716-446655440001', 1, 'Tech Labs', 'Frontend разработчик', 'Новосибирск', '2021-02-01', '2023-02-01', 'Делал личный кабинет и административную панель для B2B-платформы.'),
    ('20abd798-80a1-4723-b69a-7f241fc12a9b', '550e8400-e29b-41d4-a716-446655440001', 2, 'Freelance', 'Web-разработчик', 'Удалённо', '2019-06-01', '2021-01-01', 'Лендинги и небольшие интернет-магазины для малого бизнеса.'),
    ('55225d89-0e1d-44c9-ac06-df03e3c6bc32', '550e8400-e29b-41d4-a716-446655440001', 3, 'Студия «Мозаика»', 'Junior Frontend разработчик', 'Новосибирск', '2018-07-01', '2019-05-01', 'Вёрстка и доработка корпоративных сайтов в небольшой веб-студии.'),
    ('6a0b3015-c4ad-4f4e-9076-42577bfc1593', '550e8400-e29b-41d4-a716-446655440001', 4, 'ИП Кузнецов', 'Стажёр-разработчик', 'Новосибирск', '2017-09-01', '2018-06-01', NULL),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', '550e8400-e29b-41d4-a716-446655440002', 0, 'Северный ветер', 'Senior Frontend разработчик', 'Удалённо', '2023-03-01', NULL, 'Разрабатываю клиентские приложения для финтех-продуктов и веду фронтенд-команду.'),
    ('f4658d87-99ec-455a-96ea-e6452e8bff16', '550e8400-e29b-41d4-a716-446655440002', 1, 'Tech Labs', 'Frontend разработчик', 'Новосибирск', '2021-02-01', '2023-02-01', 'Делал личный кабинет и административную панель для B2B-платформы.'),
    ('9dccf1a9-c4ae-4bc5-9261-06cdcdee4678', '550e8400-e29b-41d4-a716-446655440002', 2, 'Freelance', 'Web-разработчик', 'Удалённо', '2019-06-01', '2021-01-01', 'Лендинги и небольшие интернет-магазины для малого бизнеса.'),
    ('eb5879d6-45df-4090-a882-d2a4daf51ec5', '550e8400-e29b-41d4-a716-446655440002', 3, 'Студия «Мозаика»', 'Junior Frontend разработчик', 'Новосибирск', '2018-07-01', '2019-05-01', 'Вёрстка и доработка корпоративных сайтов в небольшой веб-студии.'),
    ('4d846059-40e8-47b1-b477-f6bb43852de1', '550e8400-e29b-41d4-a716-446655440002', 4, 'ИП Кузнецов', 'Стажёр-разработчик', 'Новосибирск', '2017-09-01', '2018-06-01', NULL);

INSERT INTO experience_achievements (experience_id, position, value) VALUES
    ('4c6c35da-9165-407e-b44d-7db5843188a2', 0, 'Перевёл монолитный SPA на модульную архитектуру, время сборки сократилось в 3 раза'),
    ('4c6c35da-9165-407e-b44d-7db5843188a2', 1, 'Внедрил дизайн-систему на MUI, которую используют 4 продуктовые команды'),
    ('4c6c35da-9165-407e-b44d-7db5843188a2', 2, 'Довёл Lighthouse-показатели ключевых страниц с 62 до 95'),
    ('5cacc2c0-a10e-4fc5-a045-34a7bffa796b', 0, 'Написал библиотеку компонентов и переиспользуемых хуков на 40+ элементов'),
    ('5cacc2c0-a10e-4fc5-a045-34a7bffa796b', 1, 'Настроил CI: линтинг, тесты и предпросмотр каждой ветки'),
    ('5cacc2c0-a10e-4fc5-a045-34a7bffa796b', 2, 'Покрыл критичные сценарии интеграционными тестами'),
    ('20abd798-80a1-4723-b69a-7f241fc12a9b', 0, 'Выполнил более 20 проектов от вёрстки до деплоя'),
    ('55225d89-0e1d-44c9-ac06-df03e3c6bc32', 0, 'Перевёл верстку студии на адаптивную сетку и препроцессор'),
    ('6a0b3015-c4ad-4f4e-9076-42577bfc1593', 0, 'Собрал первый коммерческий сайт и настроил хостинг'),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', 0, 'Перевёл монолитный SPA на модульную архитектуру, время сборки сократилось в 3 раза'),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', 1, 'Внедрил дизайн-систему на MUI, которую используют 4 продуктовые команды'),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', 2, 'Довёл Lighthouse-показатели ключевых страниц с 62 до 95'),
    ('f4658d87-99ec-455a-96ea-e6452e8bff16', 0, 'Написал библиотеку компонентов и переиспользуемых хуков на 40+ элементов'),
    ('f4658d87-99ec-455a-96ea-e6452e8bff16', 1, 'Настроил CI: линтинг, тесты и предпросмотр каждой ветки'),
    ('f4658d87-99ec-455a-96ea-e6452e8bff16', 2, 'Покрыл критичные сценарии интеграционными тестами'),
    ('9dccf1a9-c4ae-4bc5-9261-06cdcdee4678', 0, 'Выполнил более 20 проектов от вёрстки до деплоя'),
    ('eb5879d6-45df-4090-a882-d2a4daf51ec5', 0, 'Перевёл верстку студии на адаптивную сетку и препроцессор'),
    ('4d846059-40e8-47b1-b477-f6bb43852de1', 0, 'Собрал первый коммерческий сайт и настроил хостинг');

INSERT INTO experience_technologies (experience_id, position, name, grade) VALUES
    ('4c6c35da-9165-407e-b44d-7db5843188a2', 0, 'React', NULL),
    ('4c6c35da-9165-407e-b44d-7db5843188a2', 1, 'TypeScript', NULL),
    ('4c6c35da-9165-407e-b44d-7db5843188a2', 2, 'Redux Toolkit', NULL),
    ('4c6c35da-9165-407e-b44d-7db5843188a2', 3, 'MUI', NULL),
    ('4c6c35da-9165-407e-b44d-7db5843188a2', 4, 'Vite', NULL),
    ('4c6c35da-9165-407e-b44d-7db5843188a2', 5, 'GSAP', NULL),
    ('5cacc2c0-a10e-4fc5-a045-34a7bffa796b', 0, 'React', NULL),
    ('5cacc2c0-a10e-4fc5-a045-34a7bffa796b', 1, 'TypeScript', NULL),
    ('5cacc2c0-a10e-4fc5-a045-34a7bffa796b', 2, 'Jest', NULL),
    ('5cacc2c0-a10e-4fc5-a045-34a7bffa796b', 3, 'Storybook', NULL),
    ('5cacc2c0-a10e-4fc5-a045-34a7bffa796b', 4, 'GitHub Actions', NULL),
    ('20abd798-80a1-4723-b69a-7f241fc12a9b', 0, 'JavaScript', NULL),
    ('20abd798-80a1-4723-b69a-7f241fc12a9b', 1, 'HTML', NULL),
    ('20abd798-80a1-4723-b69a-7f241fc12a9b', 2, 'SCSS', NULL),
    ('20abd798-80a1-4723-b69a-7f241fc12a9b', 3, 'Node.js', NULL),
    ('55225d89-0e1d-44c9-ac06-df03e3c6bc32', 0, 'JavaScript', NULL),
    ('55225d89-0e1d-44c9-ac06-df03e3c6bc32', 1, 'jQuery', NULL),
    ('55225d89-0e1d-44c9-ac06-df03e3c6bc32', 2, 'SCSS', NULL),
    ('6a0b3015-c4ad-4f4e-9076-42577bfc1593', 0, 'HTML', NULL),
    ('6a0b3015-c4ad-4f4e-9076-42577bfc1593', 1, 'CSS', NULL),
    ('6a0b3015-c4ad-4f4e-9076-42577bfc1593', 2, 'PHP', NULL),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', 0, 'React', NULL),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', 1, 'TypeScript', NULL),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', 2, 'Redux Toolkit', NULL),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', 3, 'MUI', NULL),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', 4, 'Vite', NULL),
    ('c25b9bf6-32d6-4425-8d12-24b531ef78c5', 5, 'GSAP', NULL),
    ('f4658d87-99ec-455a-96ea-e6452e8bff16', 0, 'React', NULL),
    ('f4658d87-99ec-455a-96ea-e6452e8bff16', 1, 'TypeScript', NULL),
    ('f4658d87-99ec-455a-96ea-e6452e8bff16', 2, 'Jest', NULL),
    ('f4658d87-99ec-455a-96ea-e6452e8bff16', 3, 'Storybook', NULL),
    ('f4658d87-99ec-455a-96ea-e6452e8bff16', 4, 'GitHub Actions', NULL),
    ('9dccf1a9-c4ae-4bc5-9261-06cdcdee4678', 0, 'JavaScript', NULL),
    ('9dccf1a9-c4ae-4bc5-9261-06cdcdee4678', 1, 'HTML', NULL),
    ('9dccf1a9-c4ae-4bc5-9261-06cdcdee4678', 2, 'SCSS', NULL),
    ('9dccf1a9-c4ae-4bc5-9261-06cdcdee4678', 3, 'Node.js', NULL),
    ('eb5879d6-45df-4090-a882-d2a4daf51ec5', 0, 'JavaScript', NULL),
    ('eb5879d6-45df-4090-a882-d2a4daf51ec5', 1, 'jQuery', NULL),
    ('eb5879d6-45df-4090-a882-d2a4daf51ec5', 2, 'SCSS', NULL),
    ('4d846059-40e8-47b1-b477-f6bb43852de1', 0, 'HTML', NULL),
    ('4d846059-40e8-47b1-b477-f6bb43852de1', 1, 'CSS', NULL),
    ('4d846059-40e8-47b1-b477-f6bb43852de1', 2, 'PHP', NULL);

INSERT INTO educations (id, person_id, position, location, university, grade, program_name, profile_name, faculty, course, end_date) VALUES
    ('9ded49c6-6eeb-4716-b4d2-9ccb2e702627', '550e8400-e29b-41d4-a716-446655440001', 0, 'Новосибирск', 'НГТУ', 3, 'Программная инженерия', 'Разработка веб-приложений', 'Автоматики и вычислительной техники', 2, NULL),
    ('22f85f20-6927-4526-a149-75cc89fd94f7', '550e8400-e29b-41d4-a716-446655440001', 1, 'Москва', 'МГТУ им. Баумана', 1, 'Информатика и вычислительная техника', 'Разработка ПО', 'Информатика и системы управления', NULL, '2021-01-01'),
    ('79fbe989-42fb-413f-ac6f-a7b0f07ca085', '550e8400-e29b-41d4-a716-446655440002', 0, 'Новосибирск', 'НГТУ', 3, 'Программная инженерия', 'Разработка веб-приложений', 'Автоматики и вычислительной техники', 2, NULL),
    ('371ac0f1-8e52-411a-b0ce-7909089e06f0', '550e8400-e29b-41d4-a716-446655440002', 1, 'Москва', 'МГТУ им. Баумана', 1, 'Информатика и вычислительная техника', 'Разработка ПО', 'Информатика и системы управления', NULL, '2021-01-01');

INSERT INTO courses (id, person_id, position, name, organization, end_date, hours, certificate_link) VALUES
    ('808fec20-46d8-48ee-9108-6620df06bf58', '550e8400-e29b-41d4-a716-446655440001', 0, 'Advanced React и производительность', 'Frontend Masters', '2024-05-01', 40, 'https://example.com/certificate/react'),
    ('abaf0c57-0dce-42e7-b951-5509f8226c97', '550e8400-e29b-41d4-a716-446655440001', 1, 'Архитектура фронтенд-приложений', 'Яндекс Практикум', '2023-11-01', 72, 'https://example.com/certificate/architecture'),
    ('1fbabd5d-6313-4a0b-9c13-c6fb3f8c80f0', '550e8400-e29b-41d4-a716-446655440001', 2, 'Node.js: разработка серверных приложений', 'Stepik', '2022-03-01', 36, NULL);

INSERT INTO activities (id, person_id, position, type, title, description, date, link) VALUES
    ('f8038982-9446-4dc2-8462-013ee03ae4ac', '550e8400-e29b-41d4-a716-446655440001', 0, 'openSource', 'Контрибьютор в open-source', 'Исправления и улучшения в популярных UI-библиотеках, 15+ принятых pull request', '2024-01-01', 'https://github.com/example'),
    ('832d94f5-a043-49b6-9e92-01085684be5a', '550e8400-e29b-41d4-a716-446655440001', 1, 'hackathon', 'Хакатон «Цифровой прорыв»', 'Финалист в команде из 4 человек: прототип сервиса за 48 часов', '2023-09-01', NULL),
    ('f3045fc0-9669-482f-8ee6-65d872ba2317', '550e8400-e29b-41d4-a716-446655440001', 2, 'speaker', 'Доклад «Анимации на GSAP без боли»', 'Выступление на митапе фронтенд-разработчиков Новосибирска', '2024-03-01', 'https://example.com/talk'),
    ('9f15c315-291b-4e3b-854b-a1af3f8b3c5c', '550e8400-e29b-41d4-a716-446655440001', 3, 'mentoring', 'Менторство джунов', 'Помогаю начинающим разработчикам с ревью кода и карьерным планом', NULL, NULL),
    ('2fbdd476-9e5f-4f3d-bc41-0c0f5239ffbc', '550e8400-e29b-41d4-a716-446655440001', 4, 'community', 'Организатор митапов', 'Помогаю проводить встречи локального IT-сообщества', '2022-10-01', NULL);

INSERT INTO projects (id, person_id, position, name, description, link) VALUES
    ('2804a21a-b0ba-48d8-9458-4056e6af0572', '550e8400-e29b-41d4-a716-446655440001', 0, 'CV Platform', 'Платформа для отображения резюме разработчиков', 'https://example.com'),
    ('595bfb4c-fc46-4fb4-9f74-9bb2b3c772f3', '550e8400-e29b-41d4-a716-446655440001', 1, 'Finance Tracker', 'Приложение для учёта личных финансов с аналитикой по категориям', 'https://example.com/finance'),
    ('93579362-3174-49cc-8cb9-184866c8ee94', '550e8400-e29b-41d4-a716-446655440001', 2, 'Pulse Dashboard', 'Дашборд метрик в реальном времени: графики, алерты и фильтры', 'https://example.com/pulse'),
    ('02b252f5-b516-4acf-91a1-b1ff87447c03', '550e8400-e29b-41d4-a716-446655440001', 3, 'Design Kit', 'Библиотека UI-компонентов с темами и документацией в Storybook', NULL),
    ('63605cee-7200-48bc-9e08-7d50a7a7d2e1', '550e8400-e29b-41d4-a716-446655440002', 0, 'CV Platform', 'Платформа для отображения резюме разработчиков', 'https://example.com'),
    ('b32fe05d-9ce3-4066-afec-2ee03c213865', '550e8400-e29b-41d4-a716-446655440002', 1, 'Finance Tracker', 'Приложение для учёта личных финансов с аналитикой по категориям', 'https://example.com/finance'),
    ('59514890-633b-498f-b508-16525b0bd8a4', '550e8400-e29b-41d4-a716-446655440002', 2, 'Pulse Dashboard', 'Дашборд метрик в реальном времени: графики, алерты и фильтры', 'https://example.com/pulse'),
    ('cc8b31dc-290f-48ea-8be9-0c801563acd3', '550e8400-e29b-41d4-a716-446655440002', 3, 'Design Kit', 'Библиотека UI-компонентов с темами и документацией в Storybook', NULL);

INSERT INTO project_technologies (project_id, position, name, grade) VALUES
    ('2804a21a-b0ba-48d8-9458-4056e6af0572', 0, 'React', NULL),
    ('2804a21a-b0ba-48d8-9458-4056e6af0572', 1, 'TypeScript', NULL),
    ('2804a21a-b0ba-48d8-9458-4056e6af0572', 2, 'Vite', NULL),
    ('2804a21a-b0ba-48d8-9458-4056e6af0572', 3, 'GSAP', NULL),
    ('595bfb4c-fc46-4fb4-9f74-9bb2b3c772f3', 0, 'React', NULL),
    ('595bfb4c-fc46-4fb4-9f74-9bb2b3c772f3', 1, 'Redux', NULL),
    ('595bfb4c-fc46-4fb4-9f74-9bb2b3c772f3', 2, 'Node.js', NULL),
    ('595bfb4c-fc46-4fb4-9f74-9bb2b3c772f3', 3, 'PostgreSQL', NULL),
    ('93579362-3174-49cc-8cb9-184866c8ee94', 0, 'React', NULL),
    ('93579362-3174-49cc-8cb9-184866c8ee94', 1, 'TypeScript', NULL),
    ('93579362-3174-49cc-8cb9-184866c8ee94', 2, 'WebSocket', NULL),
    ('93579362-3174-49cc-8cb9-184866c8ee94', 3, 'D3', NULL),
    ('02b252f5-b516-4acf-91a1-b1ff87447c03', 0, 'React', NULL),
    ('02b252f5-b516-4acf-91a1-b1ff87447c03', 1, 'MUI', NULL),
    ('02b252f5-b516-4acf-91a1-b1ff87447c03', 2, 'Storybook', NULL),
    ('63605cee-7200-48bc-9e08-7d50a7a7d2e1', 0, 'React', NULL),
    ('63605cee-7200-48bc-9e08-7d50a7a7d2e1', 1, 'TypeScript', NULL),
    ('63605cee-7200-48bc-9e08-7d50a7a7d2e1', 2, 'Vite', NULL),
    ('63605cee-7200-48bc-9e08-7d50a7a7d2e1', 3, 'GSAP', NULL),
    ('b32fe05d-9ce3-4066-afec-2ee03c213865', 0, 'React', NULL),
    ('b32fe05d-9ce3-4066-afec-2ee03c213865', 1, 'Redux', NULL),
    ('b32fe05d-9ce3-4066-afec-2ee03c213865', 2, 'Node.js', NULL),
    ('b32fe05d-9ce3-4066-afec-2ee03c213865', 3, 'PostgreSQL', NULL),
    ('59514890-633b-498f-b508-16525b0bd8a4', 0, 'React', NULL),
    ('59514890-633b-498f-b508-16525b0bd8a4', 1, 'TypeScript', NULL),
    ('59514890-633b-498f-b508-16525b0bd8a4', 2, 'WebSocket', NULL),
    ('59514890-633b-498f-b508-16525b0bd8a4', 3, 'D3', NULL),
    ('cc8b31dc-290f-48ea-8be9-0c801563acd3', 0, 'React', NULL),
    ('cc8b31dc-290f-48ea-8be9-0c801563acd3', 1, 'MUI', NULL),
    ('cc8b31dc-290f-48ea-8be9-0c801563acd3', 2, 'Storybook', NULL);