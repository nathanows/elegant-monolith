CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = timezone('utc', now());
	RETURN NEW;
END;
$$ language 'plpgsql';
