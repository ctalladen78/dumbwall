create table if not exists actions (
    id bigserial primary key,
    user_id bigserial not null references users (id) on delete cascade,
    post_id bigserial not null references posts (id) on delete cascade,
    action_type numeric default 0,
    created_at timestamptz default now(),
    constraint unique_action_constraint unique (user_id, post_id)
);