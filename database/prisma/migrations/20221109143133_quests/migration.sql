-- CreateTable

CREATE TABLE
    "quests" (
        "id" SERIAL NOT NULL,
        "quest_id" SERIAL NOT NULL,
        "period" VARCHAR(1) NOT NULL,
        "season" SMALLINT NOT NULL,
        "objective" TEXT NOT NULL DEFAULT '',
        "category" INTEGER NOT NULL DEFAULT 0,
        "metadata" JSONB,
        "quest_bin" BYTEA NOT NULL,
        "quest_list_bin" BYTEA NOT NULL,
        "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(3) NOT NULL,
        CONSTRAINT "quests_pkey" PRIMARY KEY ("id")
    );

-- CreateIndex

CREATE UNIQUE INDEX "quests_quest_id_period_season_key" ON "quests"("quest_id", "period", "season");