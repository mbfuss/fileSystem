/**
 * Функция для получения конфигурации с сервера
 */
export async function fetchConfig() {
        const response = await fetch(`/fs`, { method: "GET" });

        let data = await response.json();
        // Получаем нужное значение из JSON и возвращаем его
        return data.root;


}
