
CREATE TABLE IF NOT EXISTS `users` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(64) NOT NULL UNIQUE,
    `password` varchar(64) NOT NULL,
    `nickname` varchar(32) DEFAULT 'guet',
    `gender` varchar(8) DEFAULT NULL,
    `description` text DEFAULT '这是简介.',
    `email` varchar(128) DEFAULT NULL,
    `avatar` varchar(255) DEFAULT NULL,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `videos` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `length` int DEFAULT 0,
    `file_url` varchar(255) NOT NULL,
    `cover_url` varchar(255) DEFAULT '',
    `name` varchar(200) NOT NULL,
    `intro` varchar(500) DEFAULT '这是简介.',
    `owner_id` id NOT NULL,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX idx_owner_id (`owner_id`),
    INDEX idx_create_time (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=128 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;