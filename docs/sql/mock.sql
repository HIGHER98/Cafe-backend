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

INSERT INTO items (name, description, price, upload_date, tag, category) VALUES ("Coffee", "Yummy coffee", "2.50", "2020-09-14", 2, 8),("Flat white", "Yummy flat white", "3.00", "2020-09-16", 2, 8),("Water", "Yummy water", "1.00", "2020-09-01", 1, 8),("Pancakes", "Tasty pancakes", "6.00", "2020-09-26", 8, 1);
-- Light snacks
INSERT INTO items (name, description, price, tag, category) VALUES ("Scones", "With jam & cream", "2.20",7, 3), ("Savory Croissants", "", "4.50", 7, 3), ("Sweet Croissants", "", "2.5", 7, 3),("Toasted Bagels", "", "3.25", 7, 3),("Fresh Fruit Salad", "", "4.5", 7, 3),("Homemade Sausage Rolls", "Pork, Sage & apple sausage rolls served with relish", "3.0", 7, 3),("Homemade Soup of the Day", "", "4.5", 7, 2),("Baked Potatoes", "All served with side salad and coleslaw", "7.5", 7, 2),("Pasta with Checken Strips", "Serve with tomato sauce", "7", 7, 2),("Chilli-Cheese Nachos", "Tortilla chips topped with chilli beef and melted cheese. Served with sour cream, guacamole & salsa salad", "8.5", 7, 2),("Gourmet Sandwiches", "On crusty bread, brown soda bread or gluten free bread", "6.5", 7, 2),("Toasted Sandwiches", "On fresh white or brown bread", "4.95", 4, 2),("Vegetarian Panini", "Roasted red pepper, butternut squash, fresh and sundried tomato, bell peppers and goats cheese drizzled with basil pesto", "7.5", 5, 2),("Hawaiian Panini", "Home baked ham, red onion, pineapple and grated white cheddar", "7.5", 5, 2),("Turkey & Stuffing Panini", "Sliced turkey, baked ham and herb stuffing with mayo and country relish", "7.5", 5, 2),("Sweet Chilli Chicken Panini", "Breast of chicken dressed in sweet chilli mayo, sundried tomato, bell peppers, red onion and white grated cheese", "7.5", 5, 2),("Spicy Chicken Wrap", "Breast of chicken, mixed peppers, tomato, red onion, mixed leaves dressed in a sweet chilli mayo", "7.5", 6, 2),("Vegetarian Wrap", "Mixed leaves, fresh and sundried tomato, red onion, butternut squash, roasted red peppers, bell peppers & goats cheese drizzled with basil pesto", "7.5", 6,2),("Surf Wrap", "Tuna and mayo, mixed leaves, diced red onion, sweetcorn and grated white cheddar", "7.5", 6, 2),("Garlic Checken Wrap", "Breast of chicken, butter head lettuce, fresh and sundried tomato all smothered in garlic mayo", "7.5", 6, 2),("Hearty Italian Style Ciabatta", "Toasted ciabatta with basil pesto, tomato salsa, peppered salami, sundried tomato and white cheddar", "7.95", 7, 2),("C&B Ciabatta", "Toasted ciabatta with roasted chicken strips, crispy streaky bacon, sweet chilli mayo, sundried tomato, mixed leaves and white cheddar drizzled with our homemade basil pesto", "7.95", 7, 2),("Cairdes Caesar Salad", "A mix of Romaine and iceberg lettuce, roast chicken strips, crispy streaky bacon, herb croutons and chopped sundried tomatoes coated in our Caesar dressing and topped with parmesan shavings", "9", 8, 5),("Warm Goats Cheese Salad", "Mixed leaves, fresh and sundried tomato, red onion, roasted red peppers, and carmelised butternut squash drizzled in basil pesto and topped with warm goats cheese on a puff pastry disc", "9", 8,5)

INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (2, "Coconut milk", "0.10", "2020-09-17");
INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (2, "Soy milk", "0.20", "2020-09-17");
INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (4, "Syrup pancakes", "0", "2020-09-26");
INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (4, "Blueberry pancakes", "0.5", "2020-09-26");
INSERT INTO item_options (item_id, opt, add_price, upload_date) VALUES (4, "Bacon pancakes", "0.5", "2020-09-26");

INSERT INTO item_sizes (item_id, item_size, add_price, upload_date) VALUES (4, "2 Pancakes", "0", "2020-09-17");
INSERT INTO item_sizes (item_id, item_size, add_price, upload_date) VALUES (4, "4 Pancakes", "1", "2020-09-17");
INSERT INTO item_sizes (item_id, item_size, add_price, upload_date) VALUES (4, "6 Pancakes", "2", "2020-09-17");

INSERT INTO purchases (email, cust_name, date_time, collection_time) VALUES ("garhyer@mail.com","Gary H" , "2020-09-14 00:00:00", "2020-09-20 16:04:00");
INSERT INTO purchases (email, cust_name, date_time, collection_time) VALUES ("jimihendrix@mail.com","Jimi Hendrix" , "2020-09-14 00:00:00", "2020-09-21 16:04:00");
INSERT INTO purchases (email, cust_name, date_time, collection_time, notes) VALUES ("richardnixon@mail.com","Richard Nixon" , "2020-09-14 00:00:00", "2020-09-21 17:00:00", "gluten free bread please");

INSERT INTO purchase_items (purchase_id, item_id) VALUES (1, 2);
INSERT INTO purchase_items (purchase_id, item_id, item_options_id) VALUES (1, 2, 1);
INSERT INTO purchase_items (purchase_id, item_id) VALUES (2, 3);
INSERT INTO purchase_items (purchase_id, item_id, item_options_id) VALUES (2, 2, 2);

-- 6 Blueberry pancakes
INSERT INTO purchase_items (purchase_id, item_id, item_options_id, item_size_id) VALUES (3, 4, 4, 3);
INSERT INTO purchase_items (purchase_id, item_id) VALUES (3, 3);
