ALTER TABLE users
  ADD COLUMN auth_method VARCHAR(20) NOT NULL DEFAULT 'password';

UPDATE users
  SET auth_method = 'google'
  WHERE password = '' OR password IS NULL;

UPDATE users
  SET auth_method = 'password'
  WHERE password != '' AND password IS NOT NULL;
