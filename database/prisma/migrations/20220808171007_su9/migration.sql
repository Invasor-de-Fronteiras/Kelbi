/*
  Warnings:

  - You are about to drop the column `trophy` on the `characters` table. All the data in the column will be lost.
  - The primary key for the `gacha_shop` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to alter the column `hash` on the `gacha_shop` table. The data in that column could be lost. The data in that column will be cast from `BigInt` to `Integer`.
  - You are about to drop the column `gook0status` on the `gook` table. All the data in the column will be lost.
  - You are about to drop the column `gook1status` on the `gook` table. All the data in the column will be lost.
  - You are about to drop the column `gook2status` on the `gook` table. All the data in the column will be lost.
  - You are about to drop the column `gook3status` on the `gook` table. All the data in the column will be lost.
  - You are about to drop the column `gook4status` on the `gook` table. All the data in the column will be lost.
  - You are about to drop the column `gook5` on the `gook` table. All the data in the column will be lost.
  - You are about to drop the column `gook5status` on the `gook` table. All the data in the column will be lost.
  - The primary key for the `login_boost_state` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `ID` on the `login_boost_state` table. All the data in the column will be lost.
  - The primary key for the `schema_migrations` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to alter the column `version` on the `schema_migrations` table. The data in that column could be lost. The data in that column will be cast from `BigInt` to `Integer`.
  - The primary key for the `servers` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `event_expiration` on the `servers` table. All the data in the column will be lost.
  - You are about to drop the column `event_id` on the `servers` table. All the data in the column will be lost.
  - You are about to drop the column `server_name` on the `servers` table. All the data in the column will be lost.
  - You are about to drop the column `auth_token_num` on the `sign_sessions` table. All the data in the column will be lost.
  - You are about to drop the column `auth_token_str` on the `sign_sessions` table. All the data in the column will be lost.
  - Made the column `season` on table `servers` required. This step will fail if there are existing NULL values in that column.
  - Made the column `current_players` on table `servers` required. This step will fail if there are existing NULL values in that column.
  - Added the required column `token` to the `sign_sessions` table without a default value. This is not possible if the table is not empty.
  - Made the column `user_id` on table `sign_sessions` required. This step will fail if there are existing NULL values in that column.

*/

-- DropForeignKey
ALTER TABLE "sign_sessions" DROP CONSTRAINT "sign_sessions_user_id_fkey";

-- DropIndex
DROP INDEX "normal_shop_items_itemhash_key";

-- AlterTable
ALTER TABLE "characters" DROP COLUMN "trophy",
ADD COLUMN     "blocked" TEXT NOT NULL DEFAULT E'',
ADD COLUMN     "deleted" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN     "friends" TEXT NOT NULL DEFAULT E'',
ADD COLUMN     "house" BYTEA,
ADD COLUMN     "savefavoritequest" BYTEA,
ADD COLUMN     "scenariodata" BYTEA;

-- AlterTable
ALTER TABLE "gacha_shop" DROP CONSTRAINT "gacha_shop_pkey",
ALTER COLUMN "hash" SET DATA TYPE INTEGER,
ADD CONSTRAINT "gacha_shop_pkey" PRIMARY KEY ("hash", "reqgr", "reqhr", "gachaname", "gachalink0", "gachalink1", "gachalink2", "extraicon", "gachatype", "hideflag");

-- AlterTable
ALTER TABLE "gook" DROP COLUMN "gook0status",
DROP COLUMN "gook1status",
DROP COLUMN "gook2status",
DROP COLUMN "gook3status",
DROP COLUMN "gook4status",
DROP COLUMN "gook5",
DROP COLUMN "gook5status";

-- AlterTable
ALTER TABLE "guild_characters" ADD COLUMN     "recruiter" BOOLEAN NOT NULL DEFAULT false;

-- AlterTable
ALTER TABLE "guilds" ADD COLUMN     "pugi_name_1" VARCHAR(12) DEFAULT E'',
ADD COLUMN     "pugi_name_2" VARCHAR(12) DEFAULT E'',
ADD COLUMN     "pugi_name_3" VARCHAR(12) DEFAULT E'',
ADD COLUMN     "recruiting" BOOLEAN NOT NULL DEFAULT true;

-- AlterTable
ALTER TABLE "login_boost_state" DROP CONSTRAINT "login_boost_state_pkey",
DROP COLUMN "ID",
ADD COLUMN     "last_week" INTEGER;

