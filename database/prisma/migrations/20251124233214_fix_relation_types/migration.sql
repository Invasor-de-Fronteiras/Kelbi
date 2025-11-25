/*
  Warnings:

  - You are about to alter the column `user_id` on the `characters` table. The data in that column could be lost. The data in that column will be cast from `BigInt` to `Integer`.
  - You are about to alter the column `guild_id` on the `guild_characters` table. The data in that column could be lost. The data in that column will be cast from `BigInt` to `Integer`.
  - You are about to alter the column `character_id` on the `guild_characters` table. The data in that column could be lost. The data in that column will be cast from `BigInt` to `Integer`.

*/
-- DropForeignKey
ALTER TABLE "characters" DROP CONSTRAINT "characters_user_id_fkey";

-- DropForeignKey
ALTER TABLE "guild_characters" DROP CONSTRAINT "guild_characters_character_id_fkey";

-- DropForeignKey
ALTER TABLE "guild_characters" DROP CONSTRAINT "guild_characters_guild_id_fkey";

-- AlterTable
ALTER TABLE "characters" ALTER COLUMN "user_id" SET DATA TYPE INTEGER;

-- AlterTable
ALTER TABLE "guild_characters" ALTER COLUMN "guild_id" SET DATA TYPE INTEGER,
ALTER COLUMN "character_id" SET DATA TYPE INTEGER;

-- AddForeignKey
ALTER TABLE "characters" ADD CONSTRAINT "characters_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "guild_characters" ADD CONSTRAINT "guild_characters_character_id_fkey" FOREIGN KEY ("character_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "guild_characters" ADD CONSTRAINT "guild_characters_guild_id_fkey" FOREIGN KEY ("guild_id") REFERENCES "guilds"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
