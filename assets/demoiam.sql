# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.5.52-0+deb8u1)
# Datenbank: demoiam
# Erstellt am: 2017-06-20 17:03:37 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Export von Tabelle clients
# ------------------------------------------------------------

DROP TABLE IF EXISTS `clients`;

CREATE TABLE `clients` (
  `id` varchar(36) NOT NULL DEFAULT '',
  `secret` varchar(255) NOT NULL DEFAULT '',
  `created` datetime NOT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Export von Tabelle clients_domains
# ------------------------------------------------------------

DROP TABLE IF EXISTS `clients_domains`;

CREATE TABLE `clients_domains` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `client_id` varchar(36) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `dom_cli` (`client_id`),
  CONSTRAINT `dom_cli` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Export von Tabelle clients_scopes
# ------------------------------------------------------------

DROP TABLE IF EXISTS `clients_scopes`;

CREATE TABLE `clients_scopes` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `client_id` varchar(36) NOT NULL DEFAULT '',
  `scope` varchar(30) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `cli_sco` (`client_id`),
  CONSTRAINT `cli_sco` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Export von Tabelle refresh_tokens
# ------------------------------------------------------------

DROP TABLE IF EXISTS `refresh_tokens`;

CREATE TABLE `refresh_tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(60) NOT NULL DEFAULT '',
  `token` text NOT NULL,
  `expiry` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `refcodeq` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
