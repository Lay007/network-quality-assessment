-- phpMyAdmin SQL Dump
-- version 4.9.2
-- https://www.phpmyadmin.net/
--
-- Хост: 127.0.0.1
-- Время создания: Июн 02 2020 г., 20:43
-- Версия сервера: 10.4.10-MariaDB
-- Версия PHP: 7.3.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- База данных: `server_sfp_sla`
--

-- --------------------------------------------------------

--
-- Структура таблицы `global_config`
--

CREATE TABLE `global_config` (
  `server_IP` varchar(20) NOT NULL,
  `net_interface_name` varchar(50) DEFAULT NULL,
  `zabbix_server_name` varchar(255) DEFAULT NULL,
  `zabbix_port` int(11) DEFAULT NULL,
  `VLAN` tinyint(1) DEFAULT NULL,
  `VLAN_number` int(8) DEFAULT NULL,
  `QinQ` tinyint(1) DEFAULT NULL,
  `QinQ_number` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `global_config`
--

INSERT INTO `global_config` (`server_IP`, `net_interface_name`, `zabbix_server_name`, `zabbix_port`, `VLAN`, `VLAN_number`, `QinQ`, `QinQ_number`) VALUES
('10.0.10.115', '0', 'localhost', 10051, 0, 100, 0, 0);

-- --------------------------------------------------------

--
-- Структура таблицы `modules_sfp_sla`
--

CREATE TABLE `modules_sfp_sla` (
  `id` int(11) NOT NULL,
  `mac` bigint(20) NOT NULL,
  `name` varchar(50) NOT NULL,
  `address_ip` varchar(20) NOT NULL,
  `version` varchar(10) NOT NULL,
  `zabbix_node_name` varchar(255) DEFAULT NULL,
  `location` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `modules_sfp_sla_load_rez`
--

CREATE TABLE `modules_sfp_sla_load_rez` (
  `id` int(11) NOT NULL,
  `module_id` int(11) NOT NULL,
  `datatime` datetime NOT NULL,
  `load_to_lazer` int(11) NOT NULL,
  `load_to_com` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `net_interfaces_from_server_sla`
--

CREATE TABLE `net_interfaces_from_server_sla` (
  `id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `address_IP` varchar(50) NOT NULL,
  `address_mac` varchar(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `test_bert`
--

CREATE TABLE `test_bert` (
  `id` int(11) NOT NULL,
  `test_type` int(11) NOT NULL,
  `module_first` int(11) NOT NULL,
  `module_second` int(11) NOT NULL,
  `thr_begin` int(11) NOT NULL,
  `count_prob_packs` int(11) NOT NULL,
  `count_probs` int(11) NOT NULL,
  `rez_64` int(11) DEFAULT NULL,
  `rez_128` int(11) DEFAULT NULL,
  `rez_256` int(11) DEFAULT NULL,
  `rez_512` int(11) DEFAULT NULL,
  `rez_1024` int(11) DEFAULT NULL,
  `rez_1280` int(11) DEFAULT NULL,
  `rez_1518` int(11) DEFAULT NULL,
  `rez_4096` int(11) DEFAULT NULL,
  `rez_9000` int(11) DEFAULT NULL,
  `datetime_start` datetime DEFAULT NULL,
  `datetime_end` datetime DEFAULT NULL,
  `status` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



--
-- Структура таблицы `test_frame_loss`
--

CREATE TABLE `test_frame_loss` (
  `id` int(11) NOT NULL,
  `test_type` int(11) NOT NULL,
  `module_first` int(11) NOT NULL,
  `module_second` int(11) NOT NULL,
  `thr_begin` int(11) NOT NULL,
  `step` int(11) NOT NULL,
  `count_frames` int(11) NOT NULL,
  `count_steps` int(11) NOT NULL,
  `datetime_start` datetime DEFAULT NULL,
  `datetime_end` datetime DEFAULT NULL,
  `status` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


--
-- Структура таблицы `test_frame_loss_rez`
--

CREATE TABLE `test_frame_loss_rez` (
  `id` int(11) NOT NULL,
  `id_test` int(11) NOT NULL,
  `step_number` int(11) DEFAULT NULL,
  `rez_64` float DEFAULT NULL,
  `rez_128` float DEFAULT NULL,
  `rez_256` float DEFAULT NULL,
  `rez_512` float DEFAULT NULL,
  `rez_1024` float DEFAULT NULL,
  `rez_1280` float DEFAULT NULL,
  `rez_1518` float DEFAULT NULL,
  `rez_4096` float DEFAULT NULL,
  `rez_9000` float DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `test_latency`
--

CREATE TABLE `test_latency` (
  `id` int(11) NOT NULL,
  `test_type` int(11) DEFAULT NULL,
  `module_first` int(11) DEFAULT NULL,
  `module_second` int(11) NOT NULL,
  `thr_begin` int(11) NOT NULL,
  `count_packs` int(11) NOT NULL,
  `count_tests` int(11) NOT NULL,
  `rez_64` float DEFAULT NULL,
  `rez_64_max` float DEFAULT NULL,
  `rez_64_min` float DEFAULT NULL,
  `rez_128` float DEFAULT NULL,
  `rez_128_max` float DEFAULT NULL,
  `rez_128_min` float DEFAULT NULL,
  `rez_256` float DEFAULT NULL,
  `rez_256_max` float DEFAULT NULL,
  `rez_256_min` float DEFAULT NULL,
  `rez_512` float DEFAULT NULL,
  `rez_512_max` float DEFAULT NULL,
  `rez_512_min` float DEFAULT NULL,
  `rez_1024` float DEFAULT NULL,
  `rez_1024_max` float DEFAULT NULL,
  `rez_1024_min` float DEFAULT NULL,
  `rez_1280` float DEFAULT NULL,
  `rez_1280_max` float DEFAULT NULL,
  `rez_1280_min` float DEFAULT NULL,
  `rez_1518` float DEFAULT NULL,
  `rez_1518_max` float DEFAULT NULL,
  `rez_1518_min` float DEFAULT NULL,
  `datetime_start` datetime DEFAULT NULL,
  `datetime_end` datetime DEFAULT NULL,
  `status` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


--
-- Структура таблицы `test_sla_real`
--

CREATE TABLE `test_sla_real` (
  `id` int(11) NOT NULL,
  `test_type` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `module_first` int(11) NOT NULL,
  `module_second` int(11) NOT NULL,
  `block_size` int(11) NOT NULL,
  `clock` int(11) NOT NULL,
  `count` int(11) NOT NULL,
  `node_zabbix` varchar(50) NOT NULL,
  `test_delay` tinyint(1) DEFAULT NULL,
  `test_delay_jitter` tinyint(1) DEFAULT NULL,
  `test_loss` tinyint(1) DEFAULT NULL,
  `test_delay_1` tinyint(1) DEFAULT NULL,
  `test_delay1_jitter` tinyint(1) DEFAULT NULL,
  `test_load_sfp1_to_laser` tinyint(1) DEFAULT NULL,
  `test_load_sfp1_to_com` tinyint(1) DEFAULT NULL,
  `test_load_sfp2_to_laser` tinyint(1) DEFAULT NULL,
  `test_load_sfp2_to_com` tinyint(1) DEFAULT NULL,
  `data_start` datetime DEFAULT NULL,
  `delay_max` float DEFAULT NULL,
  `jitter_max` float DEFAULT NULL,
  `delay1_max` float DEFAULT NULL,
  `jitter1_max` float DEFAULT NULL,
  `loss_max` float DEFAULT NULL,
  `status` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


--
-- Структура таблицы `test_sla_real_alarm`
--

CREATE TABLE `test_sla_real_alarm` (
  `id` int(11) NOT NULL,
  `id_test` int(11) NOT NULL,
  `datatime` datetime NOT NULL,
  `id_var` int(11) DEFAULT NULL,
  `Value` float DEFAULT NULL,
  `message` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `test_sla_real_rez`
--

CREATE TABLE `test_sla_real_rez` (
  `id_solve` int(11) NOT NULL,
  `datetime` datetime NOT NULL,
  `test_id` int(11) NOT NULL,
  `delay_rez` float DEFAULT NULL,
  `delay_to_rez` float DEFAULT NULL,
  `delay_un_rez` float DEFAULT NULL,
  `jitter_delay_rez` float DEFAULT NULL,
  `jitter_delay_to` float DEFAULT NULL,
  `jitter_delay_un` float DEFAULT NULL,
  `packet_loss` float DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `test_throughput`
--

CREATE TABLE `test_throughput` (
  `id` int(11) NOT NULL,
  `test_type` int(11) DEFAULT NULL,
  `module_first` int(11) NOT NULL,
  `module_second` int(11) NOT NULL,
  `thr_begin` int(11) DEFAULT NULL,
  `count` int(11) DEFAULT NULL,
  `ch_type` int(11) DEFAULT NULL,
  `max_loss` int(11) DEFAULT NULL,
  `rez_64` int(11) DEFAULT NULL,
  `rez_128` int(11) DEFAULT NULL,
  `rez_256` int(11) DEFAULT NULL,
  `rez_512` int(11) DEFAULT NULL,
  `rez_1024` int(11) DEFAULT NULL,
  `rez_1280` int(11) DEFAULT NULL,
  `rez_1518` int(11) DEFAULT NULL,
  `rez_4096` int(11) DEFAULT NULL,
  `rez_9000` int(11) DEFAULT NULL,
  `datetime_start` datetime DEFAULT NULL,
  `datetime_end` datetime DEFAULT NULL,
  `status` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


--
-- Структура таблицы `test_y1564`
--

CREATE TABLE `test_y1564` (
  `id` int(11) NOT NULL,
  `test_type` int(11) NOT NULL,
  `module_first` int(11) NOT NULL,
  `module_second` int(11) NOT NULL,
  `block_size` int(11) NOT NULL,
  `ToS` int(11) DEFAULT NULL,
  `VLAN_priority` int(11) DEFAULT NULL,
  `CIR` int(11) DEFAULT NULL,
  `EIR` int(11) DEFAULT NULL,
  `TP` int(11) DEFAULT NULL,
  `period` int(11) DEFAULT NULL,
  `step_count` int(11) DEFAULT NULL,
  `max_FTD` float DEFAULT NULL,
  `max_FVD` float DEFAULT NULL,
  `max_FLR` float DEFAULT NULL,
  `rez_IR_s1` float DEFAULT NULL,
  `rez_FTD_s1` float DEFAULT NULL,
  `rez_FVD_s1` float DEFAULT NULL,
  `rez_FLR_s1` float DEFAULT NULL,
  `rez_IR_s2` float DEFAULT NULL,
  `rez_FTD_s2` float DEFAULT NULL,
  `rez_FVD_s2` float DEFAULT NULL,
  `rez_FLR_s2` float DEFAULT NULL,
  `rez_IR_s3` float DEFAULT NULL,
  `rez_FTD_s3` float DEFAULT NULL,
  `rez_FVD_s3` float DEFAULT NULL,
  `rez_FLR_s3` float DEFAULT NULL,
  `rez_IR_s4` float DEFAULT NULL,
  `rez_FTD_s4` float DEFAULT NULL,
  `rez_FVD_s4` float DEFAULT NULL,
  `rez_FLR_s4` float DEFAULT NULL,
  `rez_IR_eir` float DEFAULT NULL,
  `rez_FTD_eir` float DEFAULT NULL,
  `rez_FVD_eir` float DEFAULT NULL,
  `rez_FLR_eir` float DEFAULT NULL,
  `rez_IR_tp` float DEFAULT NULL,
  `rez_FTD_tp` float DEFAULT NULL,
  `rez_FVD_tp` float DEFAULT NULL,
  `rez_FLR_tp` float DEFAULT NULL,
  `datetime_start` datetime DEFAULT NULL,
  `datetime_end` datetime DEFAULT NULL,
  `status` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


--
-- Структура таблицы `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `login` varchar(50) NOT NULL,
  `password` varchar(50) NOT NULL,
  `type` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `users`
--

INSERT INTO `users` (`id`, `login`, `password`, `type`) VALUES
(1, 'admin', 'adminfibertrade', 'Администратор');

-- --------------------------------------------------------

--
-- Структура таблицы `zabbix_node`
--

CREATE TABLE `zabbix_node` (
  `id` int(11) NOT NULL,
  `name` varchar(50) NOT NULL,
  `description` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `zabbix_node`
--

INSERT INTO `zabbix_node` (`id`, `name`, `description`) VALUES
(1, 'SFP-SLA_4400', ''),
(2, 'SFP-SLA_4401', '');

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `global_config`
--
ALTER TABLE `global_config`
  ADD PRIMARY KEY (`server_IP`);

--
-- Индексы таблицы `modules_sfp_sla`
--
ALTER TABLE `modules_sfp_sla`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `id` (`id`),
  ADD KEY `id_2` (`id`);

--
-- Индексы таблицы `modules_sfp_sla_load_rez`
--
ALTER TABLE `modules_sfp_sla_load_rez`
  ADD PRIMARY KEY (`id`),
  ADD KEY `module_id` (`module_id`);

--
-- Индексы таблицы `net_interfaces_from_server_sla`
--
ALTER TABLE `net_interfaces_from_server_sla`
  ADD PRIMARY KEY (`id`);

--
-- Индексы таблицы `test_bert`
--
ALTER TABLE `test_bert`
  ADD PRIMARY KEY (`id`),
  ADD KEY `module_first` (`module_first`,`module_second`),
  ADD KEY `module_second` (`module_second`);

--
-- Индексы таблицы `test_frame_loss`
--
ALTER TABLE `test_frame_loss`
  ADD PRIMARY KEY (`id`),
  ADD KEY `module_first` (`module_first`,`module_second`),
  ADD KEY `id` (`id`),
  ADD KEY `id_2` (`id`);

--
-- Индексы таблицы `test_frame_loss_rez`
--
ALTER TABLE `test_frame_loss_rez`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id` (`id`,`id_test`),
  ADD KEY `id_test` (`id_test`);

--
-- Индексы таблицы `test_latency`
--
ALTER TABLE `test_latency`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id` (`id`,`module_first`),
  ADD KEY `module_first` (`module_first`),
  ADD KEY `module_second` (`module_second`);

--
-- Индексы таблицы `test_sla_real`
--
ALTER TABLE `test_sla_real`
  ADD PRIMARY KEY (`id`),
  ADD KEY `module_first` (`module_first`,`module_second`),
  ADD KEY `module_second` (`module_second`);

--
-- Индексы таблицы `test_sla_real_alarm`
--
ALTER TABLE `test_sla_real_alarm`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `id` (`id`),
  ADD KEY `id_test` (`id_test`);

--
-- Индексы таблицы `test_sla_real_rez`
--
ALTER TABLE `test_sla_real_rez`
  ADD PRIMARY KEY (`id_solve`),
  ADD UNIQUE KEY `id_solve` (`id_solve`),
  ADD KEY `test_id` (`test_id`);

--
-- Индексы таблицы `test_throughput`
--
ALTER TABLE `test_throughput`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id` (`id`,`module_first`,`module_second`),
  ADD KEY `module_first` (`module_first`),
  ADD KEY `module_second` (`module_second`);

--
-- Индексы таблицы `test_y1564`
--
ALTER TABLE `test_y1564`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id` (`id`,`module_first`,`module_second`),
  ADD KEY `module_first` (`module_first`),
  ADD KEY `module_second` (`module_second`);

--
-- Индексы таблицы `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id` (`id`);

--
-- Индексы таблицы `zabbix_node`
--
ALTER TABLE `zabbix_node`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `id_3` (`id`),
  ADD KEY `id` (`id`),
  ADD KEY `id_2` (`id`);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `modules_sfp_sla`
--
ALTER TABLE `modules_sfp_sla`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=12;

--
-- AUTO_INCREMENT для таблицы `modules_sfp_sla_load_rez`
--
ALTER TABLE `modules_sfp_sla_load_rez`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `net_interfaces_from_server_sla`
--
ALTER TABLE `net_interfaces_from_server_sla`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `test_bert`
--
ALTER TABLE `test_bert`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT для таблицы `test_frame_loss`
--
ALTER TABLE `test_frame_loss`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT для таблицы `test_frame_loss_rez`
--
ALTER TABLE `test_frame_loss_rez`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `test_latency`
--
ALTER TABLE `test_latency`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT для таблицы `test_sla_real`
--
ALTER TABLE `test_sla_real`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

--
-- AUTO_INCREMENT для таблицы `test_sla_real_alarm`
--
ALTER TABLE `test_sla_real_alarm`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `test_sla_real_rez`
--
ALTER TABLE `test_sla_real_rez`
  MODIFY `id_solve` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT для таблицы `test_throughput`
--
ALTER TABLE `test_throughput`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=22;

--
-- AUTO_INCREMENT для таблицы `test_y1564`
--
ALTER TABLE `test_y1564`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT для таблицы `zabbix_node`
--
ALTER TABLE `zabbix_node`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- Ограничения внешнего ключа сохраненных таблиц
--

--
-- Ограничения внешнего ключа таблицы `modules_sfp_sla_load_rez`
--
-- ALTER TABLE `modules_sfp_sla_load_rez`
--  ADD CONSTRAINT `modules_sfp_sla_load_rez_ibfk_1` FOREIGN KEY (`module_id`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `test_bert`
--
ALTER TABLE `test_bert`
  ADD CONSTRAINT `test_bert_ibfk_1` FOREIGN KEY (`module_first`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `test_bert_ibfk_2` FOREIGN KEY (`module_second`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `test_frame_loss`
--
ALTER TABLE `test_frame_loss`
  ADD CONSTRAINT `test_frame_loss_ibfk_1` FOREIGN KEY (`module_first`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `test_frame_loss_rez`
--
ALTER TABLE `test_frame_loss_rez`
  ADD CONSTRAINT `test_frame_loss_rez_ibfk_1` FOREIGN KEY (`id_test`) REFERENCES `test_frame_loss` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `test_latency`
--
ALTER TABLE `test_latency`
  ADD CONSTRAINT `test_latency_ibfk_1` FOREIGN KEY (`module_first`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `test_latency_ibfk_2` FOREIGN KEY (`module_second`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `test_sla_real`
--
ALTER TABLE `test_sla_real`
  ADD CONSTRAINT `test_sla_real_ibfk_1` FOREIGN KEY (`module_first`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `test_sla_real_ibfk_2` FOREIGN KEY (`module_second`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `test_sla_real_alarm`
--
ALTER TABLE `test_sla_real_alarm`
  ADD CONSTRAINT `test_sla_real_alarm_ibfk_1` FOREIGN KEY (`id_test`) REFERENCES `test_sla_real` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `test_sla_real_rez`
--
ALTER TABLE `test_sla_real_rez`
  ADD CONSTRAINT `test_sla_real_rez_ibfk_1` FOREIGN KEY (`test_id`) REFERENCES `test_sla_real` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `test_throughput`
--
ALTER TABLE `test_throughput`
  ADD CONSTRAINT `test_throughput_ibfk_1` FOREIGN KEY (`module_first`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `test_throughput_ibfk_2` FOREIGN KEY (`module_second`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Ограничения внешнего ключа таблицы `test_y1564`
--
ALTER TABLE `test_y1564`
  ADD CONSTRAINT `test_y1564_ibfk_1` FOREIGN KEY (`module_first`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `test_y1564_ibfk_2` FOREIGN KEY (`module_second`) REFERENCES `modules_sfp_sla` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
