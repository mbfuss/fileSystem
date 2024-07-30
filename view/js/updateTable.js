// Экспортируемая функция updateTable используется для обновления таблицы с файлами на странице
// Параметры:
// - data: массив данных о файлах и директориях
// - fileTableBody: элемент таблицы для отображения файлов
// - navigateToDirectory: функция для перехода в директорию при клике на имя директории

export const updateTable = (data, fileTableBody, navigateToDirectory) => {
    // Очищаем текущее содержимое таблицы.
    fileTableBody.innerHTML = '';

    // Для каждого файла создаем строку в таблице
    // Для каждого файла создаем строку в таблице
    data.forEach(file => {
        // Создаем элементы строки и ячеек
        const row = document.createElement('tr');
        const cellIsFile = document.createElement('td');
        const cellFileName = document.createElement('td');
        const cellFileSize = document.createElement('td');

        // Устанавливаем текстовое содержимое ячеек
        cellIsFile.textContent = file.is_file;
        cellFileName.textContent = file.name;
        cellFileSize.textContent = file.format_size;

        // Добавляем класс для ячейки имени файла
        cellFileName.classList.add('file-name');

        // Если это директория, добавляем класс и обработчик клика
        if (file.is_file === "Дир") {
            row.classList.add('directory'); // Добавляем класс для директории
            cellFileName.addEventListener('click', () => navigateToDirectory(file.name));
        }

        // Добавляем ячейки в строку
        row.appendChild(cellIsFile);
        row.appendChild(cellFileName);
        row.appendChild(cellFileSize);

        // Добавляем строку в тело таблицы
        fileTableBody.appendChild(row);
    });
};
