// main.js

// Импортируем необходимые модули и функции
import { fetchAndUpdateTable, toggleControls } from './fetchAndUpdateTable.js';
import { addEventHandlers } from './eventHandlers.js';
import { fetchConfig } from './envConfigLoad.js';
import { navigateBack } from './navigate.js';
import { getDomElements } from './elementsDOM.js';

// Ждем загрузки DOM перед выполнением скрипта
document.addEventListener("DOMContentLoaded", async () => {

    // Получаем элементы управления из DOM.
    const { sortOrderSlider, fileTableBody, currentPath, cancelButton, loader } = getDomElements();

    // Переменная для хранения корневого пути
    const rootDir = await fetchConfig();
    // Переменная для хранения текущего пути
    let currentRoot = rootDir;
    // Функция для обновления текущего пути
    const setCurrentRoot = (newRoot) => {
        currentRoot = newRoot;
    };
    // Функция для получения текущего пути
    const getCurrentRoot = () => currentRoot;
    // Передача функций и их вызов в качестве аргументов для функции addEventHandlers
    const updateTableCallback = () => fetchAndUpdateTable(
        currentRoot,
        setCurrentRoot,
        sortOrderSlider,
        fileTableBody,
        currentPath,
        loader,
        (isDisabled) => toggleControls(isDisabled, sortOrderSlider, cancelButton, fileTableBody)
    );

    const navigateBackCallback = navigateBack(
        getCurrentRoot,
        setCurrentRoot,
        updateTableCallback,
        rootDir
    );

    addEventHandlers(
        sortOrderSlider,
        cancelButton,
        updateTableCallback,
        navigateBackCallback
    );

    // Начальный запрос данных при загрузке страницы
    fetchAndUpdateTable(currentRoot,
        setCurrentRoot,
        sortOrderSlider,
        fileTableBody,
        currentPath,
        loader, (isDisabled) => toggleControls(isDisabled, sortOrderSlider, cancelButton, fileTableBody));

});
