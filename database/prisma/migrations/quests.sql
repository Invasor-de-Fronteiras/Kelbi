CREATE TABLE "quests" (
    "id" SERIAL NOT NULL,
    "quest_id" SERIAL NOT NULL,
    "period" VARCHAR(1) NOT NULL,
    "season" SMALLINT NOT NULL,
    "name" TEXT NOT NULL,
    "objective" TEXT NOT NULL,
    "category" INTEGER NOT NULL DEFAULT 0,
    "enabled" BOOLEAN NOT NULL DEFAULT false,
    "metadata" TEXT NOT NULL DEFAULT "",
    "quest_bin" BYTEA NOT NULL,
    "quest_bin_size" INTEGER NOT NULL,
    "quest_list_bin" BYTEA NOT NULL,
    "quest_list_bin_size" INTEGER NOT NULL,

    CONSTRAINT "id" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "quests_key" ON "quests"("quest_id", "period", "season");

INSERT INTO public.quests(
	quest_id, period, season, name, objective, category, enabled, metadata, quest_bin, quest_bin_size, quest_list_bin, quest_list_bin_size)
	VALUES (
		,
		?,
		?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?
    );