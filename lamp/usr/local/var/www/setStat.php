<?php
// Устанавливаем параметры подключения к базе данных
$servername = "localhost";
$username = "root";
$password = "12346";
$dbname = "fsdatabase";

// Создаем подключение
$conn = new mysqli($servername, $username, $password, $dbname);

// Проверяем подключение
if ($conn->connect_error) {
    die("Connection failed: " . $conn->connect_error);
}

// Получаем JSON-данные из POST-запроса
$data = file_get_contents('php://input');
$data = json_decode($data, true);

// Проверяем корректность данных
if (isset($data['path']) && isset($data['size']) && isset($data['duration']) && isset($data['request_time'])) {
    $path = $conn->real_escape_string($data['path']);
    $size = (int)$data['size']; // Преобразуем к целому числу
    $duration = (int)$data['duration']; // Преобразуем к целому числу
    $request_time = $conn->real_escape_string($data['request_time']);

    // Вставляем данные в базу данных
    $sql = "INSERT INTO file_info (path, size, duration, request_time) VALUES ('$path', $size, $duration, '$request_time')";
    
    if ($conn->query($sql) === TRUE) {
        echo "New record created successfully";
    } else {
        echo "Error: " . $sql . "<br>" . $conn->error;
    }
} else {
    echo "Invalid data";
}

// Закрываем соединение
$conn->close();
?>
