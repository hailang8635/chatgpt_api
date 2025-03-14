CREATE TABLE `t_keywords` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `catalog` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '类别',
  `fromuser` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户',
  `keyword` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '关键字',
  `answer` text COLLATE utf8mb4_unicode_ci COMMENT '答案',
  `url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '长问答的url地址',
  `labels` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '标签',
  `is_finished` tinyint(255) DEFAULT '-1' COMMENT '是否已返回给用户',
  `is_done` tinyint(255) DEFAULT '-1' COMMENT 'GPT是否已查得',
  `create_time` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `finish_time` datetime(3) DEFAULT NULL COMMENT '结束时间',
  `update_time` datetime(3) DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user` (`fromuser`) USING BTREE,
  KEY `idx_keyword` (`keyword`) USING BTREE,
  KEY `idx_create_time` (`create_time`) USING BTREE,
  FULLTEXT KEY `idx_answer` (`answer`)
) ENGINE=InnoDB AUTO_INCREMENT=2064 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;