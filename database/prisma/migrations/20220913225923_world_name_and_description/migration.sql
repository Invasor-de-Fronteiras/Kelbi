/*
  Warnings:

  - Added the required column `world_description` to the `servers` table without a default value. This is not possible if the table is not empty.
  - Added the required column `world_name` to the `servers` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "servers" ADD COLUMN     "world_description" TEXT NOT NULL,
ADD COLUMN     "world_name" TEXT NOT NULL;
