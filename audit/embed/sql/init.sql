CREATE TABLE IF NOT EXISTS audit_record
(
    `id`         BIGINT      NOT NULL AUTO_INCREMENT,
    `module`     VARCHAR(64) NOT NULL,
    `event`      VARCHAR(64) NOT NULL,
    `from`       VARCHAR(64) NOT NULL,
    `content`    MEDIUMTEXT,
    `created_at` DATETIME    NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_module` (`module`),
    KEY `idx_event` (`event`)
);