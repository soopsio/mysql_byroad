-- phpMyAdmin SQL Dump
-- version 4.1.9
-- http://www.phpmyadmin.net
--
-- Host: localhost
-- Generation Time: 2016-08-05 15:56:47
-- 服务器版本： 5.6.30-1+deb.sury.org~xenial+2-log
-- PHP Version: 5.6.24-1+deb.sury.org~xenial+1


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- Database: `byroad`
--
CREATE DATABASE IF NOT EXISTS `byroad` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `byroad`;

-- --------------------------------------------------------

--
-- 表的结构 `notify_field`
--

CREATE TABLE IF NOT EXISTS `notify_field_kafka` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `schema` varchar(120) NOT NULL,
  `table` varchar(120) NOT NULL,
  `column` varchar(120) NOT NULL,
  `send` int(11) NOT NULL,
  `task_id` int(11) NOT NULL,
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- 表的结构 `task`
--

CREATE TABLE IF NOT EXISTS `task_kafka` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(120) NOT NULL,
  `apiurl` varchar(255) NOT NULL,
  `event` varchar(120) NOT NULL,
  `stat` varchar(120) NOT NULL,
  `create_time` datetime NOT NULL,
  `create_user` varchar(120) NOT NULL,
  `routine_count` int(11) NOT NULL,
  `re_routine_count` int(11) NOT NULL,
  `re_send_time` int(11) NOT NULL,
  `retry_count` int(11) NOT NULL,
  `timeout` int(11) NOT NULL,
  `desc` varchar(255) DEFAULT NULL,
  `pack_protocal` int(11) DEFAULT NULL,
  `db_instance_name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- 表的结构 `tasklog`
--

CREATE TABLE IF NOT EXISTS `tasklog_kafka` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `task_id` int(11) NOT NULL,
  `message` varchar(1000) DEFAULT NULL,
  `reason` varchar(1000) DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

ALTER TABLE `byroad`.`task_kafka` 
ADD COLUMN `phone_numbers` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '报警电话，使用分号隔开' AFTER `db_instance_name`,
ADD COLUMN `emails` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '报警邮箱，使用分号隔开' AFTER `phone_numbers`,
ADD COLUMN `alert` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否开启报警' AFTER `emails`;
ADD COLUMN `subscribe_stat` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否开启kafka consumer' AFTER `alert` ,
ADD COLUMN `push_stat` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否开启消息推送' AFTER `subscribe_stat` ;