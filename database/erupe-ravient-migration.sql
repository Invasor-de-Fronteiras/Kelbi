BEGIN;

CREATE TABLE IF NOT EXISTS public.raviregister
(
    id integer NOT NULL,
    refid integer NOT NULL,
    nextravi integer NOT NULL,
    ravistarted integer,
    raviposttime integer,
    ravitype integer,
    maxplayers integer,
    ravikilled integer,
    carvequest integer,
    register1 integer,
    register2 integer,
    register3 integer,
    register4 integer,	
    register5 integer,
    CONSTRAINT raviregister_pkey PRIMARY KEY (id)
)
 
TABLESPACE pg_default;
 
CREATE SEQUENCE IF NOT EXISTS public.raviregister_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1
    OWNED BY raviregister.id;
	
ALTER TABLE public.raviregister
ALTER id SET DEFAULT nextval('raviregister_id_seq'::regclass);

CREATE TABLE IF NOT EXISTS public.ravistate
(
    id integer NOT NULL,
    refid integer NOT NULL,
    phase1hp integer NOT NULL,
    phase2hp integer,
    phase3hp integer,
    phase4hp integer,
    phase5hp integer,
    phase6hp integer,
    phase7hp integer,
    phase8hp integer,
    phase9hp integer,
    unknown1 integer,
    unknown2 integer,
    unknown3 integer,
    unknown4 integer,
    unknown5 integer,
    unknown6 integer,
    unknown7 integer,
    unknown8 integer,
    unknown9 integer,
    unknown10 integer,
    unknown11 integer,
    unknown12 integer,
    unknown13 integer,
    unknown14 integer,
    unknown15 integer,
    unknown16 integer,
    unknown17 integer,
    unknown18 integer,
    unknown19 integer,
    unknown20 integer,
    damagemultiplier integer,
    CONSTRAINT ravistate_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

CREATE SEQUENCE IF NOT EXISTS public.ravistate_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1
    OWNED BY ravistate.id;

ALTER TABLE public.ravistate
ALTER id SET DEFAULT nextval('ravistate_id_seq'::regclass);

CREATE TABLE IF NOT EXISTS public.ravisupport
(
    id integer NOT NULL,
    refid integer NOT NULL,
    support1 integer NOT NULL,
    support2 integer,
    support3 integer,
    support4 integer,
    support5 integer,
    support6 integer,
    support7 integer,
    support8 integer,
    support9 integer,
    support10 integer,
    support11 integer,
    support12 integer,
    support13 integer,
    support14 integer,
    support15 integer,
    support16 integer,
    support17 integer,
    support18 integer,
    support19 integer,
    support20 integer,
    support21 integer,
    support22 integer,
    support23 integer,
    support24 integer,
    support25 integer,
    CONSTRAINT ravisupport_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

CREATE SEQUENCE IF NOT EXISTS public.ravisupport_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1
    OWNED BY ravisupport.id;

ALTER TABLE public.ravisupport
ALTER id SET DEFAULT nextval('ravisupport_id_seq'::regclass);

CREATE OR REPLACE PROCEDURE public.ravireset(
	IN zeroed integer)
LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
 
UPDATE public.ravistate
	SET refid=29, phase1hp=zeroed, phase2hp=zeroed, phase3hp=zeroed, phase4hp=zeroed, phase5hp=zeroed, phase6hp=zeroed, phase7hp=zeroed, phase8hp=zeroed, phase9hp=zeroed, unknown1=zeroed, unknown2=zeroed, unknown3=zeroed, unknown4=zeroed, unknown5=zeroed, unknown6=zeroed, unknown7=zeroed, unknown8=zeroed, unknown9=zeroed, unknown10=zeroed, unknown11=zeroed, unknown12=zeroed, unknown13=zeroed, unknown14=zeroed, unknown15=zeroed, unknown16=zeroed, unknown17=zeroed, unknown18=zeroed, unknown19=zeroed, unknown20=zeroed, damagemultiplier=1
	WHERE refid = 29;

UPDATE public.raviregister
	SET refid=12, nextravi=zeroed, ravistarted=zeroed, raviposttime=zeroed, ravitype=zeroed, maxplayers=zeroed, ravikilled=zeroed, carvequest=zeroed, register1=zeroed, register2=zeroed, register3=zeroed, register4=zeroed, register5=zeroed
	WHERE refid = 12;

UPDATE public.ravisupport
	SET refid=25, support1=zeroed, support2=zeroed, support3=zeroed, support4=zeroed, support5=zeroed, support6=zeroed, support7=zeroed, support8=zeroed, support9=zeroed, support10=zeroed, support11=zeroed, support12=zeroed, support13=zeroed, support14=zeroed, support15=zeroed, support16=zeroed, support17=zeroed, support18=zeroed, support19=zeroed, support20=zeroed, support21=zeroed, support22=zeroed, support23=zeroed, support24=zeroed, support25=zeroed
	WHERE refid = 25;

    COMMIT;
END;
$BODY$;

CREATE OR REPLACE PROCEDURE public.raviinit()
LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
 
INSERT INTO public.raviregister(
	refid, nextravi, ravistarted, raviposttime, ravitype, maxplayers, ravikilled, carvequest, register1, register2, register3, register4, register5)
	VALUES (12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);
	
INSERT INTO public.ravistate(
	refid, phase1hp, phase2hp, phase3hp, phase4hp, phase5hp, phase6hp, phase7hp, phase8hp, phase9hp, unknown1, unknown2, unknown3, unknown4, unknown5, unknown6, unknown7, unknown8, unknown9, unknown10, unknown11, unknown12, unknown13, unknown14, unknown15, unknown16, unknown17, unknown18, unknown19, unknown20, damagemultiplier)
	VALUES (29, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1);
	
INSERT INTO public.ravisupport(
	refid, support1, support2, support3, support4, support5, support6, support7, support8, support9, support10, support11, support12, support13, support14, support15, support16, support17, support18, support19, support20, support21, support22, support23, support24, support25)
	VALUES (25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);

    COMMIT;
END;
$BODY$;

ALTER TABLE IF EXISTS public.characters
    ADD COLUMN guild_post_checked timestamp without time zone NOT NULL DEFAULT now(),
    ADD COLUMN time_played int NOT NULL DEFAULT 0,
    ADD COLUMN weapon_id int NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.characters
    DROP COLUMN gr_override_mode,
    DROP COLUMN small_gr_level,
    DROP COLUMN gr_override_unk0,
    DROP COLUMN gr_override_unk1;

ALTER TABLE IF EXISTS public.characters
    RENAME COLUMN gr_override_level TO gr;

ALTER TABLE IF EXISTS public.characters
    RENAME COLUMN exp TO hrp;

ALTER TABLE IF EXISTS public.characters
    RENAME COLUMN weapon TO weapon_type;

ALTER TABLE IF EXISTS public.guilds
    ADD COLUMN item_box bytea,
    ADD COLUMN event_rp int NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS public.guilds
    RENAME COLUMN rp TO rank_rp;

ALTER TABLE IF EXISTS public.guilds
    DROP COLUMN guild_hall;
	
ALTER TABLE IF EXISTS public.users
    ADD COLUMN item_box bytea,
	ADD COLUMN rights integer NOT NULL DEFAULT 14;

CREATE TABLE IF NOT EXISTS public.guild_posts
(
    id serial NOT NULL PRIMARY KEY,
    guild_id int NOT NULL,
    author_id int NOT NULL,
    post_type int NOT NULL,
    stamp_id int NOT NULL,
    title text NOT NULL,
    body text NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    liked_by text NOT NULL DEFAULT ''
);

END;