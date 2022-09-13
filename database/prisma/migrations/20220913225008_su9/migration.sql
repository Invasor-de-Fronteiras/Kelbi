-- CreateEnum
CREATE TYPE "RequestStatus" AS ENUM ('PENDING', 'REJECTED', 'ACCEPT_RAW', 'ACCEPT_AND_MODIFIED');

-- DropIndex
DROP INDEX "mail_recipient_deleted_created_id_index";

-- AlterTable
ALTER TABLE "sign_sessions" ADD COLUMN     "id" SERIAL NOT NULL,
ADD CONSTRAINT "sign_sessions_pkey" PRIMARY KEY ("id");

-- CreateTable
CREATE TABLE "request_import_save" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER NOT NULL,
    "status" "RequestStatus" NOT NULL DEFAULT 'PENDING',
    "reviewerId" TEXT,
    "messageId" TEXT NOT NULL,
    "channelId" TEXT NOT NULL,
    "created_at" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "request_import_save_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "request_import_save_messageId_key" ON "request_import_save"("messageId");

-- CreateIndex
CREATE UNIQUE INDEX "request_import_save_channelId_key" ON "request_import_save"("channelId");

-- CreateIndex
CREATE INDEX "mail_recipient_deleted_created_id_index" ON "mail"("recipient_id", "deleted", "created_at", "id");

-- AddForeignKey
ALTER TABLE "request_import_save" ADD CONSTRAINT "request_import_save_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
