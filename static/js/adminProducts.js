// Добавляем обработчики событий для кнопок обновления и удаления
window.onload = function () {
    fetchData(); // Загружаем данные при загрузке страницы
}

function fetchData() {
    fetch('/data/products') // Отправляем GET запрос на сервер для получения данных
        .then(response => response.json())
        .then(data => {
            displayData(data); // Отображаем полученные данные на странице
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}