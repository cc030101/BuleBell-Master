/* CREATE TABLE `user`(
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    `user_id`   BIGINT(20)  NOT NULL DEFAULT '0',
    'username'  VARCHAR(64) NOT NULL DEFAULT '',
    "password" VARCHAR(64) NOT  NULL    DEFAULT '',
    "email" VARCHAR(64) DEFAULT NULL,
    "gender"    TINYINT(3) DEFAULT '0',
    "create_time"   DATETIME    DEFAULT NULL,
    "update_time"   DATETIME    DEFAULT NULL,
    PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8; */


CREATE TABLE `user`
(
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT(20)    NOT NULL DEFAULT `0`,
    `username` VARCHAR(64)  NOT NULL  DEFAULT   '',
    `password`  VARCHAR(64) not NULL DEFAULT    '',
    `email` VARCHAR(64) DEFAULT NULL,
    `gender`   TINYINT(3)   DEFAULT '0',
    `create_time`   DATETIME    DEFAULT NULL,
    `update_time`   DATETIME    DEFAULT NULL,
    PRIMARY KEY (`id`)

)ENGINE=INNODB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8


DROP TABLE IF EXISTS    `community`;
CREATE TABLE `community`
(
    `id`    int(11) NOT NULL AUTO_INCREMENT,
    `community_id`  int(10) UNSIGNED NOT NULL,
    `community_name`    VARCHAR(128)    COLLATE utf8mb4_general_ci NOT NULL,
    `introduction` VARCHAR(256) COLLATE utf8mb4_general_ci NOT NULL,
    `create_time`   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`   TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE  KEY `idx_community_id`  (`community_id`)
    UNIQUE  KEY `idx_community_name`    (`community_name`)
)ENGINE=INNODB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
INSERT INTO `community`
VALUES ('1', '1', 'Go', 'Golang', '2024-11-01 08:10:10', '2024-11-01 08:10:10');
INSERT INTO `community`
VALUES ('2', '2', 'Run', '韩跑跑', '2024-11-01 09:00:00', '2024-11-01 09:00:00');
INSERT INTO `community`
VALUES ('3', '3', 'xttcc', '×梯梯cc', '2024-11-02 08:30:00', '2024-11-02 08:30:00');
INSERT INTO `community`
VALUES ('4', '4', 'hahaha', '哈哈哈', '2024-11-02 10:00:00', '2024-11-02 10:00:00');