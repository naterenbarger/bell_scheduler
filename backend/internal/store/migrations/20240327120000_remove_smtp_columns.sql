-- Remove SMTP columns from settings table
ALTER TABLE settings DROP COLUMN smtp_host;
ALTER TABLE settings DROP COLUMN smtp_port;
ALTER TABLE settings DROP COLUMN smtp_username;
ALTER TABLE settings DROP COLUMN smtp_password;
ALTER TABLE settings DROP COLUMN smtp_from; 