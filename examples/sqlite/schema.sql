-- Drop and recreate companies table
DROP TABLE If EXISTS companies;
CREATE TABLE companies (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  created_at TEXT
);

-- Insert example records
INSERT INTO companies (id, name, created_at) VALUES(1, 'Google', CURRENT_TIMESTAMP), (2, 'Apple', CURRENT_TIMESTAMP), (3, 'Microsoft', CURRENT_TIMESTAMP);

-- Drop and recreate company_attributes table
DROP TABLE IF EXISTS company_attributes;
CREATE TABLE company_attributes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  company_id INT,
  url TEXT
);

-- Insert example records
INSERT INTO company_attributes (company_id, url) VALUES (1, 'https://google.com'), (2, 'https://apple.com'), (3, 'https://microsoft.com');
