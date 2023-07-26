CREATE TABLE "staff_transaction" (
  "id" UUID PRIMARY KEY NOT NULL,
  "sales_id" UUID NOT NULL REFERENCES "sales"("id"),
  "type" VARCHAR NOT NULL,
  "source_type" VARCHAR NOT NULL,
  "text" TEXT NOT NULL,
  "amount" INT NOT NULL,
  "staff_id" UUID NOT NULL  REFERENCES "staff"("id"),
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp,
  "deleted" BOOLEAN DEFAULT false,
  "deleted_at" timestamp
);
