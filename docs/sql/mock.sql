--Sample data

INSERT INTO category (name) VALUES ("Breakfast");
INSERT INTO category (name) VALUES ("Lunch");
INSERT INTO category (name) VALUES ("Light Snacks");
INSERT INTO category (name) VALUES ("Side Order");
INSERT INTO category (name) VALUES ("Salads");
INSERT INTO category (name) VALUES ("Desserts");
INSERT INTO category (name) VALUES ("Kids Menu");
INSERT INTO category (name) VALUES ("Drinks");

INSERT INTO tags (name) VALUES ("Cold drink");
INSERT INTO tags (name) VALUES ("Hot drink");
INSERT INTO tags (name) VALUES ("Soup");
INSERT INTO tags (name) VALUES ("Toasted Sandwiches");
INSERT INTO tags (name) VALUES ("Paninis");
INSERT INTO tags (name) VALUES ("Wraps");
INSERT INTO tags (name) VALUES ("Ciabatta");
INSERT INTO tags (name) VALUES ("Other");

INSERT INTO items (name, description, price, upload_date, tag, category) VALUES ("Coffee", "Yummy coffee", "2.50", "2020-09-14", 2, 8);
INSERT INTO items (name, description, price, upload_date, tag, category) VALUES ("Flat white", "Yummy flat white", "3.00", "2020-09-16", 2, 8);
INSERT INTO items (name, description, price, upload_date, tag, category) VALUES ("Water", "Yummy water", "1.00", "2020-09-01", 1, 8);
INSERT INTO items (name, description, price, upload_date, tag, category) VALUES ("Pancakes", "Tasty pancakes", "6.00", "2020-09-26", 8, 1);

INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (2, "Coconut milk", "0.10", "2020-09-17");
INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (2, "Soy milk", "0.20", "2020-09-17");
INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (4, "Syrup pancakes", "0", "2020-09-26");
INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (4, "Blueberry pancakes", "0.5", "2020-09-26");
INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (4, "Bacon pancakes", "0.5", "2020-09-26");

INSERT INTO item_sizes (item_id, item_size, add_price, upload_date) VALUES (4, "2 Pancakes", "0", "2020-09-17");
INSERT INTO item_sizes (item_id, item_size, add_price, upload_date) VALUES (4, "4 Pancakes", "1", "2020-09-17");
INSERT INTO item_sizes (item_id, item_size, add_price, upload_date) VALUES (4, "6 Pancakes", "2", "2020-09-17");

INSERT INTO purchases (email, cust_name, date_time, collection_time, status) VALUES ("garhyer@mail.com","Gary H" , "2020-09-14 00:00:00", "2020-09-20 16:04:00", 1);
INSERT INTO purchases (email, cust_name, date_time, collection_time, status) VALUES ("jimihendrix@mail.com","Jimi Hendrix" , "2020-09-14 00:00:00", "2020-09-21 16:04:00", 2);
INSERT INTO purchases (email, cust_name, date_time, collection_time, notes, status) VALUES ("richardnixon@mail.com","Richard Nixon" , "2020-09-14 00:00:00", "2020-09-21 17:00:00", "gluten free bread please", 3);

INSERT INTO purchase_items (purchase_id, item_id) VALUES (1, 2);
INSERT INTO purchase_items (purchase_id, item_id, opt_id) VALUES (1, 2, 1);
INSERT INTO purchase_items (purchase_id, item_id) VALUES (2, 3);
INSERT INTO purchase_items (purchase_id, item_id, opt_id) VALUES (2, 2, 2);

-- 6 Blueberry pancakes
INSERT INTO purchase_items (purchase_id, item_id, opt_id, size_id) VALUES (3, 4, 4, 3);
INSERT INTO purchase_items (purchase_id, item_id) VALUES (3, 3);
