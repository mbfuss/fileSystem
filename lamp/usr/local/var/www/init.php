<?php
// Настройки подключения
$servername = "localhost";  
$username = "root";        
$password = "12346";             
$dbname = "fsdatabase";   

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
    path VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    duration INT NOT NULL,
    request_time DATETIME NOT NULL
)";

if ($conn->query($sql) === TRUE) {
    echo "Таблица создана или уже существует.<br>";
} else {
    die("Ошибка создания таблицы: " . $conn->error);
}

//  Вставка начальных данных
$sql = "INSERT INTO file_info (path, size, duration, request_time)
VALUES ('/path/to/file', 12345, 60, UNIX_TIMESTAMP())";
if ($conn->query($sql) === TRUE) {
    echo "Начальные данные добавлены.<br>";
} else {
    die("Ошибка добавления данных: " . $conn->error);
}

// Закрытие подключения
$conn->close();
?>
