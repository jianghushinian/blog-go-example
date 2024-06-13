CREATE TABLE `user`
(
    `id`        BIGINT       NOT NULL AUTO_INCREMENT,
    `email`     VARCHAR(255),
    `nickname`  VARCHAR(255),
    `username`  VARCHAR(255) NOT NULL,
    `password`  VARCHAR(255) NOT NULL,
    `createdAt` DATETIME,
    `updatedAt` DATETIME,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
