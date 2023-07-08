CREATE OR REPLACE FUNCTION NOTIFY_MESSAGE_CHANNEL() 
RETURNS TRIGGER AS 
	$$ DECLARE user_json json;
	BEGIN
	SELECT
	    row_to_json(u) INTO user_json
	FROM users u
	WHERE u.id = NEW.user_id;
	PERFORM pg_notify(
	    'message_channel',
	    json_build_object(
	        'id',
	        NEW.id,
	        'content',
	        NEW.content,
	        'created_at',
	        NEW.created_at,
	        'user',
	        user_json
	    ):: text
	);
	RETURN 
NEW; 

END;

$$ LANGUAGE plpgsql;

CREATE TRIGGER MESSAGE_INSERT_TRIGGER 
	AFTER
	INSERT ON messages FOR EACH ROW
	EXECUTE
	    PROCEDURE notify_message_channel();
