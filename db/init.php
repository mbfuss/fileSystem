<?php
// Настройки подключения
$servername = "localhost";  
$username = "root";        
$password = "12346";             
$dbname = "fsDataBase";   

// Создание подключения
$conn = new mysqli($servername, $username, $password);

// Проверка подключения
if ($conn->connect_error) {
    die("Ошибка подключения: " . $conn->connect_error);
}

// Создание базы данных, если она не существует
$sql = "CREATE DATABASE IF NOT EXISTS $dbname";
if ($conn->query($sql) === TRUE) {
    echo "База данных создана или уже существует.<br>";
} else {
    die("Ошибка создания базы данных: " . $conn->error);
}

// Выбор базы данных
$conn->select_db($dbname);

// Создание таблицы, если она не существует
$sql = "CREATE TABLE IF NOT EXISTS file_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    путь VARCHAR(255) NOT NULL,
    размер INT NOT NULL,
    время_выполнения INT NOT NULL,
    во_сколько_запрос_сделан INT NOT NULL
)";

if ($conn->query($sql) === TRUE) {
    echo "Таблица создана или уже существует.<br>";
} else {
    die("Ошибка создания таблицы: " . $conn->error);
}

//  Вставка начальных данных
$sql = "INSERT INTO file_info (путь, размер, время_выполнения, во_сколько_запрос_сделан)
VALUES ('/path/to/file', 12345, 60, UNIX_TIMESTAMP())";
if ($conn->query($sql) === TRUE) {
    echo "Начальные данные добавлены.<br>";
} else {
    die("Ошибка добавления данных: " . $conn->error);
}

// Закрытие подключения
$conn->close();
?>
