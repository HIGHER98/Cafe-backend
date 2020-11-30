SET FOREIGN_KEY_CHECKS=0;

-- Allergen table
-- EG. Gluten, crustaceans, eggs, etc.
DROP TABLE IF EXISTS `allergens`;
CREATE TABLE `allergens` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`name` varchar(32) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Item category table
-- EG. Breakfast, lunch, etc
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`name` varchar(32) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Tag table
-- EG. Toasted sandwiches, Wraps, etc
DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`name` varchar(32) NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Items table
DROP TABLE IF EXISTS `items`;
CREATE TABLE `items` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT, 
	`name` varchar(32) NOT NULL,
	`description` text,
	`price` float(8, 2) NOT NULL,
	`tag` int(10) unsigned,
	`category` int(10) unsigned NOT NULL,
	`upload_date` DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`is_del` tinyint(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`),
	FOREIGN KEY (`tag`) REFERENCES tags (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Item options table
-- Lists options associated with a menu item. EG. Blueberry/Maple syrup/Bacon pancakes
DROP TABLE IF EXISTS `item_options`;
CREATE TABLE `item_options` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`item_id` int(10) unsigned NOT NULL,
	`opt` varchar(64) NOT NULL,
	`add_price` float(8, 2) NOT NULL,
	`description` text,
	`upload_date` DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`is_del` tinyint(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`),
	FOREIGN KEY (`item_id`) REFERENCES items (`id`)
) ENGINE=InnoDB DEFAULT  CHARSET=utf8;

