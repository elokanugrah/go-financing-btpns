CREATE TABLE "users" (
  "user_id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_facility_limits" (
  "facility_limit_id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL REFERENCES "users" ("user_id"),
  "limit_amount" decimal(10, 2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tenors" (
  "tenor_id" bigserial PRIMARY KEY,
  "tenor_value" integer NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_facilities" (
  "user_facility_id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL REFERENCES "users" ("user_id"),
  "facility_limit_id" bigint NOT NULL REFERENCES "user_facility_limits" ("facility_limit_id"),
  "amount" decimal(10, 2) NOT NULL,
  "tenor" integer NOT NULL,
  "start_date" timestamptz NULL,
  "monthly_installment" decimal(10, 2) NOT NULL,
  "total_margin" decimal(10, 2) NOT NULL,
  "total_payment" decimal(10, 2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_facility_details" (
  "detail_id" bigserial PRIMARY KEY,
  "user_facility_id" bigint NOT NULL REFERENCES "user_facilities" ("user_facility_id"),
  "due_date" timestamptz NULL,
  "installment_amount" decimal(10, 2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);