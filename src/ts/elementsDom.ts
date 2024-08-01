// Возвращаем элементы DOM
export const getDomElements = () => {
    return {
        sortOrderSlider: document.getElementById('sortOrder') as HTMLInputElement, // Слайдер для сортировки
        fileTableBody: document.querySelector('#fileTable tbody') as HTMLTableSectionElement, // Тело таблицы для отображения файлов
        currentPath: document.getElementById('currentPath') as HTMLElement, // Элемент для отображения текущего пути
        cancelButton: document.getElementById('cancelButton') as HTMLButtonElement, // Кнопка для перехода назад
        loader: document.getElementById('loader') as HTMLElement // Индикатор загрузки
    };
};
