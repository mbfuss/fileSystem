import { fetchData } from './fetchData';
import { updateTable } from './updateTable';
import { navigateToDirectory, navigateBack } from './navigate';
import { addEventHandlers } from './eventHandlers';
import { fetchConfig } from './envConfigLoad';
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
    const rootDir: string = <string>await fetchConfig();
    console.log('Root directory from config:', rootDir);

    // Переменная для хранения текущего пути
    let currentRoot: string = rootDir;

    // Функция для обновления текущего пути
    const setCurrentRoot = (newRoot: string) => {
        currentRoot = newRoot;
    };

    // Функция для получения текущего пути
    const getCurrentRoot = (): string => currentRoot;

    // Создаем функцию для получения данных с сервера и обновления таблицы
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

    // Добавляем обработчики событий к элементам управления
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
