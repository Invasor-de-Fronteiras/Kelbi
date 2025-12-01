-- CreateEnum
CREATE TYPE "ConnectionProviders" AS ENUM ('DISCORD');

-- CreateEnum
CREATE TYPE "QuestPeriod" AS ENUM ('DAY', 'NIGHT');

-- CreateEnum
CREATE TYPE "QuestSeason" AS ENUM ('WARM', 'COLD', 'BREED');

-- CreateEnum
CREATE TYPE "UserLang" AS ENUM ('EN', 'JP');

-- CreateEnum
CREATE TYPE "event_type" AS ENUM ('festa', 'diva', 'vs', 'mezfes');

-- CreateEnum
CREATE TYPE "festival_colour" AS ENUM ('none', 'red', 'blue');

-- CreateEnum
CREATE TYPE "guild_application_type" AS ENUM ('applied', 'invited');

-- CreateEnum
CREATE TYPE "logkind" AS ENUM ('NEW_USER', 'NEW_CHAR', 'USER_LOGIN', 'CHAR_LOGIN', 'ROAD_CHANGES');

-- CreateEnum
CREATE TYPE "prize_type" AS ENUM ('personal', 'guild');

-- CreateTable
CREATE TABLE "accounts" (
    "id" SERIAL NOT NULL,
    "discord_id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Account_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "achievements" (
    "id" INTEGER NOT NULL,
    "ach0" INTEGER DEFAULT 0,
    "ach1" INTEGER DEFAULT 0,
    "ach2" INTEGER DEFAULT 0,
    "ach3" INTEGER DEFAULT 0,
    "ach4" INTEGER DEFAULT 0,
    "ach5" INTEGER DEFAULT 0,
    "ach6" INTEGER DEFAULT 0,
    "ach7" INTEGER DEFAULT 0,
    "ach8" INTEGER DEFAULT 0,
    "ach9" INTEGER DEFAULT 0,
    "ach10" INTEGER DEFAULT 0,
    "ach11" INTEGER DEFAULT 0,
    "ach12" INTEGER DEFAULT 0,
    "ach13" INTEGER DEFAULT 0,
    "ach14" INTEGER DEFAULT 0,
    "ach15" INTEGER DEFAULT 0,
    "ach16" INTEGER DEFAULT 0,
    "ach17" INTEGER DEFAULT 0,
    "ach18" INTEGER DEFAULT 0,
    "ach19" INTEGER DEFAULT 0,
    "ach20" INTEGER DEFAULT 0,
    "ach21" INTEGER DEFAULT 0,
    "ach22" INTEGER DEFAULT 0,
    "ach23" INTEGER DEFAULT 0,
    "ach24" INTEGER DEFAULT 0,
    "ach25" INTEGER DEFAULT 0,
    "ach26" INTEGER DEFAULT 0,
    "ach27" INTEGER DEFAULT 0,
    "ach28" INTEGER DEFAULT 0,
    "ach29" INTEGER DEFAULT 0,
    "ach30" INTEGER DEFAULT 0,
    "ach31" INTEGER DEFAULT 0,
    "ach32" INTEGER DEFAULT 0,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "achievements_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "cafe_accepted" (
    "cafe_id" INTEGER NOT NULL,
    "character_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" SERIAL NOT NULL,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "cafe_accepted_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "cafebonus" (
    "id" SERIAL NOT NULL,
    "time_req" INTEGER NOT NULL,
    "item_type" INTEGER NOT NULL,
    "item_id" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "cafebonus_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "characters" (
    "id" SERIAL NOT NULL,
    "user_id" BIGINT,
    "is_female" BOOLEAN,
    "is_new_character" BOOLEAN,
    "name" VARCHAR(15),
    "unk_desc_string" VARCHAR(31),
    "gr" INTEGER,
    "hrp" INTEGER,
    "weapon_type" INTEGER,
    "last_login" INTEGER,
    "savedata" BYTEA,
    "decomyset" BYTEA,
    "hunternavi" BYTEA,
    "otomoairou" BYTEA,
    "partner" BYTEA,
    "platebox" BYTEA,
    "platedata" BYTEA,
    "platemyset" BYTEA,
    "rengokudata" BYTEA,
    "savemercenary" BYTEA,
    "restrict_guild_scout" BOOLEAN NOT NULL DEFAULT false,
    "minidata" BYTEA,
    "gacha_items" BYTEA,
    "daily_time" TIMESTAMPTZ(6),
    "house_info" BYTEA,
    "login_boost" BYTEA,
    "skin_hist" BYTEA,
    "kouryou_point" INTEGER,
    "gcp" INTEGER,
    "guild_post_checked" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "time_played" INTEGER NOT NULL DEFAULT 0,
    "weapon_id" INTEGER NOT NULL DEFAULT 0,
    "scenariodata" BYTEA,
    "house" BYTEA,
    "savefavoritequest" BYTEA,
    "friends" TEXT NOT NULL DEFAULT '',
    "blocked" TEXT NOT NULL DEFAULT '',
    "deleted" BOOLEAN NOT NULL DEFAULT false,
    "cafe_time" INTEGER DEFAULT 0,
    "netcafe_points" INTEGER DEFAULT 0,
    "boost_time" TIMESTAMPTZ(6),
    "trophy" BYTEA,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "cafe_reset" TIMESTAMPTZ(6),
    "bonus_quests" INTEGER NOT NULL DEFAULT 0,
    "daily_quests" INTEGER NOT NULL DEFAULT 0,
    "promo_points" INTEGER NOT NULL DEFAULT 0,
    "rasta_id" INTEGER,
    "pact_id" INTEGER,
    "stampcard" INTEGER NOT NULL DEFAULT 0,
    "mezfes" BYTEA,

    CONSTRAINT "characters_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "distribution" (
    "id" SERIAL NOT NULL,
    "character_id" INTEGER,
    "type" INTEGER NOT NULL,
    "deadline" TIMESTAMPTZ(6),
    "event_name" TEXT NOT NULL DEFAULT 'GM Gift!',
    "description" TEXT NOT NULL DEFAULT '~C05You received a gift!',
    "times_acceptable" INTEGER NOT NULL DEFAULT 1,
    "min_hr" INTEGER NOT NULL DEFAULT 65535,
    "max_hr" INTEGER NOT NULL DEFAULT 65535,
    "min_sr" INTEGER NOT NULL DEFAULT 65535,
    "max_sr" INTEGER NOT NULL DEFAULT 65535,
    "min_gr" INTEGER NOT NULL DEFAULT 65535,
    "max_gr" INTEGER NOT NULL DEFAULT 65535,
    "data" BYTEA NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "distribution_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "distributions_accepted" (
    "distribution_id" INTEGER,
    "character_id" INTEGER,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" SERIAL NOT NULL,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "distributions_accepted_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "events" (
    "id" SERIAL NOT NULL,
    "start_time" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "event_type" "event_type" NOT NULL,

    CONSTRAINT "events_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "feature_weapon" (
    "start_time" TIMESTAMPTZ(6) NOT NULL,
    "featured" INTEGER NOT NULL
);

-- CreateTable
CREATE TABLE "festa_prizes" (
    "id" SERIAL NOT NULL,
    "type" "prize_type" NOT NULL,
    "tier" INTEGER NOT NULL,
    "souls_req" INTEGER NOT NULL,
    "item_id" INTEGER NOT NULL,
    "num_item" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "festa_prizes_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "festa_prizes_accepted" (
    "prize_id" INTEGER NOT NULL,
    "character_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" SERIAL NOT NULL,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "festa_prizes_accepted_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "festa_registrations" (
    "guild_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" SERIAL NOT NULL,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "team" "festival_colour" NOT NULL,

    CONSTRAINT "festa_registrations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "festa_trials" (
    "id" SERIAL NOT NULL,
    "objective" INTEGER NOT NULL,
    "goal_id" INTEGER NOT NULL,
    "times_req" INTEGER NOT NULL,
    "locale_req" INTEGER NOT NULL DEFAULT 0,
    "reward" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "festa_trials_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "fpoint_items" (
    "id" SERIAL NOT NULL,
    "item_type" INTEGER,
    "item_id" INTEGER,
    "quantity" INTEGER,
    "fpoints" INTEGER,
    "trade_type" INTEGER,

    CONSTRAINT "fpoint_items_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "gacha_box" (
    "gacha_id" INTEGER,
    "entry_id" INTEGER,
    "character_id" INTEGER
);

-- CreateTable
CREATE TABLE "gacha_entries" (
    "id" SERIAL NOT NULL,
    "gacha_id" INTEGER,
    "entry_type" INTEGER,
    "item_type" INTEGER,
    "item_number" INTEGER,
    "item_quantity" INTEGER,
    "weight" INTEGER,
    "rarity" INTEGER,
    "rolls" INTEGER,
    "frontier_points" INTEGER,
    "daily_limit" INTEGER,

    CONSTRAINT "gacha_entries_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "gacha_items" (
    "id" SERIAL NOT NULL,
    "entry_id" INTEGER,
    "item_type" INTEGER,
    "item_id" INTEGER,
    "quantity" INTEGER,

    CONSTRAINT "gacha_items_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "gacha_shop" (
    "id" SERIAL NOT NULL,
    "min_gr" INTEGER,
    "min_hr" INTEGER,
    "name" TEXT,
    "url_banner" TEXT,
    "url_feature" TEXT,
    "url_thumbnail" TEXT,
    "wide" BOOLEAN,
    "recommended" BOOLEAN,
    "gacha_type" INTEGER,
    "hidden" BOOLEAN,

    CONSTRAINT "gacha_shop_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "gacha_stepup" (
    "gacha_id" INTEGER,
    "step" INTEGER,
    "character_id" INTEGER
);

-- CreateTable
CREATE TABLE "gook" (
    "id" SERIAL NOT NULL,
    "gook0" BYTEA,
    "gook1" BYTEA,
    "gook2" BYTEA,
    "gook3" BYTEA,
    "gook4" BYTEA,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "gook_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_adventures" (
    "id" SERIAL NOT NULL,
    "guild_id" INTEGER NOT NULL,
    "destination" INTEGER NOT NULL,
    "charge" INTEGER NOT NULL DEFAULT 0,
    "depart" INTEGER NOT NULL,
    "return" INTEGER NOT NULL,
    "collected_by" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "guild_adventures_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_alliances" (
    "id" SERIAL NOT NULL,
    "name" VARCHAR(24) NOT NULL,
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "parent_id" INTEGER NOT NULL,
    "sub1_id" INTEGER,
    "sub2_id" INTEGER,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "guild_alliances_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_applications" (
    "id" SERIAL NOT NULL,
    "guild_id" INTEGER NOT NULL,
    "character_id" INTEGER NOT NULL,
    "actor_id" INTEGER NOT NULL,
    "application_type" "guild_application_type" NOT NULL,
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "guild_applications_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_characters" (
    "id" SERIAL NOT NULL,
    "guild_id" BIGINT,
    "character_id" BIGINT,
    "joined_at" TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,
    "avoid_leadership" BOOLEAN NOT NULL DEFAULT false,
    "order_index" INTEGER NOT NULL DEFAULT 1,
    "recruiter" BOOLEAN NOT NULL DEFAULT false,
    "souls" INTEGER DEFAULT 0,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "rp_today" INTEGER DEFAULT 0,
    "rp_yesterday" INTEGER DEFAULT 0,

    CONSTRAINT "guild_characters_pkey" PRIMARY KEY ("id")
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
    "hunters" TEXT NOT NULL DEFAULT '',
    "treasure" TEXT NOT NULL DEFAULT '',
    "hunt_data" BYTEA NOT NULL,
    "cats_used" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "guild_hunts_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_meals" (
    "id" SERIAL NOT NULL,
    "guild_id" INTEGER NOT NULL,
    "meal_id" INTEGER NOT NULL,
    "level" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "expires_at" TIMESTAMPTZ(6),

    CONSTRAINT "guild_meals_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_posts" (
    "id" SERIAL NOT NULL,
    "guild_id" INTEGER NOT NULL,
    "author_id" INTEGER NOT NULL,
    "post_type" INTEGER NOT NULL,
    "stamp_id" INTEGER NOT NULL,
    "title" TEXT NOT NULL,
    "body" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "liked_by" TEXT NOT NULL DEFAULT '',
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "guild_posts_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guilds" (
    "id" SERIAL NOT NULL,
    "name" VARCHAR(24),
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "leader_id" INTEGER NOT NULL,
    "main_motto" INTEGER DEFAULT 0,
    "rank_rp" INTEGER NOT NULL DEFAULT 0,
    "comment" VARCHAR(255) NOT NULL DEFAULT '',
    "icon" BYTEA,
    "sub_motto" INTEGER DEFAULT 0,
    "item_box" BYTEA,
    "event_rp" INTEGER NOT NULL DEFAULT 0,
    "pugi_name_1" VARCHAR(12) DEFAULT '',
    "pugi_name_2" VARCHAR(12) DEFAULT '',
    "pugi_name_3" VARCHAR(12) DEFAULT '',
    "recruiting" BOOLEAN NOT NULL DEFAULT true,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "pugi_outfit_1" INTEGER NOT NULL DEFAULT 0,
    "pugi_outfit_2" INTEGER NOT NULL DEFAULT 0,
    "pugi_outfit_3" INTEGER NOT NULL DEFAULT 0,
    "pugi_outfits" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "guilds_pkey" PRIMARY KEY ("id")
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

-- CreateTable
CREATE TABLE "login_boost" (
    "char_id" INTEGER,
    "week_req" INTEGER,
    "expiration" TIMESTAMPTZ(6),
    "reset" TIMESTAMPTZ(6)
);

-- CreateTable
CREATE TABLE "mail" (
    "id" SERIAL NOT NULL,
    "sender_id" INTEGER NOT NULL,
    "recipient_id" INTEGER NOT NULL,
    "subject" VARCHAR NOT NULL DEFAULT '',
    "body" VARCHAR NOT NULL DEFAULT '',
    "read" BOOLEAN NOT NULL DEFAULT false,
    "attached_item_received" BOOLEAN NOT NULL DEFAULT false,
    "attached_item" INTEGER,
    "attached_item_amount" INTEGER NOT NULL DEFAULT 1,
    "is_guild_invite" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted" BOOLEAN NOT NULL DEFAULT false,
    "locked" BOOLEAN NOT NULL DEFAULT false,
    "is_sys_message" BOOLEAN DEFAULT false,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "mail_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "questlist" (
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

    CONSTRAINT "questlist_pkey" PRIMARY KEY ("quest_id","period","season")
);

-- CreateTable
CREATE TABLE "quests" (
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

    CONSTRAINT "quests_pkey" PRIMARY KEY ("quest_id","period","season")
);

-- CreateTable
CREATE TABLE "rengoku_score" (
    "character_id" INTEGER NOT NULL,
    "max_stages_mp" INTEGER,
    "max_points_mp" INTEGER,
    "max_stages_sp" INTEGER,
    "max_points_sp" INTEGER,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "rengoku_score_pkey" PRIMARY KEY ("character_id")
);

-- CreateTable
CREATE TABLE "servers" (
    "server_id" INTEGER NOT NULL,
    "season" INTEGER NOT NULL,
    "current_players" INTEGER NOT NULL,
    "land" INTEGER,
    "world_name" TEXT,
    "world_description" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" SERIAL NOT NULL,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "servers_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "shop_items" (
    "shop_type" INTEGER,
    "shop_id" INTEGER,
    "id" SERIAL NOT NULL,
    "item_id" INTEGER,
    "cost" INTEGER,
    "quantity" INTEGER,
    "min_hr" INTEGER,
    "min_sr" INTEGER,
    "min_gr" INTEGER,
    "store_level" INTEGER,
    "max_quantity" INTEGER,
    "road_floors" INTEGER,
    "road_fatalis" INTEGER,
    "enable_weeks" VARCHAR(8),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "shop_items_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "shop_items_bought" (
    "character_id" INTEGER,
    "shop_item_id" INTEGER,
    "bought" INTEGER,
    "week" INTEGER
);

-- CreateTable
CREATE TABLE "sign_sessions" (
    "user_id" INTEGER NOT NULL,
    "char_id" INTEGER,
    "token" VARCHAR(16) NOT NULL,
    "server_id" INTEGER,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "id" SERIAL NOT NULL,

    CONSTRAINT "sign_sessions_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "stamps" (
    "character_id" INTEGER NOT NULL,
    "hl_total" INTEGER DEFAULT 0,
    "hl_redeemed" INTEGER DEFAULT 0,
    "hl_next" TIMESTAMPTZ(6),
    "ex_total" INTEGER DEFAULT 0,
    "ex_redeemed" INTEGER DEFAULT 0,
    "ex_next" TIMESTAMPTZ(6),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "stamps_pkey" PRIMARY KEY ("character_id")
);

-- CreateTable
CREATE TABLE "titles" (
    "id" INTEGER NOT NULL,
    "char_id" INTEGER NOT NULL,
    "unlocked_at" TIMESTAMPTZ(6),
    "updated_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CreateTable
CREATE TABLE "user_binary" (
    "id" SERIAL NOT NULL,
    "type2" BYTEA,
    "type3" BYTEA,
    "house_tier" BYTEA,
    "house_state" INTEGER,
    "house_password" TEXT,
    "house_data" BYTEA,
    "house_furniture" BYTEA,
    "bookshelf" BYTEA,
    "gallery" BYTEA,
    "tore" BYTEA,
    "garden" BYTEA,
    "mission" BYTEA,

    CONSTRAINT "user_binary_pkey" PRIMARY KEY ("id")
);

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

-- CreateTable
CREATE TABLE "users" (
    "id" SERIAL NOT NULL,
    "username" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "rights" INTEGER NOT NULL DEFAULT 332,
    "item_box" BYTEA,
    "last_character" INTEGER DEFAULT 0,
    "last_login" TIMESTAMPTZ(6),
    "return_expires" TIMESTAMPTZ(6),
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "dev" BOOLEAN NOT NULL DEFAULT false,
    "language" "UserLang" NOT NULL DEFAULT 'EN',
    "account_id" INTEGER,
    "gacha_premium" INTEGER,
    "gacha_trial" INTEGER,
    "frontier_points" INTEGER,
    "psn_id" TEXT,
    "wiiu_key" TEXT,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "warehouse" (
    "character_id" INTEGER NOT NULL,
    "item0" BYTEA,
    "item1" BYTEA,
    "item2" BYTEA,
    "item3" BYTEA,
    "item4" BYTEA,
    "item5" BYTEA,
    "item6" BYTEA,
    "item7" BYTEA,
    "item8" BYTEA,
    "item9" BYTEA,
    "item10" BYTEA,
    "item0name" TEXT,
    "item1name" TEXT,
    "item2name" TEXT,
    "item3name" TEXT,
    "item4name" TEXT,
    "item5name" TEXT,
    "item6name" TEXT,
    "item7name" TEXT,
    "item8name" TEXT,
    "item9name" TEXT,
    "equip0" BYTEA,
    "equip1" BYTEA,
    "equip2" BYTEA,
    "equip3" BYTEA,
    "equip4" BYTEA,
    "equip5" BYTEA,
    "equip6" BYTEA,
    "equip7" BYTEA,
    "equip8" BYTEA,
    "equip9" BYTEA,
    "equip10" BYTEA,
    "equip0name" TEXT,
    "equip1name" TEXT,
    "equip2name" TEXT,
    "equip3name" TEXT,
    "equip4name" TEXT,
    "equip5name" TEXT,
    "equip6name" TEXT,
    "equip7name" TEXT,
    "equip8name" TEXT,
    "equip9name" TEXT,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "warehouse_pkey" PRIMARY KEY ("character_id")
);

-- CreateIndex
CREATE UNIQUE INDEX "unique_discord_id" ON "accounts"("discord_id");

-- CreateIndex
CREATE INDEX "guild_application_type_index" ON "guild_applications"("application_type");

-- CreateIndex
CREATE UNIQUE INDEX "guild_application_character_id" ON "guild_applications"("guild_id", "character_id");

-- CreateIndex
CREATE UNIQUE INDEX "guild_character_unique_index" ON "guild_characters"("character_id");

-- CreateIndex
CREATE INDEX "mail_recipient_deleted_created_id_index" ON "mail"("recipient_id", "deleted", "created_at" DESC, "id" DESC);

-- CreateIndex
CREATE UNIQUE INDEX "shop_items_bought_char_id_item_id" ON "shop_items_bought"("character_id", "shop_item_id");

-- CreateIndex
CREATE UNIQUE INDEX "user_connections_id_provider" ON "user_connections"("user_id", "provider");

-- CreateIndex
CREATE UNIQUE INDEX "user_connections_provider_id_provider" ON "user_connections"("provider_id", "provider");

-- CreateIndex
CREATE UNIQUE INDEX "users_username_key" ON "users"("username");

-- AddForeignKey
ALTER TABLE "characters" ADD CONSTRAINT "characters_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "guild_applications" ADD CONSTRAINT "guild_applications_actor_id_fkey" FOREIGN KEY ("actor_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "guild_applications" ADD CONSTRAINT "guild_applications_character_id_fkey" FOREIGN KEY ("character_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "guild_applications" ADD CONSTRAINT "guild_applications_guild_id_fkey" FOREIGN KEY ("guild_id") REFERENCES "guilds"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "guild_characters" ADD CONSTRAINT "guild_characters_character_id_fkey" FOREIGN KEY ("character_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "guild_characters" ADD CONSTRAINT "guild_characters_guild_id_fkey" FOREIGN KEY ("guild_id") REFERENCES "guilds"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "mail" ADD CONSTRAINT "mail_recipient_id_fkey" FOREIGN KEY ("recipient_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "rengoku_score" ADD CONSTRAINT "character_id" FOREIGN KEY ("character_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "user_connections" ADD CONSTRAINT "user_connections_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "users" ADD CONSTRAINT "User_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "accounts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
