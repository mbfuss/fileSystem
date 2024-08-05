<?php
function getDatabaseConnection() {
    // Устанавливаем параметры подключения к базе данных
    $servername = "localhost";
    $username = "root";
    $password = "12346";
    $dbname = "fsdatabase";

    try {
        // Создаем подключение
        $conn = new mysqli($servername, $username, $password, $dbname);

        // Проверяем подключение
        if ($conn->connect_error) {
            throw new Exception("Ошибка подключение к БД: " . $conn->connect_error);
        }

        return $conn;
    } catch (Exception $e) {
        // Ловим исключение и выводим сообщение об ошибке
        echo "Error: " . htmlspecialchars($e->getMessage());
        return null;
    }
}
?>
