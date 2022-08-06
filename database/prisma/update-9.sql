BEGIN;

ALTER TABLE IF EXISTS public.users
    ADD COLUMN IF NOT EXISTS last_character int DEFAULT 0;


CREATE SEQUENCE IF NOT EXISTS public.airou_id_seq;

ALTER TABLE IF EXISTS public.gook
    DROP COLUMN IF EXISTS gook0status;

ALTER TABLE IF EXISTS public.gook
    DROP COLUMN IF EXISTS gook1status;

ALTER TABLE IF EXISTS public.gook
    DROP COLUMN IF EXISTS gook2status;

ALTER TABLE IF EXISTS public.gook
    DROP COLUMN IF EXISTS gook3status;

ALTER TABLE IF EXISTS public.gook
    DROP COLUMN IF EXISTS gook4status;

ALTER TABLE IF EXISTS public.gook
    DROP COLUMN IF EXISTS gook5status;

ALTER TABLE IF EXISTS public.gook
    DROP COLUMN IF EXISTS gook5;

UPDATE public.gook SET gook0=NULL, gook1=NULL, gook2=NULL, gook3=NULL, gook4=NULL;

ALTER TABLE IF EXISTS public.guilds
    ADD COLUMN IF NOT EXISTS pugi_name_1 varchar(12) DEFAULT '';

ALTER TABLE IF EXISTS public.guilds
    ADD COLUMN IF NOT EXISTS pugi_name_2 varchar(12) DEFAULT '';

ALTER TABLE IF EXISTS public.guilds
    ADD COLUMN IF NOT EXISTS pugi_name_3 varchar(12) DEFAULT '';

CREATE TABLE IF NOT EXISTS public.guild_adventures
(
    id serial NOT NULL PRIMARY KEY,
    guild_id int NOT NULL,
    destination int NOT NULL,
    charge int NOT NULL DEFAULT 0,
    depart int NOT NULL,
    return int NOT NULL,
    collected_by text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS public.guild_meals
(
    id serial NOT NULL PRIMARY KEY,
    guild_id int NOT NULL,
    meal_id int NOT NULL,
    level int NOT NULL,
    expires int NOT NULL
);

CREATE TABLE IF NOT EXISTS public.guild_hunts
(
    id serial NOT NULL PRIMARY KEY,
    guild_id int NOT NULL,
    host_id int NOT NULL,
    destination int NOT NULL,
    level int NOT NULL,
    return int NOT NULL,
    acquired bool NOT NULL DEFAULT false,
    claimed bool NOT NULL DEFAULT false,
    hunters text NOT NULL DEFAULT '',
    treasure text NOT NULL DEFAULT '',
    hunt_data bytea NOT NULL,
    cats_used text NOT NULL
);

ALTER TABLE IF EXISTS public.characters
    ADD COLUMN IF NOT EXISTS house bytea;

ALTER TABLE IF EXISTS public.characters
    ADD COLUMN IF NOT EXISTS scenariodata bytea;

ALTER TABLE IF EXISTS public.characters
    ADD COLUMN IF NOT EXISTS savefavoritequest bytea;

ALTER TABLE IF EXISTS public.characters
    ADD COLUMN IF NOT EXISTS friends text NOT NULL DEFAULT '';

ALTER TABLE IF EXISTS public.characters
    ADD COLUMN IF NOT EXISTS blocked text NOT NULL DEFAULT '';

ALTER TABLE IF EXISTS public.characters
    ADD COLUMN IF NOT EXISTS deleted boolean NOT NULL DEFAULT false;

DROP TABLE IF EXISTS public.sign_sessions;
CREATE TABLE IF NOT EXISTS public.sign_sessions
(
    user_id int NOT NULL,
    char_id int,
    token varchar(16) NOT NULL,
    server_id integer
);

DROP TABLE IF EXISTS public.servers;
CREATE TABLE IF NOT EXISTS public.servers
(
    server_id int NOT NULL,
    season int NOT NULL,
    current_players int NOT NULL
);

CREATE TABLE IF NOT EXISTS public.user_binaries
(
    id int PRIMARY KEY,
    type2 bytea,
    type3 bytea
);

ALTER TABLE IF EXISTS public.guilds
    ADD COLUMN IF NOT EXISTS recruiting bool NOT NULL DEFAULT true;

END;