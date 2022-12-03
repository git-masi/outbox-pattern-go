CREATE OR REPLACE FUNCTION handle_order_fulfillment_event() RETURNS TRIGGER AS $$ BEGIN PERFORM pg_notify(
		'fulfillment_event',
		json_build_object(
			-- `TG_OP` is the statement (e.g. 'INSERT' or 'UPDATE')
			'operation',
			TG_OP,
			-- `NEW` is the value of the row after the change
			'record',
			row_to_json(NEW)
		)::text
	);
RETURN NULL;
-- The return value is ignored since this is an AFTER trigger
END;
$$ LANGUAGE plpgsql;
-- The trigger will call the function whenever there is an insert or update to the order_fulfillment_messages table
CREATE TRIGGER handle_order_fulfillment_event
AFTER
INSERT
	OR
UPDATE ON order_fulfillment_messages FOR EACH ROW EXECUTE FUNCTION handle_order_fulfillment_event();