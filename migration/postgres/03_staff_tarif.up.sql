CREATE TABLE "staff_tarif" (
  "id" UUID PRIMARY KEY NOT NULL,
  "name" VARCHAR(35) NOT NULL,
  "type" VARCHAR(35) NOT NULL DEFAULT 'fixed',
    "amountforcash" INT NOT NULL,
    "amountforcard" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted" BOOLEAN DEFAULT false,
    "deleted_at" TIMESTAMP
);