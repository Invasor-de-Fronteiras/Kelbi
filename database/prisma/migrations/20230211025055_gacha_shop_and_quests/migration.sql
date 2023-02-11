/*
  Warnings:

  - You are about to drop the column `created_at` on the `fpoint_items` table. All the data in the column will be lost.
  - You are about to drop the column `hash` on the `fpoint_items` table. All the data in the column will be lost.
  - You are about to drop the column `itemid` on the `fpoint_items` table. All the data in the column will be lost.
  - You are about to drop the column `itemtype` on the `fpoint_items` table. All the data in the column will be lost.
  - You are about to drop the column `itemvalue` on the `fpoint_items` table. All the data in the column will be lost.
  - You are about to drop the column `quant` on the `fpoint_items` table. All the data in the column will be lost.
  - You are about to drop the column `tradetype` on the `fpoint_items` table. All the data in the column will be lost.
  - You are about to drop the column `updated_at` on the `fpoint_items` table. All the data in the column will be lost.
  - The primary key for the `gacha_shop` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `extraicon` on the `gacha_shop` table. All the data in the column will be lost.
  - You are about to drop the column `gachalink0` on the `gacha_shop` table. All the data in the column will be lost.
  - You are about to drop the column `gachalink1` on the `gacha_shop` table. All the data in the column will be lost.
  - You are about to drop the column `gachalink2` on the `gacha_shop` table. All the data in the column will be lost.
  - You are about to drop the column `gachaname` on the `gacha_shop` table. All the data in the column will be lost.
  - You are about to drop the column `gachatype` on the `gacha_shop` table. All the data in the column will be lost.
  - You are about to drop the column `hash` on the `gacha_shop` table. All the data in the column will be lost.
  - You are about to drop the column `hideflag` on the `gacha_shop` table. All the data in the column will be lost.
  - You are about to drop the column `reqgr` on the `gacha_shop` table. All the data in the column will be lost.
  - You are about to drop the column `reqhr` on the `gacha_shop` table. All the data in the column will be lost.
  - The primary key for the `quests` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `category` on the `quests` table. All the data in the column will be lost.
  - You are about to drop the column `enable` on the `quests` table. All the data in the column will be lost.
  - You are about to drop the column `id` on the `quests` table. All the data in the column will be lost.
  - You are about to drop the column `metadata` on the `quests` table. All the data in the column will be lost.
  - You are about to drop the column `name` on the `quests` table. All the data in the column will be lost.
  - You are about to drop the column `objective` on the `quests` table. All the data in the column will be lost.
  - You are about to drop the column `quest_list_bin` on the `quests` table. All the data in the column will be lost.
  - You are about to drop the `questlists` table. If the table is not empty, all the data it contains will be lost.
  - Added the required column `main_objective` to the `quests` table without a default value. This is not possible if the table is not empty.
  - Added the required column `sub_a_objective` to the `quests` table without a default value. This is not possible if the table is not empty.
  - Added the required column `sub_b_objective` to the `quests` table without a default value. This is not possible if the table is not empty.
  - Added the required column `title` to the `quests` table without a default value. This is not possible if the table is not empty.
  - Made the column `quest_bin` on table `quests` required. This step will fail if there are existing NULL values in that column.

*/
-- DropIndex
DROP INDEX "quests_quest_id_period_season_key";

-- AlterTable
ALTER TABLE "fpoint_items" DROP COLUMN "created_at",
DROP COLUMN "hash",
DROP COLUMN "itemid",
DROP COLUMN "itemtype",
DROP COLUMN "itemvalue",
DROP COLUMN "quant",
DROP COLUMN "tradetype",
DROP COLUMN "updated_at",
ADD COLUMN     "fpoints" INTEGER,
ADD COLUMN     "item_id" INTEGER,
ADD COLUMN     "item_type" INTEGER,
ADD COLUMN     "quantity" INTEGER,
ADD COLUMN     "trade_type" INTEGER;

-- AlterTable
ALTER TABLE "gacha_shop" DROP CONSTRAINT "gacha_shop_pkey",
DROP COLUMN "extraicon",
DROP COLUMN "gachalink0",
DROP COLUMN "gachalink1",
DROP COLUMN "gachalink2",
DROP COLUMN "gachaname",
DROP COLUMN "gachatype",
DROP COLUMN "hash",
DROP COLUMN "hideflag",
DROP COLUMN "reqgr",
DROP COLUMN "reqhr",
ADD COLUMN     "hide" BOOLEAN,
ADD COLUMN     "icon" INTEGER,
ADD COLUMN     "id" SERIAL NOT NULL,
ADD COLUMN     "link1" TEXT,
ADD COLUMN     "link2" TEXT,
ADD COLUMN     "link3" TEXT,
ADD COLUMN     "min_gr" INTEGER,
ADD COLUMN     "min_hr" INTEGER,
ADD COLUMN     "name" TEXT,
ADD COLUMN     "type" INTEGER,
ADD CONSTRAINT "gacha_shop_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "quests" DROP CONSTRAINT "quests_pkey",
DROP COLUMN "category",
DROP COLUMN "enable",
DROP COLUMN "id",
DROP COLUMN "metadata",
DROP COLUMN "name",
DROP COLUMN "objective",
DROP COLUMN "quest_list_bin",
ADD COLUMN     "main_objective" TEXT NOT NULL,
ADD COLUMN     "reward_item1" INTEGER,
ADD COLUMN     "reward_item2" INTEGER,
ADD COLUMN     "reward_item3" INTEGER,
ADD COLUMN     "sub_a_objective" TEXT NOT NULL,
ADD COLUMN     "sub_b_objective" TEXT NOT NULL,
ADD COLUMN     "title" TEXT NOT NULL,
ALTER COLUMN "quest_bin" SET NOT NULL,
ADD CONSTRAINT "quests_pkey" PRIMARY KEY ("quest_id", "period", "season");

-- AlterTable
ALTER TABLE "users" ADD COLUMN     "dev" BOOLEAN NOT NULL DEFAULT false;

-- DropTable
DROP TABLE "questlists";

-- CreateTable
CREATE TABLE "questlist" (
    "quest_id" SERIAL NOT NULL,
    "period" "QuestPeriod" NOT NULL,
    "season" "QuestSeason" NOT NULL,
    "category" INTEGER NOT NULL DEFAULT 0,
    "title" TEXT NOT NULL,
    "enable" BOOLEAN NOT NULL DEFAULT false,
    "position" INTEGER NOT NULL DEFAULT 0,
    "only_dev" BOOLEAN NOT NULL DEFAULT false,
    "metadata" JSONB,
    "questlist_bin" BYTEA NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "questlist_pkey" PRIMARY KEY ("quest_id","period","season")
);
