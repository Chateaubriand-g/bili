
CREATE TABLE IF NOT EXISTS `users` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(64) NOT NULL UNIQUE,
    `password` varchar(64) NOT NULL,
    `nickname` varchar(32) DEFAULT 'guet',
    `gender` varchar(8) DEFAULT NULL,
    `description` text,
    `email` varchar(128) DEFAULT NULL,
    `avatar` varchar(255) DEFAULT NULL,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `videos` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `length` bigint unsigned DEFAULT 0,
    `file_url` varchar(255) NOT NULL,
    `cover_url` varchar(255) DEFAULT '',
    `name` varchar(128) NOT NULL,
    `intro` varchar(512) DEFAULT '这是简介.',
    `owner_id` bigint unsigned NOT NULL,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX idx_owner_id (`owner_id`),
    INDEX idx_create_time (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=128 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `media_objects` (
    `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
    `owner_id` bigint unsigned NULL,
    `file_md5` varchar(64) NULL,
    `object_key` varchar(512) NOT NULL,
    `size` bigint DEFAULT 0,
    `status` tinyint DEFAULT 0,
    `create_at` datetime DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `media_chunks` (
    `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
    `file_md5` varchar(64) NOT NULL,
    `chunk_index` int NOT NULL,
    `size` bigint,
    `uploaded_at` datetime DEFAULT CURRENT_TIMESTAMP,
    UNIQUE
)

CREATE TABLE IF NOT EXISTS `notifications` (
    `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
    `user_id` bigint unsigned NOT NULL,
    `type` tinyint NOT NULL,  -- 0:reply,1:@,3:system,4:pm,5:follow
    `from_user_id` bigint unsigned DEFAULT NULL,
    `biz_id` bigint unsigned DEFAULT NULL,
    `payload` json NOT NULL DEFAULT (`{}`),
    `is_read` tinyint unsigned NOT NULL DEFAULT 0,
    `deleted_at` datetime DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_user_unread_create` (`user_id`,`is_read`,`created_at` DESC),
    INDEX `idx_user_biz` (`user_id`,`biz_id`,`type`),
    INDEX `idx_from_user` (`from_user_id`,`created_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `comments` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `video_id`  bigint unsigned NOT NULL,
    `user_id` bigint unsigned NOT NULL,
    `content` text NOT NULL,
    `parent_id` bigint unsigned DEFAULT 0, //父评论id,0表示顶级评论
    `like_count` int unsigned DEFAULT 0,
    `create_at` datatime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX `idx_video` (`video_id`)
    INDEX `idx_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS comments_like (
    user_id bigint unsigned NOT NULL,
    comment_id bigint unsigned NOT NULL,
    create_at datatime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(user_id,comment_id)
    INDEX idx_comment (comment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS videos_like (
    user_id bigint unsigned NOT NULL,
    video_id bigint unsigned NOT NULL,
    create_at datatime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(user_id,video_id)
    INDEX idx_video (video_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS folders (
    id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id bigint unsigned NOT NULL,
    folder_name varcha(128) NOT NULL DEFAULT 'default',
    create_at datatime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_usr (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS floder_items (
    folder_id bigint unsigned NOT NULL,
    video_id bigint unsigned NOT NULL,
    add_at datatime NOT NUL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (folder_id,video_id),
    INDEX idx_video (video_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 