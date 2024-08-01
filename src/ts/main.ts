// Импортируем необходимые модули и функции
import { fetchData } from './fetchData';
import { updateTable } from './updateTable';
import { navigateToDirectory, navigateBack } from './navigate';
import { addEventHandlers } from './eventHandlers';
import { fetchConfig } from './envConfigLoad';
import "../styles.css";

const asc: 'asc' = "asc";
const desc: 'desc' = "desc";

// Интерфейсы для элементов DOM
interface FileTableRow {
    is_file: string;
    name: string;
    format_size: string;
}

// Ждем загрузки DOM перед выполнением скрипта
document.addEventListener("DOMContentLoaded", async () => {
    // Получаем элементы управления из DOM
    const sortOrderSlider = document.getElementById('sortOrder') as HTMLInputElement; // Слайдер для сортировки
    const fileTableBody = document.querySelector('#fileTable tbody') as HTMLTableSectionElement; // Тело таблицы для отображения файлов
    const currentPath = document.getElementById('currentPath') as HTMLElement; // Элемент для отображения текущего пути
    const cancelButton = document.getElementById('cancelButton') as HTMLButtonElement; // Кнопка для перехода назад
    const loader = document.getElementById('loader') as HTMLElement;

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

    // Функция для получения данных с сервера и обновления таблицы
    const fetchAndUpdateTable = () => {
        // Получаем порядок сортировки из слайдера
        const sortOrder: 'asc' | 'desc' = sortOrderSlider.value === "0" ? desc : asc;
        // Блокируем элементы управления до завершения загрузки конфигурации
        toggleControls(true);
        // Показать индикатор загрузки
        loader.style.display = 'block';
        // Получаем данные с сервера и обновляем таблицу
        fetchData(
            currentRoot,
            sortOrder,
            (data: FileTableRow[]) => {
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
    const toggleControls = (isDisabled: boolean) => {
        sortOrderSlider.classList.toggle('disabled', isDisabled);
        cancelButton.classList.toggle('disabled', isDisabled);
        fileTableBody.classList.toggle('disabled', isDisabled);
    };

    // Добавляем обработчики событий к элементам управления.
    addEventHandlers(
        sortOrderSlider,
        cancelButton,
        fetchAndUpdateTable,
        navigateBack(getCurrentRoot, setCurrentRoot, fetchAndUpdateTable, rootDir),
        currentRoot
    );
    // Начальный запрос данных при загрузке страницы
    fetchAndUpdateTable();
});
