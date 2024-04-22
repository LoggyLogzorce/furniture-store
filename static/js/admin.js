window.onload = function() {
    fetchData(); // Загружаем данные при загрузке страницы
}

function fetchData() {
    fetch('/data') // Отправляем GET запрос на сервер для получения данных
        .then(response => response.json())
        .then(data => {
            displayData(data); // Отображаем полученные данные на странице
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}

function displayData(data) {
    const tbody = document.getElementById('data-body');
    tbody.innerHTML = ''; // Очищаем содержимое tbody

    // Создаем строку для каждого элемента массива
    data.forEach(item => {
        const tr = document.createElement('tr');
        // Для каждого элемента массива создаем ячейку данных
        Object.values(item).forEach(value => {
            const td = document.createElement('td');
            td.textContent = value; // Устанавливаем текстовое содержимое ячейки
            tr.appendChild(td); // Добавляем ячейку в строку
        });
        tbody.appendChild(tr); // Добавляем строку в tbody
    });
}

