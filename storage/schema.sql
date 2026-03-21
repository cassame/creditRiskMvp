CREATE TABLE IF NOT EXISTS strategies (
    code TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
INSERT INTO strategies (code) VALUES
    ('resident_first_time', 'Resident: First time'),
    ('resident_repeat', 'Resident: Repeat'),
    ('nonresident_first_time', 'Нерезидент: Первый раз'),
    ('nonresident_repeat', 'Нерезидент: Повторно');
create table if not exists application_statuses (
    code text primary key,
    description text not null,
    created_at TIMESTAMPTZ not null default now()
);
INSERT INTO application_statuses (code, description) VALUES
    ('approved', 'Application was approved'),
    ('rejected', 'Application was rejected'),
    ('manual_review', 'Application requires manual review');

CREATE TABLE IF NOT EXISTS check_statuses (
    code TEXT PRIMARY KEY,
    description text not null,
    created_at TIMESTAMPTZ not null default now()
);
INSERT INTO check_statuses (code, description) VALUES
    ('passed', 'Check passed successfully'),
    ('failed', 'Check failed'),
    ('error', 'External service error');
CREATE TABLE IF NOT EXISTS check_names (
    code TEXT PRIMARY KEY,
    description TEXT
);
INSERT INTO check_names (code) VALUES
    ('age>=18'), ('valid_phone'), ('valid_passport'), ('has_patronymic'),
    ('approve_amount'), ('credit_history'), ('terrorist'), ('bankruptcy');

CREATE TABLE IF NOT EXISTS applications (
    id UUID primary key,
    created_at TIMESTAMPTZ not null default now(),
    strategy TEXT NOT NULL REFERENCES strategies(code),
    status text NOT NULL REFERENCES application_statuses(code),
    payload jsonb not null
);
CREATE TABLE IF NOT EXISTS check_results (
    id BIGSERIAL PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES applications(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    check_name TEXT NOT NULL REFERENCES check_names(code),
    status TEXT NOT NULL REFERENCES check_statuses(code),
    reason TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_check_results_application_id
    ON check_results(application_id);
CREATE INDEX  idx_applications_status
    ON applications(status);


