// configLoader.js

/**
 * Функция для получения конфигурации с сервера
 */
export async function fetchConfig() {
    try {
        const response = await fetch(`/config`, { method: "GET" });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const data = await response.json();

        // Получаем нужное значение из JSON и возвращаем его
        const someValue = data.rootDir; //
        return someValue;
    } catch (error) {
        // Обработка ошибок
        console.error('There has been a problem with your fetch operation:', error);
    }
}
