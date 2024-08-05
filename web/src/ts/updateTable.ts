// Тип данных, получаемый с сервера
interface FileData {
    is_file: string; // Указывает, является ли элемент файлом или директорией ("Дир" для директорий)
    name: string; // Имя файла или директории
    format_size: string; // Размер файла в формате строки
}

// Тип для функции навигации
type NavigateToDirectoryFunction = (dirName: string) => void;

// Экспортируемая функция updateTable используется для обновления таблицы с файлами на странице
// Параметры:
// - data: массив данных о файлах и директориях
// - fileTableBody: элемент таблицы для отображения файлов
// - navigateToDirectory: функция для перехода в директорию при клике на имя директории
export const updateTable = (
    data: FileData[],
    fileTableBody: HTMLTableSectionElement,
    navigateToDirectory: NavigateToDirectoryFunction
): void => {
    // Очищаем текущее содержимое таблицы
    fileTableBody.innerHTML = '';

    // Проходим по каждому файлу/директории в полученных данных
    data.forEach(file => {
        // Создаем элементы строки и ячеек таблицы
        const row = document.createElement('tr'); // Создаем элемент строки таблицы
        const cellIsFile = document.createElement('td'); // Создаем ячейку для типа элемента (файл/директория)
        const cellFileName = document.createElement('td'); // Создаем ячейку для имени файла/директории
        const cellFileSize = document.createElement('td'); // Создаем ячейку для размера файла

        // Устанавливаем текстовое содержимое ячеек
        cellIsFile.textContent = file.is_file; // Устанавливаем текст в ячейке типа элемента
        cellFileName.textContent = file.name; // Устанавливаем текст в ячейке имени файла/директории
        cellFileSize.textContent = file.format_size; // Устанавливаем текст в ячейке размера файла

        // Добавляем класс для ячейки имени файла
        cellFileName.classList.add('file-name'); // Добавляем CSS-класс для стилизации ячейки с именем файла

        // Если элемент является директорией, добавляем соответствующий класс и обработчик клика
        if (file.is_file === "Дир") {
            row.classList.add('directory'); // Добавляем CSS-класс для строки, обозначающей директорию
            // Добавляем обработчик события клика для перехода в директорию
            cellFileName.addEventListener('click', () => navigateToDirectory(file.name));
        }

        // Добавляем ячейки в строку таблицы
        row.appendChild(cellIsFile); // Добавляем ячейку типа элемента в строку
        row.appendChild(cellFileName); // Добавляем ячейку имени файла/директории в строку
        row.appendChild(cellFileSize); // Добавляем ячейку размера файла в строку

        // Добавляем строку в тело таблицы
        fileTableBody.appendChild(row); // Добавляем заполненную строку в элемент tbody таблицы
    });
};
