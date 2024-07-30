// Экспортируемая функция navigateToDirectory используется для перехода в указанную директорию
// Параметры:
// - getCurrentRoot: функция для получения текущего пути к директории
// - setCurrentRoot: функция для обновления текущего пути
// - fetchData: функция для получения данных с сервера
export const navigateToDirectory = (getCurrentRoot, setCurrentRoot, fetchAndUpdateTable) => (dirName) => {
    const currentRoot = getCurrentRoot();
    const newRoot = currentRoot === '/' ? `/${dirName}` : `${currentRoot}/${dirName}`;
    setCurrentRoot(newRoot);
    fetchAndUpdateTable();
};

// Экспортируемая функция navigateBack используется для возврата к предыдущей директории
// Параметры:
// - getCurrentRoot: функция для получения текущего пути к директории
// - setCurrentRoot: функция для обновления текущего пути
// - fetchData: функция для получения данных с сервера.
export const navigateBack = (getCurrentRoot, setCurrentRoot, fetchAndUpdateTable,rootDir) => () => {
    const currentRoot = getCurrentRoot();
    if (currentRoot === rootDir) return;

    const pathParts = currentRoot.split('/').filter(part => part.length > 0);
    pathParts.pop();
    const newRoot = pathParts.length > 0 ? `/${pathParts.join('/')}` : '/';
    setCurrentRoot(newRoot);
    fetchAndUpdateTable();
};
