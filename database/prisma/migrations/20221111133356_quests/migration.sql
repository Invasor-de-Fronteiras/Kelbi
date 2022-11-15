-- CreateEnum

CREATE TYPE "QuestSeason" AS ENUM ('WARM', 'COLD', 'BREED');

-- CreateEnum

CREATE TYPE "QuestPeriod" AS ENUM ('DAY', 'NIGHT');

-- CreateTable

CREATE TABLE
    "quests" (
        "id" SERIAL NOT NULL,
        "quest_id" SERIAL NOT NULL,
        "period" "QuestPeriod" NOT NULL,
        "season" "QuestSeason" NOT NULL,
        "name" TEXT NOT NULL,
        "objective" TEXT NOT NULL DEFAULT '',
        "category" INTEGER NOT NULL DEFAULT 0,
        "enable" BOOLEAN NOT NULL DEFAULT false,
        "metadata" JSONB,
        "quest_bin" BYTEA,
        "quest_list_bin" BYTEA NOT NULL,
        "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "quests_pkey" PRIMARY KEY ("id")
    );