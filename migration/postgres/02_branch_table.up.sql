CREATE TABLE "branch" (
  "id" UUID PRIMARY KEY NOT NULL,
  "name" VARCHAR(50),
  "address" VARCHAR(65),
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp,
  "deleted" BOOLEAN DEFAULT false,
  "deleted_at" timestamp
);
