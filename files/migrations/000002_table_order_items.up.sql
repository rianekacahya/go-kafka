BEGIN;

CREATE TABLE order_items (
    id int(10) auto_increment NOT NULL,
	purchase_code varchar(100) NULL,
	sku varchar(100) NULL,
	quantity int(10) NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

COMMIT;
