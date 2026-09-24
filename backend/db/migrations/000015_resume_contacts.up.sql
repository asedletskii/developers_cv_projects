CREATE TABLE resume_contacts_v2 (
    id uuid PRIMARY KEY,
    resume_id uuid NOT NULL UNIQUE REFERENCES resume_roots(id) ON DELETE CASCADE,
    phone text,
    telegram text,
    email text,
    CONSTRAINT resume_contacts_v2_presence_check CHECK (phone IS NOT NULL OR telegram IS NOT NULL OR email IS NOT NULL)
);
CREATE INDEX idx_resume_contacts_v2_resume_id ON resume_contacts_v2(resume_id);

CREATE TABLE resume_contact_hidden (
    contacts_id uuid NOT NULL REFERENCES resume_contacts_v2(id) ON DELETE CASCADE,
    value text NOT NULL,
    position integer NOT NULL,
    PRIMARY KEY (contacts_id, value),
    UNIQUE (contacts_id, position),
    CONSTRAINT resume_contact_hidden_value_check CHECK (value IN ('telegram','email','phone')),
    CONSTRAINT resume_contact_hidden_position_check CHECK (position >= 0)
);
CREATE INDEX idx_resume_contact_hidden_position ON resume_contact_hidden(contacts_id, position);

CREATE TABLE resume_contacts_text_v2 (
    contacts_id uuid PRIMARY KEY REFERENCES resume_contacts_v2(id) ON DELETE CASCADE,
    title_ru text,
    title_en text,
    subtitle_ru text,
    subtitle_en text,
    CONSTRAINT resume_contacts_text_v2_presence_check CHECK (title_ru IS NOT NULL OR title_en IS NOT NULL OR subtitle_ru IS NOT NULL OR subtitle_en IS NOT NULL)
);

INSERT INTO resume_contacts_v2 (id, resume_id, phone, telegram, email)
SELECT gen_random_uuid(), r.id, c.phone, c.telegram, c.email
FROM resume_roots r
JOIN resume_contacts c ON c.resume_id = r.id
ON CONFLICT (resume_id) DO NOTHING;

INSERT INTO resume_contact_hidden (contacts_id, value, position)
SELECT c2.id, h.value, h.position
FROM resume_hidden_contacts h
JOIN resume_contacts_v2 c2 ON c2.resume_id = h.resume_id
ON CONFLICT (contacts_id, value) DO NOTHING;

INSERT INTO resume_contacts_text_v2 (contacts_id, title_ru, title_en, subtitle_ru, subtitle_en)
SELECT c2.id, t.title_ru, t.title_en, t.subtitle_ru, t.subtitle_en
FROM resume_contacts_text t
JOIN resume_contacts_v2 c2 ON c2.resume_id = t.resume_id
ON CONFLICT (contacts_id) DO NOTHING;
