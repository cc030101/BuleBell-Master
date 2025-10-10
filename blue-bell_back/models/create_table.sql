CREATE TABLE `user`(
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    `user_id`   BIGINT(20)  DEFAULT NULL,
    'username'  VARCHAR(64) NOT NULL DEFAULT '',
    "password" VARCHAR(64) NOT  NULL    DEFAULT '',
    "email" VARCHAR(64) DEFAULT NULL,
    "gender"    TINYINT(3) DEFAULT '0',
    "create_time"   DATETIME    DEFAULT NULL,
    "update_time"   DATETIME    DEFAULT NULL,
    PRIMARY KEY (`id`)
)ENGINE=INNODB DEFAULT  CHARSET=utf8;