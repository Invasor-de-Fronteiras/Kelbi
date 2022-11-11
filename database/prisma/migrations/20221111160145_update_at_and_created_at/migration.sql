/*
  Warnings:

  - You are about to drop the column `festival_colour` on the `guilds` table. All the data in the column will be lost.
  - You are about to drop the `event_week` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `request_import_save` table. If the table is not empty, all the data it contains will be lost.
  - Changed the type of `application_type` on the `guild_applications` table. No cast exists, the column would be dropped and recreated, which cannot be done if there is data, since the column is required.
  - Made the column `created_at` on table `guilds` required. This step will fail if there are existing NULL values in that column.
  - Added the required column `updated_at` to the `lucky_box_state` table without a default value. This is not possible if the table is not empty.
  - Added the required column `updated_at` to the `ravistate` table without a default value. This is not possible if the table is not empty.
  - Changed the type of `provider` on the `user_connections` table. No cast exists, the column would be dropped and recreated, which cannot be done if there is data, since the column is required.

*/
-- CreateEnum
CREATE TYPE "event_type" AS ENUM ('festa', 'diva', 'vs', 'mezfes');

-- CreateEnum
CREATE TYPE "prize_type" AS ENUM ('personal', 'guild');

-- DropForeignKey
ALTER TABLE "mail" DROP CONSTRAINT "mail_sender_id_fkey";

-- DropForeignKey
ALTER TABLE "request_import_save" DROP CONSTRAINT "request_import_save_userId_fkey";

-- DropIndex
DROP INDEX "mail_recipient_deleted_created_id_index";

-- AlterTable
ALTER TABLE "account_ban" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "account_history" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "account_moderation" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "account_sub" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "characters" ADD COLUMN     "boost_time" TIMESTAMP(6),
ADD COLUMN     "cafe_time" INTEGER DEFAULT 0,
ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "trophy" BYTEA,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN "netcafe_points" SET DEFAULT 0;

-- AlterTable
ALTER TABLE "distribution" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "distributions_accepted" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "id" SERIAL NOT NULL,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD CONSTRAINT "distributions_accepted_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "fpoint_items" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "id" SERIAL NOT NULL,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD CONSTRAINT "fpoint_items_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "gacha_shop" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "gacha_shop_items" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "gook" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "guild_adventures" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "guild_alliances" ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "guild_applications" ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
DROP COLUMN "application_type",
ADD COLUMN     "application_type" "guild_application_type" NOT NULL,
ALTER COLUMN "created_at" SET DATA TYPE TIMESTAMP(3);

-- AlterTable
ALTER TABLE "guild_characters" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "souls" INTEGER DEFAULT 0,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "guild_hunts" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "guild_meals" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "guild_posts" ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN "created_at" SET DATA TYPE TIMESTAMP(3);

-- AlterTable
ALTER TABLE "guilds" DROP COLUMN "festival_colour",
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN "created_at" SET NOT NULL,
ALTER COLUMN "created_at" SET DATA TYPE TIMESTAMP(3);

-- AlterTable
ALTER TABLE "history" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "login_boost_state" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "id" SERIAL NOT NULL,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD CONSTRAINT "login_boost_state_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "lucky_box_state" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "id" SERIAL NOT NULL,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL,
ADD CONSTRAINT "lucky_box_state_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "mail" ADD COLUMN     "is_sys_message" BOOLEAN DEFAULT false,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "normal_shop_items" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "questlists" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "raviregister" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "ravistate" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL;

-- AlterTable
ALTER TABLE "ravisupport" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "servers" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "id" SERIAL NOT NULL,
ADD COLUMN     "land" INTEGER,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN "world_description" DROP NOT NULL,
ALTER COLUMN "world_name" DROP NOT NULL,
ADD CONSTRAINT "servers_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "shop_item_state" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "id" SERIAL NOT NULL,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD CONSTRAINT "shop_item_state_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "sign_sessions" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "stepup_state" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "id" SERIAL NOT NULL,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD CONSTRAINT "stepup_state_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "user_binaries" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- AlterTable
ALTER TABLE "user_connections" DROP COLUMN "provider",
ADD COLUMN     "provider" "ConnectionProviders" NOT NULL;

-- AlterTable
ALTER TABLE "users" ADD COLUMN     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN     "last_login" TIMESTAMP(6),
ADD COLUMN     "return_expires" TIMESTAMP(6),
ADD COLUMN     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- DropTable
DROP TABLE "event_week";

-- DropTable
DROP TABLE "request_import_save";

-- DropEnum
DROP TYPE "RequestStatus";

