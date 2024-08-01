// fetchAndUpdateTable.js

import { fetchData } from './fetchData.js';
import { updateTable } from './updateTable.js';
import { navigateToDirectory } from './navigate.js';

const asc = "asc";
const desc = "desc";



export const fetchAndUpdateTable = (currentRoot, setCurrentRoot, sortOrderSlider, fileTableBody, currentPath, loader, toggleControls) => {

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
            updateTable(data, fileTableBody, navigateToDirectory(currentRoot, setCurrentRoot, () => fetchAndUpdateTable(currentRoot, setCurrentRoot, sortOrderSlider, fileTableBody, currentPath, loader, toggleControls)));
        },
        currentPath
    ).finally(() => {
        // Скрыть индикатор загрузки
        loader.style.display = 'none';
        // Разблокируем элементы управления после завершения загрузки конфигурации и данных
        toggleControls(false);
    });
};

export const toggleControls = (isDisabled, sortOrderSlider, cancelButton, fileTableBody) => {
    sortOrderSlider.classList.toggle('disabled', isDisabled);
    cancelButton.classList.toggle('disabled', isDisabled);
    fileTableBody.classList.toggle('disabled', isDisabled);
};
