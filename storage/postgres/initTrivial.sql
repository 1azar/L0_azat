-- Таблица для заказов
CREATE TABLE orders (
                        order_uid VARCHAR(255) PRIMARY KEY,
                        track_number VARCHAR(255),
                        entry VARCHAR(255),
                        delivery_info JSONB,
                        payment_info JSONB,
                        locale VARCHAR(10),
                        internal_signature VARCHAR(255),
                        customer_id VARCHAR(255),
                        delivery_service VARCHAR(255),
                        shardkey VARCHAR(10),
                        sm_id INT,
                        date_created TIMESTAMP,
                        oof_shard VARCHAR(10)
);

-- Таблица для товаров в заказах
CREATE TABLE order_items (
                             order_uid VARCHAR(255),
                             chrt_id INT,
                             price REAL,
                             rid VARCHAR(255),
                             name VARCHAR(255),
                             sale INT,
                             size VARCHAR(50),
                             total_price REAL,
                             nm_id INT,
                             brand VARCHAR(255),
                             status INT,
                             PRIMARY KEY (order_uid, chrt_id),
                             FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
);

-- Тестовое заполнение данными
INSERT INTO orders (
    order_uid, track_number, entry, delivery_info, payment_info,
    locale, internal_signature, customer_id, delivery_service,
    shardkey, sm_id, date_created, oof_shard
) VALUES (
             'b563feb7b2b84b6test', 'WBILMTESTTRACKtest', 'WBILtest',
             '{"name": "Test Testov", "phone": "+9990000000", "zip": "419420", "city": "TestBurg", "address": "Test street 15", "region": "Test", "email": "test@wb.ru"}',
             '{"transaction": "b563feb7b2b84b6test", "request_id": "", "currency": "RU", "provider": "wbpay", "amount": 419420, "payment_dt": 1637907727, "bank": "alphaTest", "delivery_cost": 1500, "goods_total": 317, "custom_fee": 0}',
             'ru', '', 'test', 'myTest', '9', 99, '2021-11-26T06:22:19Z', '1'
         );

-- Вставка данных для товара в заказе
INSERT INTO order_items (
    order_uid, chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status
) VALUES (
             'b563feb7b2b84b6test', 9934930, 453, 'ab4219087a764ae0btest', 'Mascaras', 30, '0', 317, 2389212, 'Vivienne Sabo', 202
         );