-- CreateTable
CREATE TABLE "achievements" (
    "id" INTEGER NOT NULL,
    "ach0" INTEGER DEFAULT 0,
    "ach1" INTEGER DEFAULT 0,
    "ach2" INTEGER DEFAULT 0,
    "ach3" INTEGER DEFAULT 0,
    "ach4" INTEGER DEFAULT 0,
    "ach5" INTEGER DEFAULT 0,
    "ach6" INTEGER DEFAULT 0,
    "ach7" INTEGER DEFAULT 0,
    "ach8" INTEGER DEFAULT 0,
    "ach9" INTEGER DEFAULT 0,
    "ach10" INTEGER DEFAULT 0,
    "ach11" INTEGER DEFAULT 0,
    "ach12" INTEGER DEFAULT 0,
    "ach13" INTEGER DEFAULT 0,
    "ach14" INTEGER DEFAULT 0,
    "ach15" INTEGER DEFAULT 0,
    "ach16" INTEGER DEFAULT 0,
    "ach17" INTEGER DEFAULT 0,
    "ach18" INTEGER DEFAULT 0,
    "ach19" INTEGER DEFAULT 0,
    "ach20" INTEGER DEFAULT 0,
    "ach21" INTEGER DEFAULT 0,
    "ach22" INTEGER DEFAULT 0,
    "ach23" INTEGER DEFAULT 0,
    "ach24" INTEGER DEFAULT 0,
    "ach25" INTEGER DEFAULT 0,
    "ach26" INTEGER DEFAULT 0,
    "ach27" INTEGER DEFAULT 0,
    "ach28" INTEGER DEFAULT 0,
    "ach29" INTEGER DEFAULT 0,
    "ach30" INTEGER DEFAULT 0,
    "ach31" INTEGER DEFAULT 0,
    "ach32" INTEGER DEFAULT 0,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "achievements_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "cafe_accepted" (
    "cafe_id" INTEGER NOT NULL,
    "character_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" SERIAL NOT NULL,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "cafe_accepted_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "cafebonus" (
    "id" SERIAL NOT NULL,
    "time_req" INTEGER NOT NULL,
    "item_type" INTEGER NOT NULL,
    "item_id" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "cafebonus_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "events" (
    "id" SERIAL NOT NULL,
    "start_time" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "event_type" "event_type" NOT NULL,

    CONSTRAINT "events_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "festa_prizes" (
    "id" SERIAL NOT NULL,
    "type" "prize_type" NOT NULL,
    "tier" INTEGER NOT NULL,
    "souls_req" INTEGER NOT NULL,
    "item_id" INTEGER NOT NULL,
    "num_item" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "festa_prizes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "festa_prizes_accepted" (
    "prize_id" INTEGER NOT NULL,
    "character_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" SERIAL NOT NULL,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "festa_prizes_accepted_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "festa_registrations" (
    "guild_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" SERIAL NOT NULL,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "team" "festival_colour" NOT NULL,

    CONSTRAINT "festa_registrations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "festa_trials" (
    "id" SERIAL NOT NULL,
    "objective" INTEGER NOT NULL,
    "goal_id" INTEGER NOT NULL,
    "times_req" INTEGER NOT NULL,
    "locale_req" INTEGER NOT NULL DEFAULT 0,
    "reward" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "festa_trials_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "rengoku_score" (
    "character_id" INTEGER NOT NULL,
    "max_stages_mp" INTEGER,
    "max_points_mp" INTEGER,
    "max_stages_sp" INTEGER,
    "max_points_sp" INTEGER,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "rengoku_score_pkey" PRIMARY KEY ("character_id")
);

-- CreateTable
CREATE TABLE "stamps" (
    "character_id" INTEGER NOT NULL,
    "hl_total" INTEGER DEFAULT 0,
    "hl_redeemed" INTEGER DEFAULT 0,
    "hl_next" TIMESTAMP(6),
    "ex_total" INTEGER DEFAULT 0,
    "ex_redeemed" INTEGER DEFAULT 0,
    "ex_next" TIMESTAMP(6),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "stamps_pkey" PRIMARY KEY ("character_id")
);

-- CreateTable
CREATE TABLE "titles" (
    "id" INTEGER NOT NULL,
    "char_id" INTEGER NOT NULL,
    "unlocked_at" TIMESTAMP(6),
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CreateTable
CREATE TABLE "warehouse" (
    "character_id" INTEGER NOT NULL,
    "item0" BYTEA,
    "item1" BYTEA,
    "item2" BYTEA,
    "item3" BYTEA,
    "item4" BYTEA,
    "item5" BYTEA,
    "item6" BYTEA,
    "item7" BYTEA,
    "item8" BYTEA,
    "item9" BYTEA,
    "item10" BYTEA,
    "item0name" TEXT,
    "item1name" TEXT,
    "item2name" TEXT,
    "item3name" TEXT,
    "item4name" TEXT,
    "item5name" TEXT,
    "item6name" TEXT,
    "item7name" TEXT,
    "item8name" TEXT,
    "item9name" TEXT,
    "equip0" BYTEA,
    "equip1" BYTEA,
    "equip2" BYTEA,
    "equip3" BYTEA,
    "equip4" BYTEA,
    "equip5" BYTEA,
    "equip6" BYTEA,
    "equip7" BYTEA,
    "equip8" BYTEA,
    "equip9" BYTEA,
    "equip10" BYTEA,
    "equip0name" TEXT,
    "equip1name" TEXT,
    "equip2name" TEXT,
    "equip3name" TEXT,
    "equip4name" TEXT,
    "equip5name" TEXT,
    "equip6name" TEXT,
    "equip7name" TEXT,
    "equip8name" TEXT,
    "equip9name" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "warehouse_pkey" PRIMARY KEY ("character_id")
);

-- CreateIndex
CREATE INDEX "guild_application_type_index" ON "guild_applications"("application_type");

-- CreateIndex
CREATE INDEX "mail_recipient_deleted_created_id_index" ON "mail"("recipient_id", "deleted", "created_at" DESC, "id" DESC);

-- CreateIndex
CREATE UNIQUE INDEX "user_connections_id_provider" ON "user_connections"("user_id", "provider");

-- CreateIndex
CREATE UNIQUE INDEX "user_connections_provider_id_provider" ON "user_connections"("provider_id", "provider");