-- Item size table
-- Eg Large coffee
DROP TABLE IF EXISTS `item_sizes`;
CREATE TABLE `item_sizes` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`item_id` int(10) unsigned NOT NULL,
	`item_size` varchar(64) NOT NULL,
	`add_price` float(8, 2) NOT NULL,
	`description` text,
	`upload_date` DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`is_del` tinyint(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`),
	FOREIGN KEY (`item_id`) REFERENCES items (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Item to allergen table
DROP TABLE IF EXISTS `item_to_allergens`;
CREATE TABLE `item_to_allergens` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`allergen_id` int(10) unsigned NOT NULL,
	`item_id` int(10) unsigned NOT NULL,
	PRIMARY KEY (`id`),
	FOREIGN KEY (`allergen_id`) REFERENCES allergens (`id`),
	FOREIGN KEY (`item_id`) REFERENCES items (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Image table
DROP TABLE IF EXISTS `image`;
CREATE TABLE `image` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`item_id` int(10) unsigned NOT NULL,
	`s3_link` varchar(100) NOT NULL,
	`is_del` tinyint(1) NOT NULL DEFAULT 0, 
	PRIMARY KEY (`id`),
	FOREIGN KEY (`item_id`) REFERENCES items (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Status table
DROP TABLE IF EXISTS `status`;
CREATE TABLE `status` (
	`id` tinyint(3) NOT NULL AUTO_INCREMENT,
	`description` VARCHAR(64),
	`is_del` tinyint(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Staff table
DROP TABLE IF EXISTS `staffs`;
CREATE TABLE `staffs` (
	`id` int(10) NOT NULL AUTO_INCREMENT,
	`name` varchar(64) NOT NULL,
	`is_del` tinyint(1) NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Purchases table
DROP TABLE IF EXISTS `purchases`;
CREATE TABLE `purchases` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`email` varchar(64) NOT NULL,
	`status` tinyint(3) NOT NULL DEFAULT 0,
	`cust_name` varchar(64) NOT NULL,
	`date_time` DATETIME NOT NULL, 
	`collection_time` DATETIME NOT NULL,
	`notes` varchar(256),
	`uuid` varchar(1) NOT NULL,
	PRIMARY KEY (`id`),
	FOREIGN KEY (`status`) REFERENCES status (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `purchase_activities`;
CREATE TABLE `purchase_activities` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`purchase_id` int(10) unsigned NOT NULL,
	`status_set` tinyint(3) NOT NULL,
	`set_by` int(10) NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`),
	FOREIGN KEY (`status_set`) REFERENCES status (`id`),
	FOREIGN KEY (`set_by`) REFERENCES staffs (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Purchase items table
-- Lists all items associated with an order
DROP TABLE IF EXISTS `purchase_items`;
CREATE TABLE `purchase_items` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`purchase_id` int(10) unsigned NOT NULL,
	`item_id` int(10) unsigned NOT NULL,
	`opt_id` int(10) unsigned, 
	`size_id` int(10) unsigned,
	PRIMARY KEY (`id`),
	FOREIGN KEY (`purchase_id`) REFERENCES purchases (`id`),
	FOREIGN KEY (`item_id`) REFERENCES items (`id`), 
	FOREIGN KEY (`opt_id`) REFERENCES item_options (`id`), 
	FOREIGN KEY (`size_id`) REFERENCES item_sizes (`id`) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Roles table
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

-- Menu view
CREATE VIEW item_views AS
	SELECT 
		items.id, items.name AS item_name, items.description, items.price,  items.category AS category_id, items.tag AS tag_id,
		item_options.id AS opt_id, item_options.opt, item_options.add_price AS option_price,  item_options.description AS option_description,
		item_sizes.id AS size_id, item_sizes.item_size, item_sizes.add_price AS size_price, item_sizes.description AS size_description,
		category.name AS category,
		tags.name AS tag
	FROM 
		items
	LEFT JOIN 
		item_options ON items.id = item_options.item_id 
	LEFT JOIN 
		item_sizes ON items.id = item_sizes.item_id
	LEFT JOIN
		category ON items.category = category.id
	LEFT JOIN 
		tags ON items.tag = tags.id
	WHERE items.is_del=0 AND (item_sizes.is_del=0 OR item_sizes.is_del IS NULL) AND (item_options.is_del=0 OR item_options.is_del IS NULL);

-- Queue view
CREATE VIEW queue_views AS
	SELECT p.id, i.name, s.description, i.price, p.collection_time, p.cust_name, p.notes, p.date_time AS order_time
	FROM purchase AS p, item AS i, status AS s
	WHERE p.item_id = i.id AND p.status=s.id;

-- Purchase views
CREATE VIEW purchase_views AS
SELECT 
	p.id AS purchases_id, p.cust_name, p.email, p.date_time, p.collection_time, p.notes,
	item_views.item_name, item_views.opt, item_views.item_size, (item_views.price + IFNULL(item_views.option_price, 0) + IFNULL(item_views.size_price, 0)) AS cost,
	pi.id AS purchase_items_id,
	s.description AS status
FROM 
	purchases AS p, 
	status AS s,
	purchase_items AS pi
LEFT JOIN
	item_views ON (
		pi.item_id = item_views.id AND
		(pi.opt_id = item_views.opt_id OR ISNULL(pi.opt_id)) AND
		(pi.size_id = item_views.size_id OR ISNULL(pi.size_id))
	)
WHERE 
	p.id = pi.purchase_id AND 
	pi.item_id = item_views.id AND
	p.status = s.id
ORDER BY p.id ASC;



INSERT INTO status (description) VALUES ('Pending transaction');
INSERT INTO status (description) VALUES ('Pending');
INSERT INTO status (description) VALUES ('Confirmed');
INSERT INTO status (description) VALUES ('Collected');

INSERT INTO roles (title) VALUES ('Admin');
INSERT INTO roles (title) VALUES ('Staff');

INSERT INTO allergen (name) VALUES ('Gluten'), ('Crustaceans'), ('Eggs'), ('Fish'), ('Peanuts'), ('Soybeans'), ('Milk'), ('Tree Nuts'), ('Celery'), ('Mustard'), ('Sesame Seeds'), ('Sulphite'), ('Lupin'), ('Molluscs'), ('Coeliac'), ('Coeliac available');

SET FOREIGN_KEY_CHECKS=1;


