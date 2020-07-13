SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+03:00";

-- --------------------------------------------------------

--
-- Table structure for table `messages`
--

CREATE TABLE `messages`
(
    `update_id`  int(10) UNSIGNED    NOT NULL,
    `reply_to`   int(10) UNSIGNED                  DEFAULT NULL,
    `text`       text COLLATE utf8mb4_bin,
    `sender`     bigint(20) UNSIGNED               DEFAULT NULL,
    `receiver`   bigint(20) UNSIGNED               DEFAULT NULL,
    `seen`       tinyint(1) UNSIGNED NOT NULL      DEFAULT '0',
    `photo`      text COLLATE utf8mb4_bin,
    `voice`      text COLLATE utf8mb4_bin,
    `audio`      text COLLATE utf8mb4_bin,
    `sticker`    text COLLATE utf8mb4_bin,
    `caption`    varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL,
    `message_id` int(11)                           DEFAULT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

-- --------------------------------------------------------

--
-- Table structure for table `suspends`
--

CREATE TABLE `suspends`
(
    `owner_id`     bigint(20) DEFAULT NULL,
    `blocked_user` bigint(20) DEFAULT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users`
(
    `id`            bigint(20) UNSIGNED NOT NULL,
    `first_name`    varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
    `last_name`     varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
    `user_name`     varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
    `language_code` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
    `is_bot`        tinyint(1)                       DEFAULT NULL,
    `link`          varchar(32) COLLATE utf8mb4_bin  DEFAULT NULL,
    `blocked_bot`   tinyint(1)                       DEFAULT NULL,
    `evil_mode`     tinyint(1) UNSIGNED NOT NULL     DEFAULT '0',
    `created_at`    datetime                         DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `messages`
--
ALTER TABLE `messages`
    ADD PRIMARY KEY (`update_id`);

--
-- Indexes for table `suspends`
--
ALTER TABLE `suspends`
    ADD UNIQUE KEY `owner_id` (`owner_id`, `blocked_user`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
    ADD PRIMARY KEY (`id`),
    ADD UNIQUE KEY `link` (`link`);

--
-- AUTO_INCREMENT for dumped tables
--