CREATE TABLE "accounts" (
    "id" bigserial PRIMARY KEY,
    "owner" varchar NOT NULL,
    "balance" bigint NOT NULL,
    "currency" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
    "id" bigserial PRIMARY KEY,
    "from_account_id" bigint NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    "to_account_id" bigint NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

CREATE FUNCTION update_updated_at ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

CREATE TRIGGER on_update_updated_at
    BEFORE UPDATE ON accounts FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at ();

CREATE TRIGGER on_update_updated_at
    BEFORE UPDATE ON entries FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at ();

CREATE TRIGGER on_update_updated_at
    BEFORE UPDATE ON transfers FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at ();
