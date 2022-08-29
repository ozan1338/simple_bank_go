CREATE TABLE `users` (
  `username` varchar(255) PRIMARY KEY,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `password_change_at` timestamp  DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

-- CREATE UNIQUE INDEX `accounts_index_1` ON `accounts` (`owner`, `currency`);
ALTER TABLE `accounts` ADD CONSTRAINT `owner_currency_key` UNIQUE (`owner`, `currency`);