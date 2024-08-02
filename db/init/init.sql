-- Создание базы данных
CREATE DATABASE IF NOT EXISTS my_database;

-- Выбор базы данных
USE my_database;

-- Создание таблицы
CREATE TABLE IF NOT EXISTS file_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    путь VARCHAR(255) NOT NULL,
    размер INT NOT NULL,
    время_выполнения INT NOT NULL,
    во_сколько_запрос_сделан INT NOT NULL
    );

-- (Необязательно) Вставка начальных данных
-- INSERT INTO file_info (путь, размер, время_выполнения, во_сколько_запрос_сделан)
-- VALUES ('/path/to/file', 12345, 60, UNIX_TIMESTAMP());
