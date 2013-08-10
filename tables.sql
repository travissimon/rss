
CREATE TABLE `Feed` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Title` varchar(1024) DEFAULT NULL,
  `Link` varchar(256) DEFAULT NULL,
  `Subtitle` varchar(2048) DEFAULT NULL,
  `Copyright` varchar(128) DEFAULT NULL,
  `Author` varchar(256) DEFAULT NULL,
  `PublishDate` datetime DEFAULT NULL,
  `Category` varchar(128) DEFAULT NULL,
  `Logo` varchar(256) DEFAULT NULL,
  `Icon` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `id_UNIQUE` (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=latin1;

CREATE TABLE `Feed` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Title` varchar(1024) DEFAULT NULL,
  `Link` varchar(256) DEFAULT NULL,
  `Subtitle` varchar(2048) DEFAULT NULL,
  `Copyright` varchar(128) DEFAULT NULL,
  `Author` varchar(256) DEFAULT NULL,
  `PublishDate` datetime DEFAULT NULL,
  `Category` varchar(128) DEFAULT NULL,
  `Logo` varchar(256) DEFAULT NULL,
  `Icon` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `id_UNIQUE` (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=latin1;
