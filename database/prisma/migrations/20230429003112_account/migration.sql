-- CreateEnum
CREATE TYPE "UserLang" AS ENUM ('EN', 'JP');

-- CreateEnum
CREATE TYPE "logkind" AS ENUM ('NEW_USER', 'NEW_CHAR', 'USER_LOGIN', 'CHAR_LOGIN', 'ROAD_CHANGES');

-- DropForeignKey
ALTER TABLE "characters" DROP CONSTRAINT "characters_user_id_fkey";

-- AlterTable
ALTER TABLE "users" ADD COLUMN     "account_id" INTEGER,
ADD COLUMN     "language" "UserLang" NOT NULL DEFAULT 'EN';

-- CreateTable
CREATE TABLE "accounts" (
    "id" SERIAL NOT NULL,
    "discord_id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Account_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "logging" (
    "id" SERIAL NOT NULL,
    "kind" "logkind" NOT NULL,
    "user_id" INTEGER NOT NULL,
    "character_id" INTEGER,
    "data" JSONB,
    "created_at" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "logging_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "unique_discord_id" ON "accounts"("discord_id");

-- AddForeignKey
ALTER TABLE "characters" ADD CONSTRAINT "characters_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "users" ADD CONSTRAINT "User_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "accounts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "rengoku_score" ADD CONSTRAINT "character_id" FOREIGN KEY ("character_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
