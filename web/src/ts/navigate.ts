// Функция для перехода в указанную директорию
// - getCurrentRoot: функция для получения текущего пути к директории
// - setCurrentRoot: функция для обновления текущего пути
// - fetchData: функция для получения данных с сервера
export const navigateToDirectory = (
    getCurrentRoot: () => string,
    setCurrentRoot: (newRoot: string) => void,
    fetchAndUpdateTable: () => void
) => (dirName: string) => {
    const currentRoot = getCurrentRoot();
    const newRoot = currentRoot === '/' ? `/${dirName}` : `${currentRoot}/${dirName}`;
    setCurrentRoot(newRoot);
    fetchAndUpdateTable();
};

// Функция для возврата к предыдущей директории
// - getCurrentRoot: функция для получения текущего пути к директории
// - setCurrentRoot: функция для обновления текущего пути
// - fetchData: функция для получения данных с сервера
export const navigateBack = (
    getCurrentRoot: () => string,
    setCurrentRoot: (newRoot: string) => void,
    fetchAndUpdateTable: () => void,
    rootDir: string
) => () => {
    const currentRoot = getCurrentRoot();
    if (currentRoot === rootDir){
        return alert("Вы находитесь в корневой директории");
    }
    // Разбиваем на массив [value,value,..]
    const pathParts = currentRoot.split('/').filter(part => part.length > 0);
    pathParts.pop();
    // Если массив не пустой, то объединяем его в строку
    const newRoot = pathParts.length > 0 ? `/${pathParts.join('/')}` : '/';
    setCurrentRoot(newRoot);
    fetchAndUpdateTable();
};

