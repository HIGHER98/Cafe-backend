SET FOREIGN_KEY_CHECKS=0;
-- Items table
DROP TABLE IF EXISTS `item`;
CREATE TABLE `item` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT, 
	`name` varchar(32) NOT NULL,
	`description` text,
	`price` float(8, 2) NOT NULL,
	`upload_date` DATE NOT NULL,
	`is_del` tinyint(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Image table
DROP TABLE IF EXISTS `image`;
CREATE TABLE `image` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`item_id` int(10) unsigned NOT NULL,
	`s3_link` varchar(100) NOT NULL,
	`is_del` tinyint(1) NOT NULL DEFAULT 0, 
	PRIMARY KEY (`id`),
	FOREIGN KEY (`item_id`) REFERENCES item (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Status table
DROP TABLE IF EXISTS `status`;
CREATE TABLE `status` (
	`id` tinyint(3) NOT NULL AUTO_INCREMENT,
	`description` VARCHAR(64),
	`is_del` tinyint(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- Purchases table
DROP TABLE IF EXISTS `purchase`;
CREATE TABLE `purchase` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`item_id` int(10) unsigned NOT NULL,
	`email` varchar(64) NOT NULL,
	`status` tinyint(3) NOT NULL,
	`cust_name` varchar(64) NOT NULL,
	`date_time` DATETIME NOT NULL, 
	`collection_time` DATETIME NOT NULL,
	`notes` varchar(256),
	PRIMARY KEY (`id`),
	FOREIGN KEY (`item_id`) REFERENCES item (`id`),
	FOREIGN KEY (`status`) REFERENCES status (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
	`id` tinyint(3) unsigned NOT NULL AUTO_INCREMENT,
	`title` varchar(30) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Users table
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`username` varchar(50) NOT NULL DEFAULT '',
	`password` varchar(256) NOT NULL DEFAULT '',
	`role` tinyint(3) unsigned NOT NULL,
	`is_del` tinyint(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`),
	FOREIGN KEY (`role`) REFERENCES roles (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- Views

-- User view
CREATE VIEW user_view AS
	SELECT u.id, u.username, r.title
	FROM users AS u, roles AS r
	WHERE u.role = r.id AND u.is_del=0;

-- Queue view
CREATE VIEW queue_view AS
	SELECT p.id, i.name, s.description, i.price, p.collection_time, p.cust_name, p.notes, p.date_time AS order_time
	FROM purchase AS p, item AS i, status AS s
	WHERE p.item_id = i.id AND p.status=s.id;

INSERT INTO status (description) VALUES ('Pending transaction');
INSERT INTO status (description) VALUES ('Pending');
INSERT INTO status (description) VALUES ('Confirmed');
INSERT INTO status (description) VALUES ('Collection');

INSERT INTO roles (title) VALUES ('Admin');
INSERT INTO roles (title) VALUES ('Staff');


SET FOREIGN_KEY_CHECKS=1;


