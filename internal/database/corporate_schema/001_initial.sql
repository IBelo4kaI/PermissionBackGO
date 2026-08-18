-- Схема внешней (общей) БД corporate. Дамп предоставлен как есть, эта БД
-- находится вне зоны ответственности этого сервиса: миграции сюда никогда не
-- применяются (не подключено к internal/database/migrate.go), файл нужен
-- только как вход для `sqlc generate`, чтобы иметь типизированный доступ к
-- реальной структуре. При изменении структуры на стороне владельца БД этот
-- файл нужно обновлять вручную, синхронно с реальностью.
CREATE TABLE `companies` (
  `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `inn` varchar(12) COLLATE utf8mb4_unicode_ci NOT NULL,
  `kpp` varchar(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `departments` (
  `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `company_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `parent_id` char(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `code` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `employments` (
  `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `person_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `company_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `department_id` char(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `position_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone_work` varchar(300) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email_work` varchar(300) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `use_timetrack` tinyint(1) NOT NULL DEFAULT '0',
  `hired_at` date NOT NULL,
  `dismissed_at` date DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `role_department` enum('head','deputy','subordinate') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `persons` (
  `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `first_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `last_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `patronymic` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `birthday` date DEFAULT NULL,
  `snils` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `auth_user_id` char(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `positions` (
  `id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `company_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


ALTER TABLE `companies`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_companies_inn` (`inn`);

ALTER TABLE `departments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_dept_company` (`company_id`),
  ADD KEY `fk_dept_parent` (`parent_id`);

ALTER TABLE `employments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_emp_person` (`person_id`),
  ADD KEY `fk_emp_company` (`company_id`),
  ADD KEY `fk_emp_department` (`department_id`),
  ADD KEY `fk_emp_position` (`position_id`);

ALTER TABLE `persons`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uk_persons_snils` (`snils`),
  ADD UNIQUE KEY `uk_persons_phone` (`phone`),
  ADD UNIQUE KEY `uk_persons_email` (`email`),
  ADD UNIQUE KEY `uk_persons_auth_user` (`auth_user_id`);

ALTER TABLE `positions`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_pos_company` (`company_id`);


ALTER TABLE `departments`
  ADD CONSTRAINT `fk_dept_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_dept_parent` FOREIGN KEY (`parent_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL;

ALTER TABLE `employments`
  ADD CONSTRAINT `fk_emp_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_emp_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `fk_emp_person` FOREIGN KEY (`person_id`) REFERENCES `persons` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_emp_position` FOREIGN KEY (`position_id`) REFERENCES `positions` (`id`) ON DELETE RESTRICT;

ALTER TABLE `positions`
  ADD CONSTRAINT `fk_pos_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE;
