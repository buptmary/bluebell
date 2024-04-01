-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id`     bigint(20) NOT NULL COMMENT '用户id',
    `username`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
    `password`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
    `email`       varchar(64) COLLATE utf8mb4_general_ci COMMENT '邮箱',
    `gender`      tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别(0-未知 1-男 2-女)',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`) USING BTREE,
    UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

-- ----------------------------
-- Table structure for community
-- ----------------------------
DROP TABLE IF EXISTS `communities`;
CREATE TABLE `communities`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `community_id`   int(10) unsigned NOT NULL COMMENT '社区id',
    `community_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '社区名称',
    `introduction`   varchar(256) COLLATE utf8mb4_general_ci NOT NULL COMMENT '社区简介',
    `create_time`    timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`    timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_community_id` (`community_id`),
    UNIQUE KEY `idx_community_name` (`community_name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

-- ----------------------------
-- Table structure for post
-- ----------------------------
DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts`
(
    `id`           bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `post_id`      bigint(20) NOT NULL COMMENT '帖子id',
    `title`        varchar(128) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '标题',
    `content`      varchar(8192) COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
    `author_id`    bigint(20) NOT NULL COMMENT '作者的用户id',
    `community_id` bigint(20) NOT NULL COMMENT '所属社区',
    `status`       tinyint(4) NOT NULL DEFAULT '1' COMMENT '帖子状态(0-禁用 1-正常)',
    `create_time`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_post_id` (`post_id`),
    KEY            `idx_author_id` (`author_id`),
    KEY            `idx_community_id` (`community_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;