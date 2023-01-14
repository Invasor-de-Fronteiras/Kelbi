
DROP TABLE IF EXISTS public.raviregister;
DROP TABLE IF EXISTS public.ravistate;
DROP TABLE IF EXISTS public.ravisupport;

DROP TABLE IF EXISTS public.user_binaries;

CREATE TABLE IF NOT EXISTS public.user_binary
(
    id serial NOT NULL PRIMARY KEY,
    type2 bytea,
    type3 bytea,
    house_tier bytea,
    house_state int,
    house_password text,
    house_data bytea,
    house_furniture bytea,
    bookshelf bytea,
    gallery bytea,
    tore bytea,
    garden bytea,
    mission bytea
);

INSERT INTO public.user_binary (id) SELECT c.id FROM characters c;

UPDATE public.user_binary
    SET house_furniture = (SELECT house FROM characters WHERE user_binary.id = characters.id);

UPDATE public.user_binary
    SET mission = (SELECT trophy FROM characters WHERE user_binary.id = characters.id);

ALTER TABLE IF EXISTS public.characters
    DROP COLUMN IF EXISTS house;

ALTER TABLE IF EXISTS public.characters
    DROP COLUMN IF EXISTS trophy;

CREATE TABLE IF NOT EXISTS public.servers
(
    server_id integer NOT NULL,
    season integer NOT NULL,
    current_players integer NOT NULL,
    world_name text COLLATE pg_catalog."default",
    world_description text,
    land integer
);

ALTER TABLE public.servers
    ADD COLUMN IF NOT EXISTS world_name text COLLATE pg_catalog."default";

ALTER TABLE public.servers
    ADD COLUMN IF NOT EXISTS world_description text;

ALTER TABLE public.servers
    ADD COLUMN IF NOT EXISTS land integer;

ALTER TABLE IF EXISTS public.guilds
    ADD COLUMN IF NOT EXISTS pugi_outfit_1 int NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.guilds
    ADD COLUMN IF NOT EXISTS pugi_outfit_2 int NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.guilds
    ADD COLUMN IF NOT EXISTS pugi_outfit_3 int NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.guilds
    ADD COLUMN IF NOT EXISTS pugi_outfits int NOT NULL DEFAULT 0;


CREATE TABLE IF NOT EXISTS public.feature_weapon
(
    start_time timestamp without time zone NOT NULL,
    featured integer NOT NULL
);

ALTER TABLE IF EXISTS public.characters
    ADD COLUMN IF NOT EXISTS cafe_reset timestamp without time zone;

-- Bundled schema
-- Bundled netcafe rewards
INSERT INTO public.cafebonus (time_req, item_type, item_id, quantity)
VALUES
    (1800, 17, 0, 250),
    (3600, 17, 0, 500),
    (7200, 17, 0, 1000),
    (10800, 17, 0, 1500),
    (18000, 17, 0, 1750),
    (28800, 17, 0, 3000),
    (43200, 17, 0, 4000);

-- Bundled festa prizes
INSERT INTO public.festa_prizes
    (type, tier, souls_req, item_id, num_item)
