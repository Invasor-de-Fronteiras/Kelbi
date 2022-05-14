-- CreateEnum
CREATE TYPE "festival_colour" AS ENUM ('none', 'red', 'blue');

-- CreateEnum
CREATE TYPE "guild_application_type" AS ENUM ('applied', 'invited');

-- CreateTable
CREATE TABLE "account_ban" (
    "user_id" INTEGER NOT NULL,
    "title" TEXT,
    "reason" TEXT,
    "date" TEXT,
    "pass_origin" TEXT,
    "pass_block" TEXT,

    CONSTRAINT "ban_pkey" PRIMARY KEY ("user_id")
);

-- CreateTable
CREATE TABLE "account_history" (
    "report_id" SERIAL NOT NULL,
    "user_id" INTEGER,
    "title" TEXT,
    "reason" TEXT,
    "date" DATE,

    CONSTRAINT "account_history_pkey" PRIMARY KEY ("report_id")
);

-- CreateTable
CREATE TABLE "account_moderation" (
    "id" SERIAL NOT NULL,
    "username" TEXT,
    "password" TEXT,
    "type" TEXT,

    CONSTRAINT "account_moderation_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "account_sub" (
    "id" SERIAL NOT NULL,
    "discord_id" TEXT,
    "erupe_account" TEXT,
    "erupe_password" TEXT,
    "date_inscription" DATE,
    "country" TEXT,
    "presentation" TEXT,

    CONSTRAINT "account_auth_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "characters" (
    "id" SERIAL NOT NULL,
    "user_id" INTEGER,
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
    "gacha_trial" INTEGER,
    "gacha_prem" INTEGER,
    "gacha_items" BYTEA,
    "daily_time" TIMESTAMP(6),
    "frontier_points" INTEGER,
    "netcafe_points" INTEGER,
    "house_info" BYTEA,
    "login_boost" BYTEA,
    "skin_hist" BYTEA,
    "kouryou_point" INTEGER,
    "gcp" INTEGER,
    "trophy" BYTEA,
    "guild_post_checked" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "time_played" INTEGER NOT NULL DEFAULT 0,
    "weapon_id" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "characters_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "event_week" (
    "id" INTEGER NOT NULL,
    "event_id" INTEGER NOT NULL,
    "date_expiration" INTEGER NOT NULL,

    CONSTRAINT "event_week_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "fpoint_items" (
    "hash" INTEGER,
    "itemtype" SMALLINT,
    "itemid" INTEGER,
    "quant" INTEGER,
    "itemvalue" INTEGER,
    "tradetype" SMALLINT
);

-- CreateTable
CREATE TABLE "gacha_shop" (
    "hash" BIGINT NOT NULL,
    "reqgr" INTEGER NOT NULL,
    "reqhr" INTEGER NOT NULL,
    "gachaname" VARCHAR(255) NOT NULL,
    "gachalink0" VARCHAR(255) NOT NULL,
    "gachalink1" VARCHAR(255) NOT NULL,
    "gachalink2" VARCHAR(255) NOT NULL,
    "extraicon" INTEGER NOT NULL,
    "gachatype" INTEGER NOT NULL,
    "hideflag" BOOLEAN NOT NULL,

    CONSTRAINT "gacha_shop_pkey" PRIMARY KEY ("hash","reqgr","reqhr","gachaname","gachalink0","gachalink1","gachalink2","extraicon","gachatype","hideflag")
);

-- CreateTable
CREATE TABLE "gacha_shop_items" (
    "shophash" INTEGER,
    "entrytype" SMALLINT,
    "itemhash" INTEGER NOT NULL,
    "currtype" SMALLINT,
    "currnumber" INTEGER,
    "currquant" INTEGER,
    "percentage" INTEGER,
    "rarityicon" SMALLINT,
    "rollscount" SMALLINT,
    "itemcount" SMALLINT,
    "dailylimit" SMALLINT,
    "itemtype" INTEGER[],
    "itemid" INTEGER[],
    "quantity" INTEGER[]
);

