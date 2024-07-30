// Экспортируемая функция navigateToDirectory используется для перехода в указанную директорию
// Параметры:
// - currentRoot: текущий путь к директории
// - setCurrentRoot: функция для обновления текущего пути
// - fetchData: функция для получения данных с сервера

export const navigateToDirectory = (currentRoot, setCurrentRoot, fetchData) => (dirName) => {
    // Обновляем текущий путь на основе выбранной директории
    const newRoot = currentRoot === '/' ? `/${dirName}` : `${currentRoot}/${dirName}`;
    // Устанавливаем новый текущий путь
    setCurrentRoot(newRoot);
    // Получаем данные для новой директории
    fetchData(newRoot);
};

// Экспортируемая функция navigateBack используется для возврата к предыдущей директории
// Параметры:
// - currentRoot: текущий путь к директории
// - setCurrentRoot: функция для обновления текущего пути
// - fetchData: функция для получения данных с сервера.

export const navigateBack = (currentRoot, setCurrentRoot, fetchData) => () => {
    // Если мы уже в корневом каталоге, ничего не делаем.
    // if (currentRoot === '/') return;

    // Разделяем путь на части и удаляем пустые элементы
    const pathParts = currentRoot.split('/').filter(part => part.length > 0);
    // Удаляем последний элемент пути (текущую директорию)
    pathParts.pop();
    // Обновляем текущий путь
    const newRoot = pathParts.length > 0 ? `/${pathParts.join('/')}` : '/';
    // Устанавливаем новый текущий путь
    setCurrentRoot(newRoot);
    // Получаем данные для новой директории
    fetchData(newRoot);
};
