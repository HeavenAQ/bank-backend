DROP TABLE IF EXISTS entries;

DROP TABLE IF EXISTS transfers;

DROP TABLE IF EXISTS accounts;

DROP FUNCTION IF EXISTS update_updated_at;

DROP TRIGGER IF EXISTS on_update_updated_at ON transfers;

DROP TRIGGER IF EXISTS on_update_updated_at ON entries;