VALUES
    ('personal', 1, 1, 9647, 7),
    ('personal', 2, 1, 9647, 7),
    ('personal', 3, 1, 9647, 7),
    ('personal', 1, 200, 11284, 4),
    ('personal', 2, 200, 11284, 4),
    ('personal', 3, 200, 11284, 4),
    ('personal', 1, 400, 11381, 3),
    ('personal', 2, 400, 11381, 3),
    ('personal', 3, 400, 11381, 3),
    ('personal', 1, 600, 11284, 8),
    ('personal', 2, 600, 11284, 8),
    ('personal', 3, 600, 11284, 8),
    ('personal', 1, 800, 11384, 3),
    ('personal', 2, 800, 11384, 3),
    ('personal', 3, 800, 11384, 3),
    ('personal', 1, 1000, 11284, 12),
    ('personal', 2, 1000, 11284, 12),
    ('personal', 3, 1000, 11284, 12),
    ('personal', 1, 1200, 11381, 5),
    ('personal', 2, 1200, 11381, 5),
    ('personal', 3, 1200, 11381, 5),
    ('personal', 1, 1400, 11284, 16),
    ('personal', 2, 1400, 11284, 16),
    ('personal', 3, 1400, 11284, 16),
    ('personal', 1, 1700, 11384, 5),
    ('personal', 2, 1700, 11384, 5),
    ('personal', 3, 1700, 11384, 5),
    ('personal', 1, 2000, 11284, 16),
    ('personal', 2, 2000, 11284, 16),
    ('personal', 3, 2000, 11284, 16),
    ('personal', 1, 2500, 11382, 4),
    ('personal', 2, 2500, 11382, 4),
    ('personal', 3, 2500, 11382, 4),
    ('personal', 1, 3000, 11284, 24),
    ('personal', 2, 3000, 11284, 24),
    ('personal', 3, 3000, 11284, 24),
    ('personal', 1, 4000, 11385, 4),
    ('personal', 2, 4000, 11385, 4),
    ('personal', 3, 4000, 11385, 4),
    ('personal', 1, 5000, 11381, 11),
    ('personal', 2, 5000, 11381, 11),
    ('personal', 3, 5000, 11381, 11),
    ('personal', 1, 6000, 5177, 5),
    ('personal', 2, 6000, 5177, 5),
    ('personal', 3, 6000, 5177, 5),
    ('personal', 1, 7000, 11384, 11),
    ('personal', 2, 7000, 11384, 11),
    ('personal', 3, 7000, 11384, 11),
    ('personal', 1, 10000, 11382, 8),
    ('personal', 2, 10000, 11382, 8),
    ('personal', 3, 10000, 11382, 8),
    ('personal', 1, 15000, 11385, 4),
    ('personal', 2, 15000, 11385, 4),
    ('personal', 3, 15000, 11385, 4),
    ('personal', 1, 20000, 11381, 13),
    ('personal', 2, 20000, 11381, 13),
    ('personal', 3, 20000, 11381, 13),
    ('personal', 1, 25000, 11385, 4),
    ('personal', 2, 25000, 11385, 4),
    ('personal', 3, 25000, 11385, 4),
    ('personal', 1, 30000, 11383, 1),
    ('personal', 2, 30000, 11383, 1),
    ('personal', 3, 30000, 11383, 1);

INSERT INTO public.festa_prizes
(type, tier, souls_req, item_id, num_item)
VALUES
    ('guild', 1, 100, 7468, 5),
    ('guild', 2, 100, 7468, 5),
    ('guild', 3, 100, 7465, 5),
    ('guild', 1, 300, 7469, 5),
    ('guild', 2, 300, 7469, 5),
    ('guild', 3, 300, 7466, 5),
    ('guild', 1, 700, 7470, 5),
    ('guild', 2, 700, 7470, 5),
    ('guild', 3, 700, 7467, 5),
    ('guild', 1, 1500, 13405, 14),
    ('guild', 1, 1500, 1520, 3),
    ('guild', 2, 1500, 13405, 14),
    ('guild', 2, 1500, 1520, 3),
    ('guild', 3, 1500, 7011, 3),
    ('guild', 3, 1500, 13405, 14),
    ('guild', 1, 3000, 10201, 10),
    ('guild', 2, 3000, 10201, 10),
    ('guild', 3, 3000, 10201, 10),
    ('guild', 1, 6000, 13895, 14),
    ('guild', 1, 6000, 1520, 6),
    ('guild', 2, 6000, 13895, 14),
    ('guild', 2, 6000, 1520, 6),
    ('guild', 3, 6000, 13895, 14),
    ('guild', 3, 6000, 7011, 4),
    ('guild', 1, 12000, 13406, 14),
    ('guild', 1, 12000, 1520, 9),
    ('guild', 2, 12000, 13406, 14),
    ('guild', 2, 12000, 1520, 9),
    ('guild', 3, 12000, 13406, 14),
    ('guild', 3, 12000, 7011, 5),
    ('guild', 1, 25000, 10207, 10),
    ('guild', 2, 25000, 10207, 10),
    ('guild', 3, 25000, 10207, 10),
    ('guild', 1, 50000, 1520, 12),
    ('guild', 1, 50000, 13896, 14),
    ('guild', 2, 50000, 1520, 12),
    ('guild', 2, 50000, 13896, 14),
    ('guild', 3, 50000, 7011, 6),
    ('guild', 3, 50000, 13896, 14),
    ('guild', 1, 100000, 10201, 10),
    ('guild', 2, 100000, 10201, 10),
    ('guild', 3, 100000, 10201, 10),
    ('guild', 1, 200000, 13406, 16),
    ('guild', 2, 200000, 13406, 16),
    ('guild', 3, 200000, 13406, 16),
    ('guild', 1, 300000, 13896, 16),
    ('guild', 2, 300000, 13896, 16),
    ('guild', 3, 300000, 13896, 16),
    ('guild', 1, 400000, 10207, 10),
    ('guild', 2, 400000, 10207, 10),
    ('guild', 3, 400000, 10207, 10),
    ('guild', 1, 500000, 13407, 6),
    ('guild', 1, 500000, 13897, 6),
    ('guild', 2, 500000, 13407, 6),
    ('guild', 2, 500000, 13897, 6),
    ('guild', 3, 500000, 13407, 6),
    ('guild', 3, 500000, 13897, 6);

