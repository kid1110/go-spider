BEGIN;
DROP DATABASE `db_news`;
CREATE DATABASE `db_news`;

USE db_news;

CREATE TABLE `t_news`
(
    `nid` int PRIMARY KEY auto_increment COMMENT '新闻id',
    `gid` VARCHAR(100) UNIQUE KEY NOT NULL COMMENT '爬取新闻时候的id',
    `title` varchar(500) NULL COMMENT '标题',
    `context` MEDIUMTEXT NULL COMMENT '文本内容',
    `date` VARCHAR(50) COMMENT '创建日期',
    UNIQUE idx_gid(gid),
    INDEX index_title(title),
    INDEX index_date(`date`)


)ENGINE = InnoDB
 DEFAULT CHARSET = utf8mb4;

 COMMIT;
