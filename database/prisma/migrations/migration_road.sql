ALTER TABLE IF EXISTS public.normal_shop_items
    ADD COLUMN enable_weeks varchar(8);
	
ALTER TABLE IF EXISTS public.shop_item_state
    ADD COLUMN week integer;