-- Drop and recreate companies table
DROP TABLE IF EXISTS companies;
CREATE TABLE companies (
  id SERIAL,
  name varchar(255) NOT NULL,
  created_at timestamp NOT NULL
);

-- Insert example records
INSERT INTO companies (id, name, created_at) VALUES (1, 'Google', current_timestamp), (2, 'Apple', current_timestamp), (3, 'Microsoft', current_timestamp);

-- Drop and recreate company_attributes table
DROP TABLE IF EXISTS company_attributes;
CREATE TABLE company_attributes (
  id SERIAL,
  company_id int NOT NULL,
  url varchar(255) NOT NULL
);

-- Insert example records
TRUNCATE TABLE company_attributes;
INSERT INTO company_attributes (id, company_id, url) VALUES (1, 1, 'https://google.com'), (2, 2, 'https://apple.com'), (3, 3, 'https://microsoft.com');
