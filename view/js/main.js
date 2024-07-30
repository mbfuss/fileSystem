// js/main.js

// Импортируем необходимые модули и функции
import { fetchData } from './fetchData.js';
import { updateTable } from './updateTable.js';
import { navigateToDirectory, navigateBack } from './navigate.js';
import { addEventHandlers } from './eventHandlers.js';

const asc = "asc";
const desc = "desc";

// Ждем загрузки DOM перед выполнением скрипта
document.addEventListener("DOMContentLoaded", () => {
    // Получаем элементы управления из DOM.
    const sortOrderSlider = document.getElementById('sortOrder'); // Слайдер для сортировки
    const fileTableBody = document.querySelector('#fileTable tbody'); // Тело таблицы для отображения файлов
    const currentPath = document.getElementById('currentPath'); // Элемент для отображения текущего пути
    const cancelButton = document.getElementById('cancelButton'); // Кнопка для перехода назад

    // Переменная для хранения текущего пути
    let currentRoot = '/';

    // Функция для обновления текущего пути
    const setCurrentRoot = (newRoot) => {
        currentRoot = newRoot;
    };

    // Функция для получения текущего пути
    const getCurrentRoot = () => currentRoot;

    // Функция для получения данных с сервера и обновления таблицы
    const fetchAndUpdateTable = () => {
        // Получаем порядок сортировки из слайдера
        const sortOrder = sortOrderSlider.value === "0" ? desc : asc;
        // Получаем данные с сервера и обновляем таблицу
        fetchData(currentRoot, sortOrder, (data) => updateTable(data, fileTableBody, navigateToDirectory(getCurrentRoot, setCurrentRoot, fetchAndUpdateTable)), currentPath);
    };

    // Добавляем обработчики событий к элементам управления.
    addEventHandlers(sortOrderSlider, cancelButton, fetchAndUpdateTable, navigateBack(getCurrentRoot, setCurrentRoot, fetchAndUpdateTable));

    // Начальный запрос данных при загрузке страницы.
    fetchAndUpdateTable();
});
