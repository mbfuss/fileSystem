// Определяем типы для параметров функции
type FetchDataFunction = (currentRoot: string, sortOrder: 'asc' | 'desc') => void;
type NavigateBackFunction = () => void;

// Экспортируемая функция addEventHandlers добавляет обработчики событий к элементам управления на странице
export const addEventHandlers = (
    sortOrderSlider: HTMLInputElement,  // Элемент слайдера
    cancelButton: HTMLButtonElement,    // Кнопка отмены
    statButton: HTMLButtonElement,      // Кнопка статистики
    fetchData: FetchDataFunction,       // Функция для запроса данных
    navigateBack: NavigateBackFunction, // Функция для возврата
    currentRoot: string                 // Текущий путь
): void => {

    // Добавляем обработчик события 'input' к слайдеру сортировки
    sortOrderSlider.addEventListener('input', () => {
        // Определяем порядок сортировки на основе значения слайдера
        const sortOrder: 'asc' | 'desc' = sortOrderSlider.value === "0" ? "desc" : "asc";

        // Вызываем функцию fetchData для получения данных с сервера с новым порядком сортировки
        fetchData(currentRoot, sortOrder);
    });

    // Добавляем обработчик события 'click' к кнопке отмены
    cancelButton.addEventListener('click', navigateBack);
    // Добавляем обработчик события 'click' к кнопке статистики
    statButton.addEventListener('click', () => {
        // Переход на указанный URL
        window.location.href = '/getfileinfo';
    });
};
