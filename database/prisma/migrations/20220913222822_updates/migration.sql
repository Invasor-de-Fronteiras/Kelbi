/*
  Warnings:

  - The primary key for the `sign_sessions` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `id` on the `sign_sessions` table. All the data in the column will be lost.

*/
-- CreateEnum
CREATE TYPE "ConnectionProviders" AS ENUM ('DISCORD');

-- DropIndex
DROP INDEX "mail_recipient_deleted_created_id_index";

-- AlterTable
ALTER TABLE "sign_sessions" DROP CONSTRAINT "sign_sessions_pkey",
DROP COLUMN "id";

-- CreateTable
CREATE TABLE "user_connections" (
    "id" SERIAL NOT NULL,
    "user_id" INTEGER NOT NULL,
    "provider" "ConnectionProviders" NOT NULL,
    "provider_id" TEXT NOT NULL,
    "created_at" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "user_connections_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "user_connections_id_provider" ON "user_connections"("user_id", "provider");

-- CreateIndex
CREATE UNIQUE INDEX "user_connections_provider_id_provider" ON "user_connections"("provider_id", "provider");

-- CreateIndex
CREATE INDEX "mail_recipient_deleted_created_id_index" ON "mail"("recipient_id", "deleted", "created_at" DESC, "id" DESC);

-- AddForeignKey
ALTER TABLE "user_connections" ADD CONSTRAINT "user_connections_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
