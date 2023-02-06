
DROP TABLE IF EXISTS "quests";
CREATE TABLE
    "quests" (
        "quest_id" SERIAL NOT NULL,
        "period" "QuestPeriod" NOT NULL,
        "season" "QuestSeason" NOT NULL,
        "title" TEXT NOT NULL,
        "main_objective" TEXT NOT NULL,
        "sub_a_objective" TEXT NOT NULL,
        "sub_b_objective" TEXT NOT NULL,
        "reward_item1" INTEGER,
        "reward_item2" INTEGER,
        "reward_item3" INTEGER,
        "quest_bin" BYTEA NOT NULL,
        "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "quests_pkey" PRIMARY KEY ("quest_id", "period", "season")
    );

DROP TABLE IF EXISTS "questlist";
CREATE TABLE
"questlist" (
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
    CONSTRAINT "questlist_pkey" PRIMARY KEY ("quest_id", "period", "season")
);

ALTER TABLE IF EXISTS public.users
    ADD COLUMN IF NOT EXISTS dev BOOLEAN NOT NULL DEFAULT false;