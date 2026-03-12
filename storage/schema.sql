create table if not exists application_statuses (
    code text primary key,
    description text not null,
    created_at timestamptz not null default now()
);
INSERT INTO application_statuses (code, description) VALUES
    ('approved', 'Application was approved'),
    ('rejected', 'Application was rejected'),
    ('manual_review', 'Application requires manual review');

create table if not exists check_statuses (
    code text primary key,
    description text not null,
    created_at timestamptz not null default now()
);
INSERT INTO check_statuses (code, description) VALUES
    ('passed', 'Check passed successfully'),
    ('failed', 'Check failed'),
    ('error', 'External service error');

create table if not exists applications (
    id text primary key,
    created_at timestamptz not null default now(),
    strategy text not null CHECK (
        strategy IN ('resident_first_time', 'resident_repeat',
        'nonresident_first_time', 'nonresident_repeat')
    ),
    status text NOT NULL REFERENCES application_statuses(code),
    payload jsonb not null
);

create table if not exists check_results (
    id bigserial primary key,
    application_id text not null references applications(id),
    created_at timestamptz not null default now(),
    check_name text not null CHECK (
        check_name IN ('age>=18', 'valid_phone', 'valid_passport',
        'has_patronymic', 'approve_amount', 'credit_history',
        'terrorist', 'bankruptcy')
    ),
    status text not null REFERENCES check_statuses(code),
    reason text not null default ''
);

CREATE INDEX idx_check_results_application_id
    ON check_results(application_id);
CREATE INDEX  idx_applications_status
    ON applications(status);


