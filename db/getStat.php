<?php
// Настройки подключения
$servername = "localhost";
$username = "root";
$password = "12346";
$dbname = "fsDataBase";

// Создание подключения
$conn = new mysqli($servername, $username, $password, $dbname);

// Проверка подключения
if ($conn->connect_error) {
    die("Ошибка подключения: " . $conn->connect_error);
}

// Выполнение SELECT-запроса
$sql = "SELECT id, путь, размер, время_выполнения, во_сколько_запрос_сделан FROM file_info";
$result = $conn->query($sql);

// Проверка наличия результатов
if ($result->num_rows > 0) {
    // Вывод данных каждой строки
    while($row = $result->fetch_assoc()) {
        echo "ID: " . $row["id"]. " - Путь: " . $row["путь"]. " - Размер: " . $row["размер"]. " - Время выполнения: " . $row["время_выполнения"]. " - Во сколько запрос сделан: " . $row["во_сколько_запрос_сделан"]. "<br>";
    }
} else {
    echo "0 результатов";
}

// Закрытие подключения
$conn->close();
?>
