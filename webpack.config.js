const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    entry: './src/main.js', // Входная точка вашего приложения
    output: {
        filename: 'bundle.js', // Имя выходного файла
        path: path.resolve(__dirname, 'dist'), // Папка для выходных файлов
        clean: true // Очистка папки dist перед каждым билдом
    },
    module: {
        rules: [
            {
                test: /\.css$/, // Регулярное выражение для поиска CSS файлов
                use: ['style-loader', 'css-loader'] // Загрузчики для обработки CSS файлов
            },
            {
                test: /\.html$/, // Регулярное выражение для поиска HTML файлов
                use: ['html-loader'] // Загрузчик для обработки HTML файлов
            }
        ]
    },
    plugins: [
        new HtmlWebpackPlugin({
            template: './src/index.html', // Шаблон HTML файла
            filename: 'index.html' // Имя выходного HTML файла
        })
    ],
    devServer: {
        compress: true,
        port: 9000
    },
    mode: 'development' // Режим сборки (development или production)
};
