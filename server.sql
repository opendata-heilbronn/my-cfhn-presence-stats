CREATE TABLE `presences` (
  `username` varchar(32) NOT NULL,
  `datetime` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `presences` ADD UNIQUE KEY `user_presence` (`username`,`datetime`);
ALTER TABLE `presences` ADD KEY `username` (`username`) USING BTREE;
ALTER TABLE `presences` ADD KEY `datetime` (`datetime`) USING BTREE;

CREATE TABLE `streaks` (
  `username` varchar(32) NOT NULL,
  `arrival` datetime NOT NULL,
  `departure` datetime NOT NULL,
  `ticks` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `streaks` ADD UNIQUE KEY `user_streak` (`username`,`arrival`);
ALTER TABLE `streaks` ADD KEY `username` (`username`) USING BTREE;
ALTER TABLE `streaks` ADD KEY `departure` (`departure`) USING BTREE;
ALTER TABLE `streaks` ADD KEY `ticks` (`ticks`) USING BTREE;
