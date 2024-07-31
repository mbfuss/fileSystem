// Импортируем необходимые модули и функции
import { fetchData } from './fetchData.js';
import { updateTable } from './updateTable.js';
import { navigateToDirectory, navigateBack } from './navigate.js';
import { addEventHandlers } from './eventHandlers.js';
import { fetchConfig } from './envConfigLoad.js';
import './styles.css'

const asc = "asc";
const desc = "desc";


// Ждем загрузки DOM перед выполнением скрипта
document.addEventListener("DOMContentLoaded", async () => {

    // Получаем элементы управления из DOM.
    const sortOrderSlider = document.getElementById('sortOrder'); // Слайдер для сортировки
    const fileTableBody = document.querySelector('#fileTable tbody'); // Тело таблицы для отображения файлов
    const currentPath = document.getElementById('currentPath'); // Элемент для отображения текущего пути
    const cancelButton = document.getElementById('cancelButton'); // Кнопка для перехода назад
    const loader = document.getElementById('loader')

    // Переменная для хранения корневого пути
    const rootDir = await fetchConfig();
    console.log('Root directory from config:', rootDir);

    // Переменная для хранения текущего пути
    let currentRoot = rootDir;


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
        // Блокируем элементы управления до завершения загрузки конфигурации
        toggleControls(true);
        // Показать индикатор загрузки
        loader.style.display = 'block';
        // Получаем данные с сервера и обновляем таблицу
        fetchData(
            currentRoot,
            sortOrder,
            (data) => {
                updateTable(data, fileTableBody, navigateToDirectory(getCurrentRoot, setCurrentRoot, fetchAndUpdateTable));
            },
            currentPath
        ).finally(() => {
            // Скрыть индикатор загрузки
            loader.style.display = 'none';
            // Разблокируем элементы управления после завершения загрузки конфигурации и данных
            toggleControls(false);
        });
    };

    // Функция для блокировки и разблокировки элементов управления
    const toggleControls = (isDisabled) => {
        sortOrderSlider.classList.toggle('disabled', isDisabled);
        cancelButton.classList.toggle('disabled', isDisabled);
        fileTableBody.classList.toggle('disabled',isDisabled)
    };


    // Добавляем обработчики событий к элементам управления.
    addEventHandlers(sortOrderSlider, cancelButton, fetchAndUpdateTable, navigateBack(getCurrentRoot, setCurrentRoot, fetchAndUpdateTable, rootDir));
    // Начальный запрос данных при загрузке страницы
    fetchAndUpdateTable();

});