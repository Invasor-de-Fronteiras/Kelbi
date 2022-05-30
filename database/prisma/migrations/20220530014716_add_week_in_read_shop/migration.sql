-- AlterTable
ALTER TABLE "mail" ADD COLUMN     "locked" BOOLEAN NOT NULL DEFAULT false;

-- AlterTable
ALTER TABLE "normal_shop_items" ADD COLUMN     "enable_weeks" VARCHAR(8);

-- AlterTable
ALTER TABLE "shop_item_state" ADD COLUMN     "week" INTEGER;
