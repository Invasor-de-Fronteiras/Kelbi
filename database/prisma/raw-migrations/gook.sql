BEGIN;

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

END;