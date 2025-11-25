-- AlterTable
ALTER TABLE "gacha_entries" ADD COLUMN     "position" INTEGER NOT NULL DEFAULT 1;

-- AlterTable
ALTER TABLE "gacha_items" ADD COLUMN     "position" INTEGER NOT NULL DEFAULT 1;

-- AlterTable
ALTER TABLE "gacha_shop" ADD COLUMN     "position" INTEGER NOT NULL DEFAULT 1;

-- AlterTable
ALTER TABLE "shop_items" ADD COLUMN     "position" INTEGER NOT NULL DEFAULT 1;
