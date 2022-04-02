-- Creates the PostgreSQL database and user
DROP DATABASE IF EXISTS cartedepoezii_dev_db;
CREATE DATABASE cartedepoezii_dev_db
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;
COMMENT ON DATABASE cartedepoezii_dev_db
    IS 'The database for cartedepoezii.';
DROP ROLE IF EXISTS cartedepoezii_dev;
CREATE ROLE
    cartedepoezii_dev
    WITH
    LOGIN
    REPLICATION
    BYPASSRLS
    PASSWORD 'cartedepoezii_dev_pwd';
GRANT ALL ON DATABASE cartedepoezii_dev_db TO cartedepoezii_dev;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
