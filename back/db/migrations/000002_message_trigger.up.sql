CREATE OR REPLACE FUNCTION notify_message_channel()
RETURNS TRIGGER AS $$
BEGIN
  PERFORM pg_notify('message_channel', row_to_json(NEW)::text);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER message_insert_trigger
AFTER INSERT ON messages
FOR EACH ROW
EXECUTE PROCEDURE notify_message_channel();
