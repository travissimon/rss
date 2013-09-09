CREATE SCHEMA 'rss';

use rss;


CREATE TABLE `Feed` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Title` varchar(1024) DEFAULT NULL,
  `Link` varchar(256) DEFAULT NULL,
  `Url` varchar(1024) DEFAULT NULL,
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

CREATE TABLE `Entry` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `FeedId` int(11) NOT NULL,
  `Title` varchar(2048) DEFAULT NULL,
  `Link` varchar(2048) DEFAULT NULL,
  `Subtitle` varchar(4096) DEFAULT NULL,
  `Guid` varchar(512) DEFAULT NULL,
  `UpdatedDate` datetime DEFAULT NULL,
  `Summary` varchar(4096) DEFAULT NULL,
  `Content` varchar(4096) DEFAULT NULL,
  `Source` varchar(4096) DEFAULT NULL,
  `Comments` varchar(2048) DEFAULT NULL,
  `Thumbnail` varchar(1024) DEFAULT NULL,
  `Length` int(11) DEFAULT NULL,
  `Type` varchar(1024) DEFAULT NULL,
  `Url` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `id_UNIQUE` (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=latin1;

ALTER TABLE `rss`.`Entry` ADD INDEX `fk_Entry_Feed_idx` (`FeedId` ASC);
ALTER TABLE `rss`.`Entry` 
  ADD CONSTRAINT `fk_Entry_Feed`
  FOREIGN KEY (`FeedId`)
  REFERENCES `rss`.`Feed` (`Id`)
  ON DELETE NO ACTION
  ON UPDATE NO ACTION;

CREATE TABLE `Subscription` (
  `UserId` int(11) NOT NULL,
  `FeedId` int(11) NOT NULL,
  `UnreadItems` int(11)
) ENGINE=InnoDB CHARSET=latin1;

ALTER TABLE `rss`.`Subscription` ADD INDEX `fk_Subscription_User_idx` (`UserId` ASC);
ALTER TABLE `rss`.`Subscription` 
  ADD CONSTRAINT `fk_Subscription_UserId`
  FOREIGN KEY (`UserId`)
  REFERENCES `rss`.`User` (`Id`)
  ON DELETE NO ACTION
  ON UPDATE NO ACTION;

ALTER TABLE `rss`.`Subscription` ADD INDEX `fk_Subscription_Feed_idx` (`FeedId` ASC);
ALTER TABLE `rss`.`Subscription` 
  ADD CONSTRAINT `fk_Subscription_FeedId`
  FOREIGN KEY (`FeedId`)
  REFERENCES `rss`.`Feed` (`Id`)
  ON DELETE NO ACTION
  ON UPDATE NO ACTION;

