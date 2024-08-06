<?php

require 'connDB.php';

try {
    // Создаем подключение
    $conn = getDatabaseConnection();

    // Извлекаем весь содержимое тела POST-запроса как строку
    $data = file_get_contents('php://input');
    // Преобразует JSON-строку в ассоциативный массив 
    $data = json_decode($data, true);

    // Проверяем содержатся ли в массиве $data ключи 'path', 'size' и 'duration'
    if (!isset($data['path']) || !isset($data['size']) || !isset($data['duration'])) {
        throw new Exception("Неверные данные");
    }

    $path = $conn->real_escape_string($data['path']); // экранирование, убираем лишнии символы 
    $size = (int)$data['size']; // Преобразуем к целому числу
    $duration = (int)$data['duration']; // Преобразуем к целому числу

    // Вставляем данные в базу данных, используя NOW() для текущего времени
    $sql = "INSERT INTO file_info (path, size, duration, request_time) VALUES ('$path', $size, $duration, NOW())";

    // Проверяем был ли выполен sql запрос
    if ($conn->query($sql) === TRUE) {
        echo "Новая запись успешно создана";
    } else {
        throw new Exception("Ошибка выполнения запроса: " . $conn->error);
    }
} catch (Exception $e) {
    echo "Ошибка: " . htmlspecialchars($e->getMessage());
} finally {
    // Закрываем соединение
    // переменная $conn существует и является экземпляром класса mysqli
    if (isset($conn) && $conn instanceof mysqli) {
        $conn->close();
    }
}
?>
