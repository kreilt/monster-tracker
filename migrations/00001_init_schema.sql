-- +goose Up

CREATE TYPE rarity AS ENUM ('ordinary','unusual','rare','very rare','legendary');

CREATE TYPE statuses AS ENUM ('disc','lim','new','reg','alc','current');

CREATE TABLE flavors
(flavor_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
title text NOT NULL,
lineup text NOT NULL,
description text,
rare rarity NOT NULL,
region varchar(50) NOT NULL,
color varchar(10) NOT NULL,
status statuses NOT NULL,
photo text,
UNIQUE(title)
);

CREATE TABLE users
(user_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
external_id text NOT NULL UNIQUE,
nickname text,
email text NOT NULL
);

CREATE TABLE flavors_users
(flavor_user_id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
user_id integer NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
flavor_id integer NOT NULL REFERENCES flavors(flavor_id),
status_tried boolean NOT NULL,
tried_at timestamptz,
user_photo text,
UNIQUE (user_id, flavor_id)
);

-- +goose Down

DROP TABLE flavors_users;
DROP TABLE users;
DROP TABLE flavors;
DROP TYPE statuses;
DROP TYPE rarity;