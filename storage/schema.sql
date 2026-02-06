create table if not exists applications (
    id text primary key,
    created_at timestamptz not null default now(),
    strategy text not null,
    status text not null,
    payload jsonb not null
);

create table if not exists check_results (
    id bigserial primary key,
    applications_id text not null references applications(id),
    created_at timestamptz not null default now(),
    check_name text not null,
    status text not null,
    reason text not null default ''
);