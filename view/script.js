// Ждем загрузки DOM перед выполнением скрипта
document.addEventListener("DOMContentLoaded", () => {
    // Получаем элементы управления из DOM
    const sortOrderSlider = document.getElementById('sortOrder'); // Слайдер для сортировки
    const fileTableBody = document.querySelector('#fileTable tbody'); // Таблица для отображения файлов
    const currentPath = document.getElementById('currentPath'); // Элемент для отображения текущего пути
    const cancelButton = document.getElementById('cancelButton'); // Кнопка для перехода назад

    // Переменная для хранения текущего пути
    let currentRoot = '/';

    // Асинхронная функция для получения данных с сервера
    const fetchData = async (root = currentRoot) => {
        // Получаем порядок сортировки из слайдера
        const sortOrder = sortOrderSlider.value === "0" ? "desc" : "asc";
        // Обновляем текст текущего пути
        currentPath.textContent = `Текущий путь: ${root}`;

        try {
            // Делаем запрос к серверу для получения списка файлов
            const response = await fetch(`http://localhost:8080/fs?root=${encodeURIComponent(root)}&sort=${sortOrder}`);
            // Проверяем успешность ответа
            if (!response.ok) throw new Error('Network response was not ok');
            // Парсим JSON ответ
            const data = await response.json();
            // Обновляем таблицу с файлами
            updateTable(data);
        } catch (error) {
            // Обрабатываем ошибку запроса
            console.error('Fetching error:', error);
        }
    };

    // Функция для обновления таблицы с файлами
    const updateTable = (data) => {
        // Очищаем текущее содержимое таблицы
        fileTableBody.innerHTML = '';
        // Для каждого файла создаем строку в таблице
        data.forEach(file => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${file.is_file}</td>
                <td class="file-name">${file.name}</td>
                <td>${file.format_size}</td>
            `;
            // Если это директория, добавляем класс и обработчик клика
            if (file.is_file === "Дир") {
                row.classList.add('directory'); // Добавляем класс для директории, чтобы работали css стили
                row.querySelector('.file-name').addEventListener('click', () => navigateToDirectory(file.name));
            }
            // Добавляем строку в таблицу
            fileTableBody.appendChild(row);
        });
    };

    // Функция для перехода в директорию
    const navigateToDirectory = (dirName) => {
        // Обновляем текущий путь
        currentRoot = currentRoot === '/' ? `/${dirName}` : `${currentRoot}/${dirName}`;
        // Запрашиваем данные для новой директории
        fetchData(currentRoot);
    };

    // Функция для возврата к предыдущей директории
    const navigateBack = () => {
        // Если мы уже в корневом каталоге, ничего не делаем
        if (currentRoot === '/') return;
        // Разделяем путь на части и удаляем пустые элементы
        const pathParts = currentRoot.split('/').filter(part => part.length > 0);
        // Удаляем последний элемент пути (текущую директорию)
        pathParts.pop();
        // Обновляем текущий путь
        currentRoot = pathParts.length > 0 ? `/${pathParts.join('/')}` : '/';
        // Запрашиваем данные для новой директории
        fetchData(currentRoot);
    };

    // Обработчик события для слайдера сортировки
    sortOrderSlider.addEventListener('input', () => fetchData(currentRoot));
    // Обработчик события для кнопки "Назад"
    cancelButton.addEventListener('click', navigateBack);

    // Начальный запрос данных при загрузке страницы
    fetchData();
});
