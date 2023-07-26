CREATE TABLE "sales" (
  "id" UUID PRIMARY KEY NOT NULL,
  "branch_id" UUID NOT NULL REFERENCES "branch"("id"),
  "shop_assistent_id" UUID REFERENCES "staff"("id"),
  "cashier_id" UUID NOT NULL REFERENCES "staff"("id"),
  "price" INT NOT NULL,
  "payment_type" VARCHAR NOT NULL,
  "status" VARCHAR NOT NULL DEFAULT 'success',
    "client_name" VARCHAR,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted" BOOLEAN DEFAULT false,
    "deleted_at" TIMESTAMP
);
