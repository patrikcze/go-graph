CREATE DATABASE IF NOT EXISTS temperature_db;
USE temperature_db;

CREATE TABLE IF NOT EXISTS data (
id INT AUTO_INCREMENT PRIMARY KEY,
time DATETIME NOT NULL,
temperature FLOAT NOT NULL,
humidity FLOAT NOT NULL,
pressure FLOAT NOT NULL
);

INSERT INTO data (time, temperature, humidity, pressure) VALUES
    ('2022-01-01 00:00:00', 20, 78, 1023),
    ('2022-01-01 01:00:00', 21, 75, 1023),
    ('2022-01-01 02:00:00', 22, 34, 1025),
    ('2022-01-01 03:00:00', 23, 55, 1022),
    ('2022-01-01 04:00:00', 24, 33, 1021),
    ('2022-01-01 05:00:00', 25, 12, 1023),
    ('2022-01-01 06:00:00', 26, 64, 1024),
    ('2022-01-01 07:00:00', 27, 23, 1019),
    ('2022-01-01 08:00:00', 28, 22, 1022),
    ('2022-01-01 09:00:00', 10, 33, 1021),
    ('2022-01-01 10:00:00', 15, 55, 1022),
    ('2022-01-01 11:00:00', 31, 66, 1024),
    ('2022-01-01 12:00:00', 25, 33, 1025),
    ('2022-01-01 13:00:00', 33, 22, 1023),
    ('2022-01-01 14:00:00', 34, 22, 1021),
    ('2022-01-01 15:00:00', 35, 33, 1025),
    ('2022-01-01 16:00:00', 11, 33, 1026),
    ('2022-01-01 17:00:00', 0, 55, 1021),
    ('2022-01-01 18:00:00', -10, 50, 1012),
    ('2022-01-01 19:00:00', -5, 65, 1010),
    ('2022-01-01 20:00:00', -3.2, 95, 1023);
