-- 菜单表
CREATE TABLE `chatgpt_plus`.`chatgpt_menus` (
    `id` INT(11) NOT NULL AUTO_INCREMENT ,
    `name` VARCHAR(30) NOT NULL COMMENT '菜单名称' ,
    `icon` VARCHAR(150) NOT NULL COMMENT '菜单图标' ,
     `url` VARCHAR(100) NOT NULL COMMENT '地址' ,
     `sort_num` SMALLINT(3) NOT NULL COMMENT '排序' ,
      `enabled` TINYINT(1) NOT NULL COMMENT '是否启用' ,
       PRIMARY KEY (`id`)) ENGINE = InnoDB COMMENT = '前端菜单表';

INSERT INTO `chatgpt_menus` (`id`, `name`, `icon`, `url`, `sort_num`, `enabled`) VALUES
                                                                                     (1, '对话聊天', '/images/menu/chat.png', '/chat', 0, 1),
                                                                                     (5, 'MJ 绘画', '/images/menu/mj.png', '/mj', 1, 1),
                                                                                     (6, 'SD 绘画', '/images/menu/sd.png', '/sd', 2, 1),
                                                                                     (7, '算力日志', '/images/menu/log.png', '/powerLog', 5, 1),
                                                                                     (8, '应用中心', '/images/menu/app.png', '/apps', 3, 1),
                                                                                     (9, '作品展示', '/images/menu/img-wall.png', '/images-wall', 4, 1),
                                                                                     (10, '会员计划', '/images/menu/member.png', '/member', 6, 1),
                                                                                     (11, '分享计划', '/images/menu/share.png', '/invite', 7, 1);