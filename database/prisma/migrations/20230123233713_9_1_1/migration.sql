/*
  Warnings:

  - The primary key for the `gacha_shop` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - The primary key for the `normal_shop_items` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `boughtquantity` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `itemhash` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `maximumquantity` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `points` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `rankreqg` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `rankreqhigh` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `rankreqlow` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `roadfloorsrequired` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `storelevelreq` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `tradequantity` on the `normal_shop_items` table. All the data in the column will be lost.
  - You are about to drop the column `weeklyfataliskills` on the `normal_shop_items` table. All the data in the column will be lost.
  - The primary key for the `schema_migrations` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - Added the required column `id` to the `normal_shop_items` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "characters" ADD COLUMN     "house" BYTEA,
ADD COLUMN     "trophy" BYTEA;

-- AlterTable
ALTER TABLE "gacha_shop" DROP CONSTRAINT "gacha_shop_pkey",
ALTER COLUMN "hash" SET DATA TYPE BIGINT,
ADD CONSTRAINT "gacha_shop_pkey" PRIMARY KEY ("hash", "reqgr", "reqhr", "gachaname", "gachalink0", "gachalink1", "gachalink2", "extraicon", "gachatype", "hideflag");

-- AlterTable
ALTER TABLE "normal_shop_items" DROP CONSTRAINT "normal_shop_items_pkey",
DROP COLUMN "boughtquantity",
DROP COLUMN "itemhash",
DROP COLUMN "maximumquantity",
DROP COLUMN "points",
DROP COLUMN "rankreqg",
DROP COLUMN "rankreqhigh",
DROP COLUMN "rankreqlow",
DROP COLUMN "roadfloorsrequired",
DROP COLUMN "storelevelreq",
DROP COLUMN "tradequantity",
DROP COLUMN "weeklyfataliskills",
ADD COLUMN     "cost" INTEGER,
ADD COLUMN     "id" INTEGER NOT NULL,
ADD COLUMN     "max_quantity" INTEGER,
ADD COLUMN     "min_gr" INTEGER,
ADD COLUMN     "min_hr" INTEGER,
ADD COLUMN     "min_sr" INTEGER,
ADD COLUMN     "quantity" INTEGER,
ADD COLUMN     "req_store_level" INTEGER,
ADD COLUMN     "road_fatalis" INTEGER,
ADD COLUMN     "road_floors" INTEGER,
ADD CONSTRAINT "normal_shop_items_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "schema_migrations" DROP CONSTRAINT "schema_migrations_pkey",
ALTER COLUMN "version" SET DATA TYPE BIGINT,
ADD CONSTRAINT "schema_migrations_pkey" PRIMARY KEY ("version");

-- AlterTable
ALTER TABLE "users" ALTER COLUMN "rights" SET DEFAULT 334;
