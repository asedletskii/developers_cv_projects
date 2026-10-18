# Публичная проекция резюме

Публичная страница открывается по **идентификатору именованной ссылки** (`linkId`), не по `resumeId`. Единственная публичная ручка получения резюме — `GET /api/v1/cv/:linkId`.

Актуальная функциональная спецификация: [FR-PUB](../../docs_new/03-functional/FR-PUB.md). Контракт посещений и аналитики: [analytics API](./analytics.md).

## `GET /api/v1/cv/:linkId`

Без тела запроса и аутентификации. Бэкенд находит ссылку, проверяет её статус, разрешает связанное резюме и возвращает его только если резюме опубликовано (`opened`).

Успешное тело — публичная проекция модуля `resume`, а не кабинетное представление:

```json
{
  "title": "Backend Developer",
  "language": "en",
  "createdAt": "2026-10-01T10:00:00Z",
  "updatedAt": "2026-10-10T12:00:00Z",
  "overview": {
    "greeting": {
      "first": "Hello",
      "second": "I'm Ada",
      "secondHighlight": "Ada",
      "third": "Backend Developer"
    },
    "sections": []
  },
  "profile": {
    "id": "profile-id",
    "name": "Ada",
    "surname": "Lovelace",
    "position": "Backend Developer",
    "photoPath": "https://example.com/photo.png",
    "dateOfBirth": "",
    "location": "",
    "relocation": false,
    "workFormat": []
  },
  "contacts": {
    "id": "contacts-id",
    "hiddenContacts": ["email"]
  },
  "skills": {
    "id": "skills-id",
    "hardSkills": ["Go"],
    "softSkills": []
  },
  "languages": [],
  "technologies": [],
  "experiences": [],
  "educations": [],
  "courses": [],
  "activities": [],
  "projects": []
}
```

The example is illustrative: optional values are omitted, IDs are child-section IDs, and the collections are empty in this example. The response does not expose a top-level resume ID, owner/user ID, resume status, or version. It includes profile, contacts, skills and overview as nested objects; collection names are `experiences` and `educations`. A hidden contact value is omitted from the JSON and listed in `contacts.hiddenContacts`.
The public projection does not include `overview.resumeLink`. The editor's overview contract is unchanged by this public projection change.

The public projection is built by `resume/application/queries.PublicResumeQueryHandler` and its reader. It intentionally differs from the full editor view. Public projects omit editor-only fields (role, Git link and contributor data); section fields and supported collection entries are described by the corresponding `Public*View` types.

### Response headers

| Header | Behavior |
|---|---|
| `Cache-Control` | `no-store` |
| `X-View-Id` | Present only when optional visit creation succeeds |
| `Set-Cookie` | A visitor cookie is set when required by successful visit creation |

Failure to create a visit does not fail the resume response. Without `X-View-Id`, the frontend does not start the visit tracker.

### Errors

| Status | Current behavior |
|---|---|
| `400` | `linkId` is not a UUID |
| `404` | Link does not exist or the associated resume is not published |
| `410` | Link exists but is inactive |

The `410` behavior for an inactive link is intentional and retained. Unknown links and unpublished resumes return `404`.

## `POST /api/v1/cv/:linkId/contact-views`

Reveals one hidden contact or returns available contacts for download. This route also uses `linkId`; it never accepts a public resume ID as an access key.

For one contact:

```json
{
  "initiator": "contact",
  "contact": "email",
  "viewId": "8c0e…",
  "turnstileToken": "..."
}
```

Response: `{"contact":"email","value":"person@example.com"}`.

For download:

```json
{
  "initiator": "download",
  "viewId": "8c0e…",
  "turnstileToken": "..."
}
```

Response: `{"contacts":[{"contact":"email","value":"person@example.com"}]}`. Only non-empty contacts are included. An omitted `initiator` is accepted for compatibility when `contact` is supplied.

Both paths use the same application resolver and status mapping as public resume retrieval: unknown links return `404`, inactive links return `410`. The resolver checks the link before asking `resume` for contact data. If a verifier is configured, a valid one-time Turnstile token is required. A `viewId` is optional; when supplied, it is used only for best-effort analytics recording and is checked against the link by the analytics recorder.

The `resume` module owns the contact data and hidden-contact rules. The link flow owns resolving the public link. Analytics records action metadata only; it must not store contact values.
