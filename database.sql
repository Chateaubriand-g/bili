
CREATE TABLE IF NOT EXISTS `users` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(64) NOT NULL UNIQUE,
    `password` varchar(64) NOT NULL,
    `nickname` varchar(32) DEFAULT 'guet',
    `email` varchar(128) DEFAULT NULL,
    `avatar` varchar(255) DEFAULT NULL,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX idx_username (`username`),
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;