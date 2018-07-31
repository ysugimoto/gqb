-- Drop and recreate companies table
DROP TABLE IF EXISTS companies;
CREATE TABLE companies (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  created_at datetime NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;

-- Insert example records
INSERT INTO companies (name, created_at) VALUES ('Google', NOW()), ('Apple', NOW()), ('Microsoft', NOW());

-- Drop and recreate company_attributes table
DROP TABLE IF EXISTS company_attributes;
CREATE TABLE company_attributes (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  company_id int(11) unsigned NOT NULL,
  url varchar(255) NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;

-- Insert example records
INSERT INTO company_attributes (company_id, url) VALUES (1, 'https://google.com'), (2, 'https://apple.com'), (3, 'https://microsoft.com');