-- AlterTable
ALTER TABLE "schema_migrations" DROP CONSTRAINT "schema_migrations_pkey",
ALTER COLUMN "version" SET DATA TYPE INTEGER,
ADD CONSTRAINT "schema_migrations_pkey" PRIMARY KEY ("version");

-- AlterTable
ALTER TABLE "servers" DROP CONSTRAINT "servers_pkey",
DROP COLUMN "event_expiration",
DROP COLUMN "event_id",
DROP COLUMN "server_name",
ALTER COLUMN "server_id" DROP DEFAULT,
ALTER COLUMN "season" SET NOT NULL,
ALTER COLUMN "current_players" SET NOT NULL;
DROP SEQUENCE "servers_server_id_seq";

-- AlterTable
ALTER TABLE "sign_sessions" DROP COLUMN "auth_token_num",
DROP COLUMN "auth_token_str",
ADD COLUMN     "char_id" INTEGER,
ADD COLUMN     "server_id" INTEGER,
ADD COLUMN     "token" VARCHAR(16) NOT NULL,
ALTER COLUMN "user_id" SET NOT NULL;

-- AlterTable
ALTER TABLE "users" ADD COLUMN     "last_character" INTEGER DEFAULT 0,
ALTER COLUMN "rights" SET DEFAULT 846;

-- CreateTable
CREATE TABLE "distribution" (
    "id" SERIAL NOT NULL,
    "character_id" INTEGER,
    "type" INTEGER NOT NULL,
    "deadline" TIMESTAMP(6),
    "event_name" TEXT NOT NULL DEFAULT E'GM Gift!',
    "description" TEXT NOT NULL DEFAULT E'~C05You received a gift!',
    "times_acceptable" INTEGER NOT NULL DEFAULT 1,
    "min_hr" INTEGER NOT NULL DEFAULT 65535,
    "max_hr" INTEGER NOT NULL DEFAULT 65535,
    "min_sr" INTEGER NOT NULL DEFAULT 65535,
    "max_sr" INTEGER NOT NULL DEFAULT 65535,
    "min_gr" INTEGER NOT NULL DEFAULT 65535,
    "max_gr" INTEGER NOT NULL DEFAULT 65535,
    "data" BYTEA NOT NULL,

    CONSTRAINT "distribution_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "distributions_accepted" (
    "distribution_id" INTEGER,
    "character_id" INTEGER
);

-- CreateTable
CREATE TABLE "guild_adventures" (
    "id" SERIAL NOT NULL,
    "guild_id" INTEGER NOT NULL,
    "destination" INTEGER NOT NULL,
    "charge" INTEGER NOT NULL DEFAULT 0,
    "depart" INTEGER NOT NULL,
    "return" INTEGER NOT NULL,
    "collected_by" TEXT NOT NULL DEFAULT E'',

    CONSTRAINT "guild_adventures_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_alliances" (
    "id" SERIAL NOT NULL,
    "name" VARCHAR(24) NOT NULL,
    "created_at" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "parent_id" INTEGER NOT NULL,
    "sub1_id" INTEGER,
    "sub2_id" INTEGER,

    CONSTRAINT "guild_alliances_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_hunts" (
    "id" SERIAL NOT NULL,
    "guild_id" INTEGER NOT NULL,
    "host_id" INTEGER NOT NULL,
    "destination" INTEGER NOT NULL,
    "level" INTEGER NOT NULL,
    "return" INTEGER NOT NULL,
    "acquired" BOOLEAN NOT NULL DEFAULT false,
    "claimed" BOOLEAN NOT NULL DEFAULT false,
    "hunters" TEXT NOT NULL DEFAULT E'',
    "treasure" TEXT NOT NULL DEFAULT E'',
    "hunt_data" BYTEA NOT NULL,
    "cats_used" TEXT NOT NULL,

    CONSTRAINT "guild_hunts_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_meals" (
    "id" SERIAL NOT NULL,
    "guild_id" INTEGER NOT NULL,
    "meal_id" INTEGER NOT NULL,
    "level" INTEGER NOT NULL,
    "expires" INTEGER NOT NULL,

    CONSTRAINT "guild_meals_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "user_binaries" (
    "id" INTEGER NOT NULL,
    "type2" BYTEA,
    "type3" BYTEA,

    CONSTRAINT "user_binaries_pkey" PRIMARY KEY ("id")
);
