CREATE TYPE "UserLang" AS ENUM ('EN', 'JP');

ALTER TABLE "users" ADD COLUMN "language" "UserLang" NOT NULL DEFAULT 'EN';