-- Bundled festa trials
INSERT INTO public.festa_trials
    (objective, goal_id, times_req, locale_req, reward)
VALUES
    (1,27,1,0,1),
    (5,53034,0,0,400),
    (5,22042,0,0,89),
    (5,23397,0,0,89),
    (1,28,1,0,1),
    (1,68,1,0,1),
    (1,6,1,0,2),
    (1,38,1,0,2),
    (1,20,1,0,3),
    (1,39,1,0,4),
    (1,48,1,0,4),
    (1,67,1,0,4),
    (1,93,1,0,4),
    (1,22,1,0,5),
    (1,52,1,0,5),
    (1,101,1,0,5),
    (1,1,1,0,5),
    (1,37,1,0,5),
    (1,15,1,0,5),
    (1,45,1,0,5),
    (1,74,1,0,5),
    (1,78,1,0,5),
    (1,103,1,0,5),
    (1,51,1,0,6),
    (1,17,1,0,6),
    (1,21,1,0,6),
    (1,92,1,0,6),
    (1,47,1,0,7),
    (1,46,1,0,7),
    (1,26,1,0,7),
    (1,14,1,0,7),
    (1,11,1,0,7),
    (1,44,1,0,8),
    (1,43,1,0,8),
    (1,49,1,0,8),
    (1,40,1,0,8),
    (1,76,1,0,8),
    (1,89,1,0,8),
    (1,94,1,0,8),
    (1,96,1,0,8),
    (1,75,1,0,8),
    (1,91,1,0,8),
    (1,53,1,0,9),
    (1,80,1,0,9),
    (1,42,1,0,9),
    (1,79,1,0,9),
    (1,81,1,0,10),
    (1,41,1,0,10),
    (1,82,1,0,10),
    (1,90,1,0,10),
    (1,149,1,0,10),
    (1,85,1,0,11),
    (1,95,1,0,11),
    (1,121,1,0,11),
    (1,142,1,0,11),
    (1,141,1,0,11),
    (1,146,1,0,12),
    (1,147,1,0,12),
    (1,148,1,0,12),
    (1,151,1,0,12),
    (1,152,1,0,12),
    (1,159,1,0,12),
    (1,153,1,0,12),
    (1,162,1,0,12),
    (1,111,1,0,13),
    (1,110,1,0,13),
    (1,112,1,0,13),
    (1,109,1,0,14),
    (1,169,1,0,15),
    (2,33,1,0,6),
    (2,104,1,0,8),
    (2,119,1,0,8),
    (2,120,1,0,8),
    (2,54,1,0,8),
    (2,59,1,0,8),
    (2,64,1,0,8),
    (2,65,1,0,8),
    (2,99,1,0,9),
    (2,83,1,0,9),
    (2,84,1,0,10),
    (2,77,1,0,10),
    (2,106,1,0,10),
    (2,55,1,0,10),
    (2,58,1,0,10),
    (2,7,1,0,10),
    (2,50,1,0,11),
    (2,131,1,0,11),
    (2,129,1,0,11),
    (2,140,1,0,11),
    (2,122,1,0,11),
    (2,126,1,0,11),
    (2,127,1,0,11),
    (2,128,1,0,11),
    (2,130,1,0,11),
    (2,139,1,0,11),
    (2,144,1,0,11),
    (2,150,1,0,11),
    (2,158,1,0,11),
    (2,164,1,0,15),
    (2,165,1,0,15),
    (2,2,1,7,15),
    (2,36,1,0,15),
    (2,71,1,0,15),
    (2,108,1,0,15),
    (2,116,1,0,15),
    (2,107,1,0,15),
    (2,154,1,0,17),
    (2,166,1,0,17),
    (2,170,1,0,18),
    (3,31,1,0,1),
    (3,8,1,0,3),
    (3,123,1,0,8),
    (3,105,1,0,9),
    (3,125,1,0,11),
    (3,115,1,0,12),
    (3,114,1,0,12),
    (3,161,1,0,12),
    (4,670,1,0,1),
    (4,671,1,0,1),
    (4,672,1,0,1),
    (4,675,1,0,1),
    (4,673,1,0,1),
    (4,674,1,0,1);
