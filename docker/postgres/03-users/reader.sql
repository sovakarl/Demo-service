DROP USER IF EXISTS user_reader;
CREATE USER  user_reader WITH LOGIN PASSWORD 'read_password_123';

GRANT role_reader TO user_reader;