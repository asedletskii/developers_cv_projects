CREATE TABLE resume_experience_tasks (
    experience_id uuid NOT NULL REFERENCES resume_experiences(id) ON DELETE CASCADE,
    position integer NOT NULL,
    value text NOT NULL,
    PRIMARY KEY (experience_id, position),
    CONSTRAINT resume_experience_tasks_position_check CHECK (position >= 0),
    CONSTRAINT resume_experience_tasks_value_check CHECK (length(btrim(value)) > 0)
);
