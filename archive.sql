-- MySQL dump 10.13  Distrib 5.7.22, for Linux (x86_64)
--
-- Host: localhost    Database: archive
-- ------------------------------------------------------
-- Server version	5.7.22-0ubuntu0.16.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Archive`
--

DROP TABLE IF EXISTS `Archive`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Archive` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `objectInstanceIdentifier` bigint(20) unsigned DEFAULT NULL,
  `element` blob,
  `area` smallint(6) DEFAULT NULL,
  `service` smallint(6) DEFAULT NULL,
  `version` tinyint(4) DEFAULT NULL,
  `number` smallint(6) DEFAULT NULL,
  `domain` text,
  `timestamp` datetime DEFAULT NULL,
  `details.related` bigint(20) DEFAULT NULL,
  `network` text,
  `provider` text,
  `details.source` blob,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Archive`
--

LOCK TABLES `Archive` WRITE;
/*!40000 ALTER TABLE `Archive` DISABLE KEYS */;
INSERT INTO `Archive` VALUES (1,6162062580224294496,'\0\0\0\0øT\»/',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0\0'),(2,1,'\0\0\0\0æ∑≥\·',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(3,2,'\0\0\0\0Ω\÷/\Ï',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(4,3,'\0\0\0\0>˛z',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(5,4,'\0\0\0\0?7\Í@',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(6,5,'\0\0\0\0?}c\ƒ',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(7,6,'\0\0\0\0øz€º',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(8,7,'\0\0\0\0?!k{',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(9,8,'\0\0\0\0æ\Ê®',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(10,9,'\0\0\0\0æZ˚â',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0	'),(11,10,'\0\0\0\0?u\“˝',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0\n'),(12,11,'\0\0\0\0øv1˘',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(13,12,'\0\0\0\0>µ\Î\Á',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(14,13,'\0\0\0\0?t\·8',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0\r'),(15,14,'\0\0\0\0æ\"(I',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(16,15,'\0\0\0\0?Dh4',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(17,16,'\0\0\0\0Ω\«y≤',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(18,17,'\0\0\0\0æÖQç',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(19,18,'\0\0\0\0øã',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(20,19,'\0\0\0\0øs¶\‹',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(21,20,'\0\0\0\0?<\ \≈',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(22,21,'\0\0\0\0ºç#',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(23,22,'\0\0\0\0>ëä',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(24,23,'\0\0\0\0>?Òe',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(25,24,'\0\0\0\0æ\‚MN',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(26,25,'\0\0\0\0>¨˛',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(27,26,'\0\0\0\0æÄ|ä',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0\Z'),(28,27,'\0\0\0\0?~D\»',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(29,28,'\0\0\0\0>ﬂá\Ê',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(30,29,'\0\0\0\0Ω$\·\‹',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(31,30,'\0\0\0\0øXî',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(32,31,'\0\0\0\0>\Ì˛|',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0'),(33,32,'\0\0\0\0>∞Xˆ',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0 '),(34,33,'\0\0\0\0æé',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0!'),(35,34,'\0\0\0\0?\Zæv',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0\"'),(36,35,'\0\0\0\0ø1r',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0#'),(37,36,'\0\0\0\0ødtà',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0$'),(38,37,'\0\0\0\0?.úX',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network1','tests/provider1','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0%'),(39,38,'\0\0\0\0Ω\Œ',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0&'),(40,39,'\0\0\0\0?V%¢',2,3,1,1,'fr.cnes.archiveservice.test','2018-05-02 17:45:40',0,'tests/network2','tests/provider2','\0\0\0\0\0\0\0\0\0fr\0\0\0cnes\0\0\0archiveservice\0\0\0test\0\0\0\0\0\0\0\'');
/*!40000 ALTER TABLE `Archive` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-05-02 19:46:39
