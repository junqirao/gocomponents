CREATE TABLE IF NOT EXISTS application_auth
(
    `id`          BIGINT       NOT NULL AUTO_INCREMENT,
    `name`        VARCHAR(64)  NOT NULL COMMENT 'display name',
    `description` VARCHAR(200) NULL COMMENT 'description',
    `app_id`      VARCHAR(64)  NOT NULL COMMENT 'unique app id, for index',
    `app_key`     VARCHAR(64)  NOT NULL COMMENT 'unique app key',
    `app_secret`  VARCHAR(32)  NOT NULL COMMENT '32 bits hashed app secret',
    `created_at`  DATETIME     NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_app_id` (`app_id`),
    UNIQUE KEY `uk_app_key` (`app_key`)
);