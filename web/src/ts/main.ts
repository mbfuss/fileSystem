import {fetchData, fetchRootConfig} from './fetchData';
import { updateTable } from './updateTable';
import { navigateToDirectory, navigateBack } from './navigate';
import { addEventHandlers } from './eventHandlers';
import { createFetchAndUpdateTable } from './fetchAndUpdateTable';
import "../styles/styles.css";
import {getDomElements} from "./elementsDom";

// Ждем загрузки DOM перед выполнением скрипта
document.addEventListener("DOMContentLoaded", async () => {
    // Получаем элементы управления из DOM
    const {
        sortOrderSlider,
        fileTableBody,
        currentPath,
        cancelButton,
        loader
    } = getDomElements();
    // Переменная для хранения корневого пути
    const rootDir: string = <string>await fetchRootConfig();

    // Переменная для хранения текущего пути
    let currentRoot: string = rootDir;

    // Функция для обновления текущего пути
    const setCurrentRoot = (newRoot: string) => {
        currentRoot = newRoot;
    };

    // Функция для получения текущего пути
    const getCurrentRoot = (): string => currentRoot;

    // Создаем функцию для получения данных с сервера и обновления таблицы
    // Запрашивает данные с сервера с помощью fetchData
    // Обновляет таблицу с помощью updateTable
    // Обрабатывает навигацию по директориям через navigateToDirectory, передавая текущий корневой путь и обновленную функцию fetchAndUpdateTable
    // Использует элементы управления  и индикатор загрузки
    // Получает текущий корневой путь
    const fetchAndUpdateTable = createFetchAndUpdateTable(
        fetchData,
        updateTable,
        (dirName: string) => navigateToDirectory(getCurrentRoot, setCurrentRoot, fetchAndUpdateTable)(dirName),
        sortOrderSlider,
        cancelButton,
        fileTableBody,
        loader,
        currentPath,
        getCurrentRoot
    );

    // Функция для привязки обработчиков событий к элементам управления
    // sortOrderSlider для обработки изменения порядка сортировки
    // cancelButton для обработки отмены действия
    // fetchAndUpdateTable для обновления таблицы
    // navigateBack для обработки навигации назад
    // currentRoot — текущий корневой путь
    addEventHandlers(
        sortOrderSlider,
        cancelButton,
        fetchAndUpdateTable,
        navigateBack(getCurrentRoot, setCurrentRoot, fetchAndUpdateTable, rootDir),
        currentRoot
    );

    // Начальный запрос данных при загрузке страницы
    await fetchAndUpdateTable();
});
