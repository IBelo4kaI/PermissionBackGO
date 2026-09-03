CREATE TABLE `smtp_settings` (
	`id` varchar(36) NOT NULL,
	`host` varchar(255) NOT NULL,
	`port` int NOT NULL,
	`username` varchar(255) NOT NULL,
	`password_encrypted` text NOT NULL,
	`from_address` varchar(255) NOT NULL,
	`from_name` varchar(255) default NULL,
	`use_tls` tinyint(1) NOT NULL default 1,
	`updated_at` datetime NOT NULL default current_timestamp on update current_timestamp
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `password_resets` (
	`id` varchar(36) NOT NULL,
	`user_id` varchar(36) NOT NULL,
	`token_hash` varchar(64) NOT NULL,
	`used` tinyint(1) NOT NULL default 0,
	`expires_at` datetime NOT NULL,
	`created_at` datetime NOT NULL default current_timestamp
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

ALTER TABLE `smtp_settings`
ADD PRIMARY KEY (`id`);

ALTER TABLE `password_resets`
ADD PRIMARY KEY (`id`),
ADD UNIQUE KEY `token_hash` (`token_hash`),
ADD KEY `user_id` (`user_id`);

ALTER TABLE `password_resets`
ADD CONSTRAINT `password_resets_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT;
