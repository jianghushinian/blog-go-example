CREATE DATABASE nightwatch
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL DEFAULT '' COMMENT '任务名称',
  `namespace` varchar(45) NOT NULL DEFAULT '' COMMENT 'k8s namespace 名称',
  `info` TEXT NOT NULL COMMENT '任务 k8s 相关信息',
  `status` varchar(45) NOT NULL DEFAULT '' COMMENT '任务状态',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户 ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name_namespace` (`name`, `namespace`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务表';

-- test data
INSERT INTO `task` (`id`, `name`, `namespace`, `info`, `status`, `user_id`) VALUES (1, 'demo-task-1', 'default', '{"image":"alpine","command":["sleep"],"args":["60"]}', 'Normal', 1);
INSERT INTO `task` (`id`, `name`, `namespace`, `info`, `status`, `user_id`) VALUES (2, 'demo-task-2', 'demo', '{"image":"busybox","command":["sleep"],"args":["3600"]}', 'Normal', 2);
