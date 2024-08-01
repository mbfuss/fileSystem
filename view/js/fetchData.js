// Экспортируемая асинхронная функция fetchData используется для получения данных с сервера
// Параметры:
// - root: текущий путь к директории, для которой нужно получить данные
// - sortOrder: порядок сортировки, который может быть 'asc' (по возрастанию) или 'desc' (по убыванию)
// - updateTable: функция, которая обновляет таблицу с файлами на основе полученных данных
// - currentPath: элемент, в котором отображается текущий путь

export const fetchData = async (root, sortOrder, updateTable, currentPath) => {
    // Обновляем текст текущего пути в элементе currentPath
    currentPath.textContent = `Текущий путь: ${root}`;

    try {
        // Делаем GET-запрос к серверу, передавая текущий путь и порядок сортировки в качестве параметров
        const response = await fetch(`/fs?root=${encodeURIComponent(root)}&sort=${sortOrder}`, { method: "GET" });
        // Проверяем успешность ответа. Если ответ не успешный, выбрасываем ошибку
        if (!response.ok) throw new Error('Ошибка ответа сети');

        // Парсим JSON ответ от сервера
        const data = await response.json();
        if (data.error_code !== 0) {
            throw new Error('Возникла проблема с операцией получения данных: ' + data.error_message);
        }
        // Обновляем таблицу с файлами, используя полученные данные
        updateTable(data.data);
    } catch (error) {
        alert(error.message);
        return null; // Возвращаем значение по умолчанию или null в случае ошибки
    }
};