DROP TABLE IF EXISTS `redirects`;

CREATE TABLE `redirects` (
  `from` varchar(255) NOT NULL,
  `to` varchar(255) NOT NULL,
  PRIMARY KEY (`from`)
);

INSERT INTO `redirects` VALUES (
  '/amazon',
  'https://www.amazon.it'
);
