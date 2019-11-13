--  User table
DROP TABLE IF EXISTS `User`;

CREATE TABLE `User`
(
    `id`           int          NOT NULL AUTO_INCREMENT COMMENT '主键',
    `username`     varchar(10)  NOT NULL COMMENT '用户名',
    `password`     varchar(42)  NOT NULL COMMENT '密码',
    `avatar`       varchar(100) NOT NULL COMMENT '头像',
    `role_id`      int          NOT NULL DEFAULT 2 COMMENT '角色id',
    `recent_time`  timestamp             DEFAULT CURRENT_TIMESTAMP COMMENT '最近登录时间',
    `created_time` timestamp             DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_time` timestamp             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '用户表';

--  Role table
DROP TABLE IF EXISTS `Role`;

CREATE TABLE `Role`
(
    `id`           int         NOT NULL COMMENT '主键',
    `name`         varchar(10) NOT NULL COMMENT '角色名称',
    `created_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '角色表';

--  Permission table
DROP TABLE IF EXISTS `Permission`;

CREATE TABLE `Permission`
(
    `id`           int         NOT NULL COMMENT '主键',
    `name`         varchar(10) NOT NULL COMMENT '权限名称',
    `role_id`      int         NOT NULL COMMENT '角色id',
    `created_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '权限表';

--  Menu_Permission table
DROP TABLE IF EXISTS `Menu_Permission`;

CREATE TABLE `Menu_Permission`
(
    `id`           int         NOT NULL COMMENT '主键',
    `name`         varchar(10) NOT NULL COMMENT '菜单名称',
    `role_id`      int         NOT NULL COMMENT '角色id',
    `created_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT '菜单权限表';
