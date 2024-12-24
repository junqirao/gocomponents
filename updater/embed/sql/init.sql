CREATE TABLE IF NOT EXISTS srv_version
(
    `id`         BIGINT      NOT NULL AUTO_INCREMENT,
    `name`       VARCHAR(64) NOT NULL,
    `type`       TINYINT     NOT NULL COMMENT '0: function, 1: sql',
    `status`     TINYINT     NOT NULL COMMENT '0: failed, 1: success',
    `cost`       BIGINT      NOT NULL COMMENT 'ms',
    `created_at` DATETIME    NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name_type` (`name`, `type`)
);