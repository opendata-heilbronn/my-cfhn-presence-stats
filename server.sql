CREATE TABLE `presences` (
  `username` varchar(32) NOT NULL,
  `datetime` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `presences` ADD UNIQUE KEY `user_presence` (`username`,`datetime`);
ALTER TABLE `presences` ADD KEY `username` (`username`) USING BTREE;
ALTER TABLE `presences` ADD KEY `datetime` (`datetime`) USING BTREE;
