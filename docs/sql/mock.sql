--Sample data
INSERT INTO item (name, description, price, upload_date) VALUES ("Coffee", "Yummy coffee", "2.50", "2020-09-14");
INSERT INTO item (name, description, price, upload_date) VALUES ("Flat white", "Yummy flat white", "3.00", "2020-09-16");
INSERT INTO item (name, description, price, upload_date) VALUES ("Water", "Yummy water", "1.00", "2020-09-01");

INSERT INTO purchase (item_id, email, cust_name, date_time, collection_time) VALUES (1, "garhyer@mail.com","Gary H" , "2020-09-14 00:00:00", "2020-09-20 16:04:00")
INSERT INTO purchase (item_id, email, cust_name, date_time, collection_time) VALUES (1, "jimihendrix@mail.com","Jimi Hendrix" , "2020-09-14 00:00:00", "2020-09-21 16:04:00")
INSERT INTO purchase (item_id, email, cust_name, date_time, collection_time, notes) VALUES (1, "richardnixon@mail.com","Richard Nixon" , "2020-09-14 00:00:00", "2020-09-21 17:00:00", "gluten free bread please")

