create table if not exists users (
        id bigserial primary key,
        login text unique not null,
        email text unique not null,
        email_verified boolean default false,
        password text not null,
        active boolean default true,
	karma numeric default 0
);

create table if not exists posts (
        id bigserial primary key,
        uid text unique not null,
        type integer not null,
        title text not null,
        body text not null,
        ups numeric not null default 0,
        downs numeric not null default 0,
        user_id bigserial not null references users (id) on delete cascade
);

create table if not exists comments (
        id bigserial primary key,
        body text not null,
        parent_id bigserial references comments (id) on delete cascade,
        post_id bigserial not null references posts (id) on delete cascade,
        user_id bigserial not null references users (id),
        ups numeric default 0,
        downs numeric default 0
);
