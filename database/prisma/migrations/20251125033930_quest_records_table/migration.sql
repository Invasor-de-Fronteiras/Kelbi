-- CreateEnum
CREATE TYPE "QuestRecordStatus" AS ENUM ('COMPLETED', 'FAILED', 'RETURNED');

-- CreateTable
CREATE TABLE "quest_records" (
    "record_id" UUID NOT NULL,
    "quest_id" INTEGER NOT NULL,
    "period" "QuestPeriod" NOT NULL,
    "season" "QuestSeason" NOT NULL,
    "user_id" INTEGER NOT NULL,
    "character_id" INTEGER NOT NULL,
    "party_size" SMALLINT NOT NULL,
    "weapon_type" INTEGER NOT NULL,
    "guild_id" INTEGER,
    "status" "QuestRecordStatus",
    "frames" INTEGER NOT NULL,
    "time" INTEGER NOT NULL,
    "started_at" TIMESTAMP(3) NOT NULL,
    "finished_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CreateIndex
CREATE INDEX "idx_quest_records_character_id" ON "quest_records"("character_id");

-- CreateIndex
CREATE INDEX "idx_quest_records_quest_id" ON "quest_records"("quest_id");

-- CreateIndex
CREATE INDEX "idx_quest_records_quest_id_record_id" ON "quest_records"("quest_id", "record_id");

-- CreateIndex
CREATE INDEX "idx_quest_records_record_id" ON "quest_records"("record_id");

-- CreateIndex
CREATE INDEX "idx_quest_records_user_id" ON "quest_records"("user_id");
