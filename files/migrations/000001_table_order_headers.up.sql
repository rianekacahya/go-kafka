BEGIN;

CREATE TABLE order_headers (
    id int(10) auto_increment NOT NULL,
    customer_email varchar(100) NULL,
    purchase_code varchar(100) NULL,
    purchase_date datetime NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

COMMIT;
