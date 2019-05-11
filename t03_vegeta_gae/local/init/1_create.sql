
CREATE TABLE IF NOT EXISTS `user` (
  `id` varchar(64) NOT NULL,
  `name` varchar(256) NOT NULL,
  `mail` varchar(256) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE utf8_bin;
