CREATE DATABASE  IF NOT EXISTS bicis_caba;

USE bicis_caba;

CREATE TABLE IF NOT EXISTS users(
	user_id INT,
	gender VARCHAR(2),
	age INT,
	sign_up_date VARCHAR(30)
);

CREATE TABLE IF NOT EXISTS station(
	station_id INT,
	latitude FLOAT,
    longitud FLOAT,
    station_name VARCHAR(30)
);

CREATE TABLE IF NOT EXISTS trips(
	user_id INT,
	origin_location_id INT,
	destiny_location_id INT,
	origin_date VARCHAR(30),
    destiny_date VARCHAR(30),
    duration INT
);