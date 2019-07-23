
--  user_table
DROP TABLE IF EXISTS `user`;

CREATE TABLE `user`(`id` int NOT NULL AUTO_INCREMENT,
        `username` varchar(10) NOT NULL,
        `password` varchar(32) NOT NULL,
        `avatar` varchar(100) NOT NULL,
        `permission` int NOT NULL DEFAULT 1,
        `createdAt` timestamp DEFAULT CURRENT_TIMESTAMP,
        `updatedAt` timestamp DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; 
