ALTER TABLE `users`
ADD COLUMN `phone` varchar(32) DEFAULT NULL AFTER `nickname`;
