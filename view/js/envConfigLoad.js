// configLoader.js

/**
 * Функция для получения конфигурации с сервера
 */
export async function fetchConfig() {
        const response = await fetch(`/fs`, { method: "GET" });
        const data = await response.json();
        // Получаем нужное значение из JSON и возвращаем его
        return data.root;
}
