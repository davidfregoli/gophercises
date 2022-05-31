DROP TABLE IF EXISTS `redirects`;

CREATE TABLE `redirects` (
  `from` varchar(255) NOT NULL,
  `to` varchar(255) NOT NULL,
  PRIMARY KEY (`from`)
);

INSERT INTO `redirects` VALUES (
  '/vods',
  'https://www.youtube.com/playlist?list=PLFs19LVskfNzQLZkGG_zf6yfYTp_3v_e6'
);
