CREATE TABLE IF NOT EXISTS events
(
    id           uuid                                not null,
    type_id      smallserial                         not null,
    aggregate_id uuid                                not null,
    payload      json                                not null,
    created_at   timestamp default CURRENT_TIMESTAMP not null,
    constraint events_pk
        primary key (id, aggregate_id, type_id)
);

CREATE FUNCTION events_after_insert_trigger()
    RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('events:created', CONCAT(NEW.id::text, ',', NEW.type_id::text, ',', NEW.aggregate_id::text));
    RETURN NULL;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER events_after_insert_trigger
    AFTER INSERT ON events
    FOR EACH ROW EXECUTE PROCEDURE events_after_insert_trigger();

CREATE TABLE IF NOT EXISTS locker
(
    event_id     uuid        not null,
    aggregate_id uuid        not null,
    lock_domain  varchar(50) not null,
    constraint locker_pk
        primary key (event_id, aggregate_id, lock_domain)
);