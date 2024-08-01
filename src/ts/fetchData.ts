// Тип данных, получаемый с сервера
interface FileData {
    is_file: string;
    name: string;
    format_size: string;
}

// Тип для функции обновления таблицы
type UpdateTableFunction = (data: FileData[]) => void;

// Экспортируемая асинхронная функция fetchData используется для получения данных с сервера
// Параметры:
// - root: текущий путь к директории, для которой нужно получить данные
// - sortOrder: порядок сортировки, который может быть 'asc' (по возрастанию) или 'desc' (по убыванию)
// - updateTable: функция, которая обновляет таблицу с файлами на основе полученных данных
// - currentPath: элемент, в котором отображается текущий путь
export const fetchData = async (
    root: string,
    sortOrder: 'asc' | 'desc',
    updateTable: UpdateTableFunction,
    currentPath: HTMLElement
): Promise<void> => {
    // Обновляем текст текущего пути в элементе currentPath
    currentPath.textContent = `Текущий путь: ${root}`;

    try {
        // Делаем GET-запрос к серверу, передавая текущий путь и порядок сортировки в качестве параметров
        const response = await fetch(`/fs?root=${encodeURIComponent(root)}&sort=${sortOrder}`, { method: "GET" });

        // Проверяем успешность ответа. Если ответ не успешный, выбрасываем ошибку
        if (!response.ok) throw new Error('Ошибка ответа сети');

        // Парсим JSON ответ от сервера
        const data = await response.json();

        // Обновляем таблицу с файлами, используя полученные данные
        updateTable(data.data);
    } catch (error) {
        // Логируем ошибку в консоль, если произошла ошибка при выполнении запроса
        console.error('Ошибка получения данных:', error);
        alert("Ошибка получения данных");
    }
};