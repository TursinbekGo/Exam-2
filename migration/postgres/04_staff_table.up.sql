CREATE TABLE "staff" (
  "id" UUID PRIMARY KEY NOT NULL,
 "branch_id" UUID NOT NULL REFERENCES "branch"("id"),
  "tarif_id" UUID NOT NULL REFERENCES "staff_tarif"("id"),
  "type" VARCHAR NOT NULL,
  "name" VARCHAR NOT NULL,
  "balance" INT NOT NULL DEFAULT 0,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp,
  "deleted" BOOLEAN DEFAULT false,
  "deleted_at" timestamp
);
