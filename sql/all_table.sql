/*
 Navicat Premium Dump SQL

 Source Server         : PostgreSQL
 Source Server Type    : PostgreSQL
 Source Server Version : 180003 (180003)
 Source Host           : localhost:5432
 Source Catalog        : postgres
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 180003 (180003)
 File Encoding         : 65001

 Date: 17/04/2026 18:01:01
*/


-- ----------------------------
-- Table structure for chat_members
-- ----------------------------
DROP TABLE IF EXISTS "public"."chat_members";
CREATE TABLE "public"."chat_members" (
  "chat_id" int4 NOT NULL,
  "user_id" int4 NOT NULL,
  "role" varchar(20) COLLATE "pg_catalog"."default" DEFAULT 'member'::character varying,
  "joined_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."chat_members" IS '聊天室與使用者關聯（成員表）';

-- ----------------------------
-- Table structure for chats
-- ----------------------------
DROP TABLE IF EXISTS "public"."chats";
CREATE TABLE "public"."chats" (
  "id" int4 NOT NULL DEFAULT nextval('chats_id_seq'::regclass),
  "type" "public"."chat_type" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default",
  "last_message_id" int4,
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "created_by" int4 NOT NULL,
  "search_key" varchar(30) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."chats"."type" IS '聊天室類型：private 或 group';
COMMENT ON COLUMN "public"."chats"."last_message_id" IS '快取最後一筆訊息 ID，加速聊天列表查詢';
COMMENT ON COLUMN "public"."chats"."created_by" IS '建立者';
COMMENT ON COLUMN "public"."chats"."search_key" IS '輔助搜尋用的key';
COMMENT ON TABLE "public"."chats" IS '聊天室（支援私訊與群組）';

-- ----------------------------
-- Indexes structure for table chats
-- ----------------------------
CREATE UNIQUE INDEX "chats_private_key_idx" ON "public"."chats" USING btree (
  "search_key" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table chats
-- ----------------------------
ALTER TABLE "public"."chats" ADD CONSTRAINT "chats_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Table structure for message_reads
-- ----------------------------
DROP TABLE IF EXISTS "public"."message_reads";
CREATE TABLE "public"."message_reads" (
  "message_id" int4 NOT NULL,
  "user_id" int4 NOT NULL,
  "read_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON TABLE "public"."message_reads" IS '訊息已讀紀錄（支援多使用者已讀）';

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS "public"."messages";
CREATE TABLE "public"."messages" (
  "id" int4 NOT NULL DEFAULT nextval('messages_id_seq'::regclass),
  "chat_id" int4 NOT NULL,
  "sender_id" int4,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamp(6),
  "type" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT 'text'::character varying
)
;
COMMENT ON COLUMN "public"."messages"."chat_id" IS '聊天室ID';
COMMENT ON COLUMN "public"."messages"."sender_id" IS '發送者ID';
COMMENT ON COLUMN "public"."messages"."content" IS '訊息';
COMMENT ON COLUMN "public"."messages"."created_at" IS '建立時間';
COMMENT ON COLUMN "public"."messages"."deleted_at" IS '非 NULL 表示訊息已被收回';
COMMENT ON COLUMN "public"."messages"."type" IS '類型';
COMMENT ON TABLE "public"."messages" IS '聊天訊息';

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" int4 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  "username" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
  "password_hash" text COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;
COMMENT ON COLUMN "public"."users"."username" IS '使用者帳號（唯一）';
COMMENT ON TABLE "public"."users" IS '使用者基本資料';

-- ----------------------------
-- Indexes structure for table chat_members
-- ----------------------------
CREATE INDEX "idx_chat_members_user" ON "public"."chat_members" USING btree (
  "user_id" "pg_catalog"."int4_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table chat_members
-- ----------------------------
ALTER TABLE "public"."chat_members" ADD CONSTRAINT "chat_members_pkey" PRIMARY KEY ("chat_id", "user_id");

-- ----------------------------
-- Primary Key structure for table chats
-- ----------------------------
ALTER TABLE "public"."chats" ADD CONSTRAINT "chats_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table message_reads
-- ----------------------------
CREATE INDEX "idx_message_reads_user" ON "public"."message_reads" USING btree (
  "user_id" "pg_catalog"."int4_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table message_reads
-- ----------------------------
ALTER TABLE "public"."message_reads" ADD CONSTRAINT "message_reads_pkey" PRIMARY KEY ("message_id", "user_id");

-- ----------------------------
-- Indexes structure for table messages
-- ----------------------------
CREATE INDEX "idx_messages_chat_id_id" ON "public"."messages" USING btree (
  "chat_id" "pg_catalog"."int4_ops" ASC NULLS LAST,
  "id" "pg_catalog"."int4_ops" DESC NULLS FIRST
);
COMMENT ON INDEX "public"."idx_messages_chat_id_id" IS '用於聊天室訊息分頁查詢';

-- ----------------------------
-- Primary Key structure for table messages
-- ----------------------------
ALTER TABLE "public"."messages" ADD CONSTRAINT "messages_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_username_key" UNIQUE ("username");

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table chat_members
-- ----------------------------
ALTER TABLE "public"."chat_members" ADD CONSTRAINT "chat_members_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "public"."chats" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."chat_members" ADD CONSTRAINT "chat_members_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table message_reads
-- ----------------------------
ALTER TABLE "public"."message_reads" ADD CONSTRAINT "message_reads_message_id_fkey" FOREIGN KEY ("message_id") REFERENCES "public"."messages" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."message_reads" ADD CONSTRAINT "message_reads_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table messages
-- ----------------------------
ALTER TABLE "public"."messages" ADD CONSTRAINT "messages_chat_id_fkey" FOREIGN KEY ("chat_id") REFERENCES "public"."chats" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."messages" ADD CONSTRAINT "messages_sender_id_fkey" FOREIGN KEY ("sender_id") REFERENCES "public"."users" ("id") ON DELETE SET NULL ON UPDATE NO ACTION;
