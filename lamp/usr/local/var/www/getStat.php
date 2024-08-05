<?php
require 'connDB.php';

// Создаем подключение
$conn = getDatabaseConnection();

// SQL-запрос для получения всех данных из таблицы file_info
$sql = "SELECT id, path, size, duration, request_time FROM file_info";
$result = $conn->query($sql);

// Начинаем формирование HTML-кода страницы
echo '<!DOCTYPE html>';
echo '<html>';
echo '<head>';
echo '<title>File Info Table and Chart</title>';
echo '<button onclick="window.history.back();">Назад</button>';
// Подключаем Chart.js
echo '<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>';
// CSS-стили для оформления таблицы и графика
echo '<style>';
echo 'table { width: 100%; border-collapse: collapse; }'; // Таблица занимает 100% ширины контейнера и границы ячеек схлопнуты
echo 'th, td { border: 1px solid white; padding: 8px; text-align: left; }'; // Границы ячеек цвет, отступы внутри ячеек 8px, текст выровнен по левому краю
echo 'th { background-color: white; color: black; }'; // Фон заголовков таблицы, цвет текста
echo 'tr:nth-child(even) { background-color: #007bff; }'; // Фон четных строк
echo 'tr:nth-child(odd) { background-color: white; }'; // Фон нечетных строк
echo '#myChart { max-width: 100%; height: 400px; margin-top: 20px; }'; // Устанавливаем высоту и отступ графика
echo '</style>';
echo '</head>';
echo '<body>';
echo '<canvas id="myChart"></canvas>'; // Контейнер для графика
//  $result — переменная, содержит объект результата запроса к базе данных, cоздаётся при выполнении SQL-запроса
// num_rows — это свойство объекта результата, которое содержит количество строк в результирующем наборе данных
if ($result->num_rows > 0) {
    echo '<html lang="ru">';
    echo '<meta charset="UTF-8">';
    // Начинаем формирование таблицы
    echo '<table>';
    echo '<tr>';
    // Выводим имена столбцов с подписями
    echo '<th>ID</th>';
    echo '<th>Путь</th>';
    echo '<th>Размер (bytes)</th>';
    echo '<th>Продолжительность (ms)</th>';
    echo '<th>Время запроса</th>';
    echo '</tr>';
    
    // Массив для хранения данных графика
    $chartData = array(); // Инициализация пустого массива для данных графика

    // Извлекаем и выводим данные каждой строки результата
    while($row = $result->fetch_assoc()) { // Извлекаем строку результата в виде ассоциативного массива
        // Добавляем данные в массив для графика
        $chartData[] = $row; // Добавляем текущую строку в массив данных графика
    }
        // Выводим таблицу
        foreach ($chartData as $row) {
            echo '<tr>'; // Открываем строку таблицы
            foreach ($row as $value) { // Перебираем значения строки
                echo '<td>' . htmlspecialchars($value) . '</td>'; // Выводим каждое значение в отдельной ячейке, экранируем специальные символы
            }
            echo '</tr>'; // Закрываем строку таблицы
        }
    // Сортируем данные по размеру (от меньшего к большему)
    usort($chartData, function($a, $b) {
    // Если результат меньше нуля, $a будет стоять перед $b
    // Если результат равен нулю, порядок элементов останется неизменным
    // Если результат больше нуля, $b будет стоять перед $a
        return $a['size'] - $b['size'];
    });

    echo '</table>'; // Закрываем таблицу

    // Преобразуем данные в формат JSON для использования в JavaScript
    $chartDataJson = json_encode($chartData); // Преобразуем массив данных в JSON строку
} else {
    // Если данных нет, выводим сообщение
    echo '<p>Данных не обнаружено</p>'; 
    $chartDataJson = json_encode([]); // Устанавливаем пустой массив в формате JSON на случай, если ответ от БД пустой
}
// Закрываем соединение с базой данных
$conn->close(); // Закрываем соединение с базой данных
?>

<script>
    // Функция для построения графика
    function renderChart(data) {
        // Извлекаем размеры и длительности из данных
        const sizes = data.map(item => item.size);
        const durations = data.map(item => item.duration);

        // Логарифмируем значения
        // const logSizes = sizes.map(size => Math.log10(size + 1)); // Используем log10 для логарифмирования и добавляем 1 для избежания log(0)
        // const logDurations = durations.map(duration => Math.log10(duration + 1)); // Используем log10 для логарифмирования и добавляем 1 для избежания log(0)

        // Получаем контекст канваса для построения графика
        const ctx = document.getElementById('myChart').getContext('2d');
        
        // Создаем новый график
        new Chart(ctx, {
            type: 'line', // Используем линейный график
            data: {
                labels: sizes, // размер в качестве меток по оси X
                datasets: [{
                    label: 'Зависимость времени обработки от размера директории',
                    data: durations, //  данные для оси Y
                    borderColor: 'rgba(54, 162, 235, 1)', // Цвет линии
                    backgroundColor: 'rgba(54, 162, 235, 1)', // Цвет фона линии
                    borderWidth: 2, // Толщина линии
                    pointBackgroundColor: 'rgba(54, 162, 235, 2)', // Цвет точек на линии
                    pointBorderColor: '#fff' // Цвет границы точек
                }]
            },
            options: {
                scales: {
                    x: {
                        // type: 'logarithmic',
                        title: {
                            display: true, // Показывать заголовок оси X
                            text: 'Размер (bytes)', // Текст заголовка оси X
                        }
                    },
                    y: {
                        title: {
                            display: true, // Показывать заголовок оси Y
                            text: 'Продолжительность (ms)' // Текст заголовка оси Y
                        }
                    }
                }
            }
        });
    }

    // Получаем данные и строим график после загрузки страницы
    document.addEventListener('DOMContentLoaded', function() {
        const chartData = <?php echo $chartDataJson; ?>; // Вставляем JSON данные из PHP в JavaScript
        renderChart(chartData); // Вызываем функцию для построения графика
    });
</script>

</body>
</html>
