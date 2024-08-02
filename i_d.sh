npm install --save-dev style-loader css-loader html-webpack-plugin

if [ $? -eq 0 ]; then
    echo "Пакеты успешно установлены."
else
    echo "Произошла ошибка при установке пакетов."
fi
