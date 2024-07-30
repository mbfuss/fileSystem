// Экспортируемая функция updateTable используется для обновления таблицы с файлами на странице
// Параметры:
// - data: массив данных о файлах и директориях
// - fileTableBody: элемент таблицы для отображения файлов
// - navigateToDirectory: функция для перехода в директорию при клике на имя директории

export const updateTable = (data, fileTableBody, navigateToDirectory) => {
    // Очищаем текущее содержимое таблицы.
    fileTableBody.innerHTML = '';

    // Для каждого файла создаем строку в таблице.
    data.forEach(file => {
        const row = document.createElement('tr');
        // Заполняем строку данными о файле.
        row.innerHTML = `
            <td>${file.is_file}</td>
            <td class="file-name">${file.name}</td>
            <td>${file.format_size}</td>
        `;

        // Если это директория, добавляем класс и обработчик клика
        if (file.is_file === "Дир") {
            row.classList.add('directory'); // Добавляем класс для директории, чтобы работали CSS стили
            row.querySelector('.file-name').addEventListener('click', () => navigateToDirectory(file.name));
        }

        // Добавляем строку в таблицу.
        fileTableBody.appendChild(row);
    });
};