-- CreateTable
CREATE TABLE "gook" (
    "id" SERIAL NOT NULL,
    "gook0" BYTEA,
    "gook1" BYTEA,
    "gook2" BYTEA,
    "gook3" BYTEA,
    "gook4" BYTEA,
    "gook5" BYTEA,
    "gook0status" BOOLEAN,
    "gook1status" BOOLEAN,
    "gook2status" BOOLEAN,
    "gook3status" BOOLEAN,
    "gook4status" BOOLEAN,
    "gook5status" BOOLEAN,

    CONSTRAINT "gook_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_applications" (
    "id" SERIAL NOT NULL,
    "guild_id" INTEGER NOT NULL,
    "character_id" INTEGER NOT NULL,
    "actor_id" INTEGER NOT NULL,
    "application_type" "guild_application_type" NOT NULL,
    "created_at" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "guild_applications_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guild_characters" (
    "id" SERIAL NOT NULL,
    "guild_id" INTEGER,
    "character_id" INTEGER,
    "joined_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "avoid_leadership" BOOLEAN NOT NULL DEFAULT false,
    "order_index" INTEGER NOT NULL DEFAULT 1,

    CONSTRAINT "guild_characters_pkey" PRIMARY KEY ("id")
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
    "created_at" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "liked_by" TEXT NOT NULL DEFAULT E'',

    CONSTRAINT "guild_posts_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "guilds" (
    "id" SERIAL NOT NULL,
    "name" VARCHAR(24),
    "created_at" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP,
    "leader_id" INTEGER NOT NULL,
    "main_motto" INTEGER DEFAULT 0,
    "rank_rp" INTEGER NOT NULL DEFAULT 0,
    "comment" VARCHAR(255) NOT NULL DEFAULT E'',
    "festival_colour" "festival_colour" DEFAULT E'none',
    "icon" BYTEA,
    "sub_motto" INTEGER DEFAULT 0,
    "item_box" BYTEA,
    "event_rp" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "guilds_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "history" (
    "user_id" INTEGER,
    "admin_id" INTEGER,
    "report_id" INTEGER NOT NULL,
    "title" TEXT,
    "reason" TEXT,

    CONSTRAINT "history_pkey" PRIMARY KEY ("report_id")
);

-- CreateTable
CREATE TABLE "login_boost_state" (
    "char_id" INTEGER,
    "week_req" SMALLINT,
    "week_count" SMALLINT,
    "available" BOOLEAN,
    "end_time" INTEGER,
    "ID" SERIAL NOT NULL,

    CONSTRAINT "login_boost_state_pkey" PRIMARY KEY ("ID")
);

-- CreateTable
CREATE TABLE "lucky_box_state" (
    "char_id" INTEGER,
    "shophash" INTEGER NOT NULL,
    "used_itemhash" INTEGER[]
);

-- CreateTable
CREATE TABLE "mail" (
    "id" SERIAL NOT NULL,
    "sender_id" INTEGER NOT NULL,
    "recipient_id" INTEGER NOT NULL,
    "subject" VARCHAR NOT NULL DEFAULT E'',
    "body" VARCHAR NOT NULL DEFAULT E'',
    "read" BOOLEAN NOT NULL DEFAULT false,
    "attached_item_received" BOOLEAN NOT NULL DEFAULT false,
    "attached_item" INTEGER,
    "attached_item_amount" INTEGER NOT NULL DEFAULT 1,
    "is_guild_invite" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "mail_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "normal_shop_items" (
    "shoptype" INTEGER,
    "shopid" INTEGER,
    "itemhash" INTEGER NOT NULL,
    "itemid" INTEGER,
    "points" INTEGER,
    "tradequantity" INTEGER,
    "rankreqlow" INTEGER,
    "rankreqhigh" INTEGER,
    "rankreqg" INTEGER,
    "storelevelreq" INTEGER,
    "maximumquantity" INTEGER,
    "boughtquantity" INTEGER,
    "roadfloorsrequired" INTEGER,
    "weeklyfataliskills" INTEGER,

    CONSTRAINT "normal_shop_items_pkey" PRIMARY KEY ("itemhash")
);

-- CreateTable
CREATE TABLE "questlists" (
    "ind" INTEGER NOT NULL,
    "questlist" BYTEA,

    CONSTRAINT "questlists_pkey" PRIMARY KEY ("ind")
);

-- CreateTable
CREATE TABLE "raviregister" (
    "id" SERIAL NOT NULL,
    "refid" INTEGER NOT NULL,
    "nextravi" INTEGER NOT NULL,
    "ravistarted" INTEGER,
    "raviposttime" INTEGER,
    "ravitype" INTEGER,
    "maxplayers" INTEGER,
    "ravikilled" INTEGER,
    "carvequest" INTEGER,
    "register1" INTEGER,
    "register2" INTEGER,
    "register3" INTEGER,
    "register4" INTEGER,
    "register5" INTEGER,

    CONSTRAINT "raviregister_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ravistate" (
    "id" SERIAL NOT NULL,
    "refid" INTEGER NOT NULL,
    "phase1hp" INTEGER NOT NULL,
    "phase2hp" INTEGER,
    "phase3hp" INTEGER,
    "phase4hp" INTEGER,
    "phase5hp" INTEGER,
    "phase6hp" INTEGER,
    "phase7hp" INTEGER,
    "phase8hp" INTEGER,
    "phase9hp" INTEGER,
    "unknown1" INTEGER,
    "unknown2" INTEGER,
    "unknown3" INTEGER,
    "unknown4" INTEGER,
    "unknown5" INTEGER,
    "unknown6" INTEGER,
    "unknown7" INTEGER,
    "unknown8" INTEGER,
    "unknown9" INTEGER,
    "unknown10" INTEGER,
    "unknown11" INTEGER,
    "unknown12" INTEGER,
    "unknown13" INTEGER,
    "unknown14" INTEGER,
    "unknown15" INTEGER,
    "unknown16" INTEGER,
    "unknown17" INTEGER,
    "unknown18" INTEGER,
    "unknown19" INTEGER,
    "unknown20" INTEGER,
    "damagemultiplier" INTEGER,

    CONSTRAINT "ravistate_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "ravisupport" (
    "id" SERIAL NOT NULL,
    "refid" INTEGER NOT NULL,
    "support1" INTEGER NOT NULL,
    "support2" INTEGER,
    "support3" INTEGER,
    "support4" INTEGER,
    "support5" INTEGER,
    "support6" INTEGER,
    "support7" INTEGER,
    "support8" INTEGER,
    "support9" INTEGER,
    "support10" INTEGER,
    "support11" INTEGER,
    "support12" INTEGER,
    "support13" INTEGER,
    "support14" INTEGER,
    "support15" INTEGER,
    "support16" INTEGER,
    "support17" INTEGER,
    "support18" INTEGER,
    "support19" INTEGER,
    "support20" INTEGER,
    "support21" INTEGER,
    "support22" INTEGER,
    "support23" INTEGER,
    "support24" INTEGER,
    "support25" INTEGER,

    CONSTRAINT "ravisupport_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "schema_migrations" (
    "version" BIGINT NOT NULL,
    "dirty" BOOLEAN NOT NULL,

    CONSTRAINT "schema_migrations_pkey" PRIMARY KEY ("version")
);

-- CreateTable
CREATE TABLE "servers" (
    "server_id" SERIAL NOT NULL,
    "server_name" TEXT,
    "season" INTEGER,
    "current_players" INTEGER,
    "event_id" INTEGER,
    "event_expiration" INTEGER,

    CONSTRAINT "servers_pkey" PRIMARY KEY ("server_id")
);

-- CreateTable
CREATE TABLE "shop_item_state" (
    "char_id" INTEGER,
    "itemhash" INTEGER NOT NULL,
    "usedquantity" INTEGER
);

-- CreateTable
CREATE TABLE "sign_sessions" (
    "id" SERIAL NOT NULL,
    "user_id" INTEGER,
    "auth_token_num" BIGINT,
    "auth_token_str" TEXT,

    CONSTRAINT "sign_sessions_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "stepup_state" (
    "char_id" INTEGER,
    "shophash" INTEGER NOT NULL,
    "step_progression" INTEGER,
    "step_time" TIMESTAMP(6)
);

-- CreateTable
CREATE TABLE "users" (
    "id" SERIAL NOT NULL,
    "username" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "item_box" BYTEA,
    "rights" INTEGER NOT NULL DEFAULT 14,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "gacha_shop_items_itemhash_key" ON "gacha_shop_items"("itemhash");

-- CreateIndex
CREATE INDEX "guild_application_type_index" ON "guild_applications"("application_type");

-- CreateIndex
CREATE UNIQUE INDEX "guild_application_character_id" ON "guild_applications"("guild_id", "character_id");

-- CreateIndex
CREATE UNIQUE INDEX "guild_character_unique_index" ON "guild_characters"("character_id");

-- CreateIndex
CREATE UNIQUE INDEX "id_week" ON "login_boost_state"("char_id", "week_req");

-- CreateIndex
CREATE UNIQUE INDEX "lucky_box_state_id_shophash" ON "lucky_box_state"("char_id", "shophash");

-- CreateIndex
CREATE INDEX "mail_recipient_deleted_created_id_index" ON "mail"("recipient_id", "deleted", "created_at", "id");

-- CreateIndex
CREATE UNIQUE INDEX "normal_shop_items_itemhash_key" ON "normal_shop_items"("itemhash");

-- CreateIndex
CREATE UNIQUE INDEX "shop_item_state_id_itemhash" ON "shop_item_state"("char_id", "itemhash");

-- CreateIndex
CREATE UNIQUE INDEX "stepup_state_id_shophash" ON "stepup_state"("char_id", "shophash");

-- CreateIndex
CREATE UNIQUE INDEX "users_username_key" ON "users"("username");

-- AddForeignKey
ALTER TABLE "characters" ADD CONSTRAINT "characters_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

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
ALTER TABLE "login_boost_state" ADD CONSTRAINT "login_boost_state_char_id_fkey" FOREIGN KEY ("char_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "lucky_box_state" ADD CONSTRAINT "lucky_box_state_char_id_fkey" FOREIGN KEY ("char_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "mail" ADD CONSTRAINT "mail_recipient_id_fkey" FOREIGN KEY ("recipient_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "mail" ADD CONSTRAINT "mail_sender_id_fkey" FOREIGN KEY ("sender_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "shop_item_state" ADD CONSTRAINT "shop_item_state_char_id_fkey" FOREIGN KEY ("char_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "sign_sessions" ADD CONSTRAINT "sign_sessions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- AddForeignKey
ALTER TABLE "stepup_state" ADD CONSTRAINT "stepup_state_char_id_fkey" FOREIGN KEY ("char_id") REFERENCES "characters"("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
