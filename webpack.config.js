const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    entry: './web/src/ts/main.ts',          // Входная точка вашего приложения
    output: {
        filename: 'bundle.[contenthash].js',            // Имя выходного файла
        path: path.resolve(__dirname, 'dist') // Директория для выходного файла
    },
    resolve: {
        extensions: ['.ts', '.js', '.css'],  // Расширения файлов, которые Webpack будет обрабатывать
    },
    module: {
        rules: [
            {
                test: /\.ts$/,                // Применение правил к файлам с расширением .ts
                use: 'ts-loader',             // Использование ts-loader для обработки TypeScript
                exclude: /node_modules/,      // Исключение папки node_modules
            },
            {
                test: /\.css$/,               // Применение правил к файлам с расширением .css
                use: ['style-loader', 'css-loader'], // Использование загрузчиков для обработки CSS
            },
        ],
    },
    plugins: [
        new HtmlWebpackPlugin({
            template: './web/src/index.html',    // Шаблон HTML файла
            filename: 'index.html'           // Имя выходного HTML файла
        })
    ],
    mode: 'development',               // Режим сборки (можно установить в 'production' для продакшн)
};
