CREATE TABLE `genders` (`id` varchar(36) NOT NULL, `name` varchar(15) NOT NULL) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `permissions` (
	`id` varchar(36) NOT NULL,
	`service_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci default NULL,
	`code` varchar(100) NOT NULL,
	`name` varchar(100) NOT NULL,
	`description` text NOT NULL,
	`created_at` datetime NOT NULL default current_timestamp
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `roles` (
	`id` varchar(36) NOT NULL,
	`service_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci default NULL,
	`name` varchar(100) NOT NULL,
	`description` text NOT NULL,
	`is_global` tinyint(1) NOT NULL,
	`created_at` datetime NOT NULL default current_timestamp
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `role_permissions` (
	`role_id` varchar(36) NOT NULL,
	`permission_id` varchar(36) NOT NULL,
	`granted_at` datetime NOT NULL default current_timestamp
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `services` (
	`id` varchar(36) NOT NULL,
	`name` varchar(100) NOT NULL,
	`prefix` varchar(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
	`theme` varchar(30) default NULL,
	`description` text NOT NULL,
	`image_url` text,
	`url` text,
	`created_at` datetime NOT NULL default current_timestamp
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `sessions` (
	`id` varchar(36) NOT NULL,
	`token_hash` varchar(64) NOT NULL,
	`user_id` varchar(36) NOT NULL,
	`expires_at` timestamp NOT NULL
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `users` (
	`id` varchar(36) NOT NULL,
	`name` varchar(100) NOT NULL,
	`surname` varchar(100) NOT NULL,
	`patronymic` varchar(100) default NULL,
	`username` varchar(100) NOT NULL,
	`password` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
	`gender_id` varchar(36) NOT NULL,
	`birthday` date NOT NULL,
	`status` enum('active', 'no-active') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL default 'active',
	`created_at` datetime NOT NULL default current_timestamp
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE `user_roles` (
	`user_id` varchar(36) NOT NULL,
	`role_id` varchar(36) NOT NULL,
	`assigned_at` datetime NOT NULL default current_timestamp
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

ALTER TABLE `genders`
ADD PRIMARY KEY (`id`),
ADD UNIQUE KEY `NAME_UNIQUE` (`name`);

ALTER TABLE `permissions`
ADD PRIMARY KEY (`id`),
ADD UNIQUE KEY `code` (`code`) USING btree,
ADD KEY `service_id_2` (`service_id`);

ALTER TABLE `roles`
ADD PRIMARY KEY (`id`),
ADD KEY `service_id` (`service_id`);

ALTER TABLE `role_permissions`
ADD UNIQUE KEY `role_id_2` (`role_id`, `permission_id`),
ADD KEY `role_id` (`role_id`),
ADD KEY `permission_id` (`permission_id`);

ALTER TABLE `services`
ADD PRIMARY KEY (`id`),
ADD UNIQUE KEY `name` (`name`),
ADD UNIQUE KEY `prefix` (`prefix`);

ALTER TABLE `sessions`
ADD PRIMARY KEY (`id`),
ADD KEY `token_hash` (`token_hash`),
ADD KEY `user_id` (`user_id`);

ALTER TABLE `users`
ADD PRIMARY KEY (`id`),
ADD UNIQUE KEY `USERNAME_UNIQUE` (`username`),
ADD KEY `gender_id` (`gender_id`) USING btree;

ALTER TABLE `user_roles`
ADD KEY `user_id` (`user_id`),
ADD KEY `role_id` (`role_id`);

ALTER TABLE `permissions`
ADD CONSTRAINT `permissions_ibfk_1` FOREIGN KEY (`service_id`) REFERENCES `services` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE `roles`
ADD CONSTRAINT `roles_ibfk_1` FOREIGN KEY (`service_id`) REFERENCES `services` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `role_permissions`
ADD CONSTRAINT `role_permissions_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
ADD CONSTRAINT `role_permissions_ibfk_2` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `sessions`
ADD CONSTRAINT `sessions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE `users`
ADD CONSTRAINT `users_ibfk_1` FOREIGN KEY (`gender_id`) REFERENCES `genders` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT;

ALTER TABLE `user_roles`
ADD CONSTRAINT `user_roles_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
ADD CONSTRAINT `user_roles_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
