create table if not exists user_profiles (
    id uuid not null references auth.users on delete cascade,
    first_name text,
    last_name text,

    primary key (id)
);