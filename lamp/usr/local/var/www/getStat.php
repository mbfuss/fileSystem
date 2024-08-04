<?php
// Настройки подключения
$servername = "localhost";
$username = "root";
$password = "12346";
$dbname = "fsdatabase";

// Создание подключения
$conn = new mysqli($servername, $username, $password, $dbname);

// Проверка подключения
if ($conn->connect_error) {
    die("Ошибка подключения: " . $conn->connect_error);
}

// Выполнение SELECT-запроса
$sql = "SELECT id, path, size, duration, request_time FROM file_info";
$result = $conn->query($sql);

// Проверка наличия результатов
if ($result->num_rows > 0) {
    // Вывод данных каждой строки
    while($row = $result->fetch_assoc()) {
        echo "ID: " . $row["id"]. " - Путь: " . $row["path"]. " - Размер: " . $row["size"]. " - Время выполнения: " . $row["duration"]. " - Во сколько запрос сделан: " . $row["request_time"]. "<br>";
    }
} else {
    echo "0 результатов";
}

// Закрытие подключения
$conn->close();

