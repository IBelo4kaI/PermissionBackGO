ALTER TABLE `users`
ADD COLUMN `nickname` varchar(100) default NULL AFTER `username`;
