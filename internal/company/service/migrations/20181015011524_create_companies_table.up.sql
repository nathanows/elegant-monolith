CREATE TABLE companies (
	id    serial PRIMARY KEY,
	name   varchar(80) UNIQUE NOT NULL CHECK (name <> ''),
	created_at timestamp without time zone NOT NULL default timezone('utc', now()),
	updated_at timestamp without time zone NOT NULL default timezone('utc', now())
);

CREATE TRIGGER set_company_updated_at BEFORE UPDATE ON companies FOR EACH ROW EXECUTE PROCEDURE set_updated_at();
