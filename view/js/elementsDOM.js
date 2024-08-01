export const getDomElements = () => {
    const sortOrderSlider = document.getElementById('sortOrder'); // Слайдер для сортировки
    const fileTableBody = document.querySelector('#fileTable tbody'); // Тело таблицы для отображения файлов
    const currentPath = document.getElementById('currentPath'); // Элемент для отображения текущего пути
    const cancelButton = document.getElementById('cancelButton'); // Кнопка для перехода назад
    const loader = document.getElementById('loader'); // Индикатор загрузки
    return { sortOrderSlider, fileTableBody, currentPath, cancelButton, loader };
};