<?php
header('Content-Type: application/json');

// Настройки подключения к базе данных
$servername = "localhost";
$username = "root";
$password = "12346";
$dbname = "fsDataBase";

// Создание подключения
$conn = new mysqli($servername, $username, $password, $dbname);

// Проверка подключения
if ($conn->connect_error) {
    echo json_encode(["error" => "Ошибка подключения к базе данных"]);
    exit;
}

// Получение данных из POST-запроса
$data = json_decode(file_get_contents('php://input'), true);

// Проверка наличия необходимых данных
if (!isset($data['путь'], $data['размер'], $data['время_выполнения'])) {
    echo json_encode(["error" => "Отсутствуют необходимые данные"]);
    exit;
}

$путь = $data['путь'];
$размер = $data['размер'];
$время_выполнения = $data['время_выполнения'];
$во_сколько_запрос_сделан = time();

// Подготовка и выполнение INSERT-запроса
$stmt = $conn->prepare("INSERT INTO file_info (путь, размер, время_выполнения, во_сколько_запрос_сделан) VALUES (?, ?, ?, ?)");
$stmt->bind_param("siis", $путь, $размер, $время_выполнения, $во_сколько_запрос_сделан);

if ($stmt->execute()) {
    echo json_encode(["success" => "Данные успешно сохранены"]);
} else {
    echo json_encode(["error" => "Ошибка выполнения запроса"]);
}

// Закрытие подключения
$stmt->close();
$conn->close();
?>
