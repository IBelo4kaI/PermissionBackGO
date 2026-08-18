CREATE TABLE `invites` (
	`id` varchar(36) NOT NULL,
	`code` varchar(64) NOT NULL,
	`email` varchar(255) default NULL,
	`person_id` varchar(36) default NULL,
	`company_id` varchar(36) default NULL,
	`department_id` varchar(36) default NULL,
	`position_id` varchar(36) default NULL,
	`created_by` varchar(36) NOT NULL,
	`user_id` varchar(36) default NULL,
	`used` tinyint(1) NOT NULL default 0,
	`revoked` tinyint(1) NOT NULL default 0,
	`expires_at` datetime NOT NULL,
	`created_at` datetime NOT NULL default current_timestamp
) engine = innodb default charset = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

ALTER TABLE `invites`
ADD PRIMARY KEY (`id`),
ADD UNIQUE KEY `code` (`code`),
ADD KEY `created_by` (`created_by`),
ADD KEY `user_id` (`user_id`),
ADD KEY `email` (`email`),
ADD KEY `company_id` (`company_id`),
ADD KEY `department_id` (`department_id`),
ADD KEY `position_id` (`position_id`);

ALTER TABLE `invites`
ADD CONSTRAINT `invites_ibfk_1` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
ADD CONSTRAINT `invites_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE RESTRICT;
