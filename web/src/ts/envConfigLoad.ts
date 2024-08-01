// Функция для получения конфигурации с сервера
export async function fetchConfig(): Promise<string | undefined> {
    try {
        const response = await fetch(`/fs`, { method: "GET" });

        const data = await response.json();
        // Получаем нужное значение из JSON и возвращаем его
        return data.root;
    } catch (error) {
        // Обработка ошибок
        console.log('Возникла проблема с операцией получения данных:', error);
        alert('Возникла проблема с операцией получения данных');
    }
}
