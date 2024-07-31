// configLoader.js

/**
 * Функция для получения конфигурации с сервера
 */
export async function fetchConfig() {
    try {
        const response = await fetch(`/fs`, { method: "GET" });

        if (!response.ok) {
            throw new Error('Ошибка ответа сети');
        }

        const data = await response.json();

        // Получаем нужное значение из JSON и возвращаем его
        return data.root;
    } catch (error) {
        // Обработка ошибок
        console.log('Возникла проблема с операцией получения данных:', error);
        alert('Возникла проблема с операцией получения данных');
    }
}