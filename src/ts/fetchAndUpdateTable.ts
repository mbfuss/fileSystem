// Интерфейсы для элементов DOM
interface FileTableRow {
    is_file: string;
    name: string;
    format_size: string;
}

// Определяем типы для параметров функции
type FetchDataFunction = (
    root: string,
    sortOrder: 'asc' | 'desc',
    updateTable: (data: FileTableRow[]) => void,
    currentPath: HTMLElement
) => Promise<void>;

type UpdateTableFunction = (
    data: FileTableRow[],
    fileTableBody: HTMLTableSectionElement,
    navigateToDirectory: (dirName: string) => void
) => void;

type ToggleControlsFunction = (isDisabled: boolean) => void;

type FetchAndUpdateTableFunction = () => Promise<void>;

// Экспортируемая функция
export const createFetchAndUpdateTable = (
    fetchData: FetchDataFunction,
    updateTable: UpdateTableFunction,
    navigateToDirectory: (dirName: string) => void,
    sortOrderSlider: HTMLInputElement,
    fileTableBody: HTMLTableSectionElement,
    loader: HTMLElement,
    currentPath: HTMLElement,
    getCurrentRoot: () => string
): FetchAndUpdateTableFunction => {
    return async () => {
        // Получаем порядок сортировки из слайдера
        const sortOrder: 'asc' | 'desc' = sortOrderSlider.value === "0" ? "desc" : "asc";

        // Блокируем элементы управления до завершения загрузки конфигурации
        toggleControls(true);

        // Показать индикатор загрузки
        loader.style.display = 'block';

        try {
            // Получаем данные с сервера и обновляем таблицу
            await fetchData(
                getCurrentRoot(),
                sortOrder,
                (data: FileTableRow[]) => {
                    updateTable(data, fileTableBody, navigateToDirectory);
                },
                currentPath
            );
        } finally {
            // Скрыть индикатор загрузки
            loader.style.display = 'none';
            // Разблокируем элементы управления после завершения загрузки конфигурации и данных
            toggleControls(false);
        }
    };
};

// Функция для блокировки и разблокировки элементов управления
const toggleControls: ToggleControlsFunction = (isDisabled: boolean) => {
    // Пример кода для блокировки элементов управления
    // sortOrderSlider.classList.toggle('disabled', isDisabled);
    // cancelButton.classList.toggle('disabled', isDisabled);
    // fileTableBody.classList.toggle('disabled', isDisabled);
};
