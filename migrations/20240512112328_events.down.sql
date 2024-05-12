DROP TRIGGER IF EXISTS events_after_insert_trigger ON events;
DROP FUNCTION IF EXISTS events_after_insert_trigger();
DROP TABLE IF EXISTS events;
drop table if exists locker;