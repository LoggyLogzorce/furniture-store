// Добавляем обработчики событий для кнопок обновления и удаления
window.onload = function () {
    fetchData(); // Загружаем данные при загрузке страницы
}

function fetchData() {
    fetch('/data/reviews') // Отправляем GET запрос на сервер для получения данных
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
    tbody.innerHTML = ''; // Очищаем таблицу перед добавлением новых данных

    data.forEach(item => {
        const row = document.createElement('tr');

        // Добавляем ячейки с данными
        for (const key in item) {
            const cell = document.createElement('td');
            cell.textContent = item[key];
            row.appendChild(cell);
        }

        // Добавляем ячейку с кнопками действий
        const actionCell = document.createElement('td');
        actionCell.classList.add('action-buttons');

        const deleteButton = document.createElement('button');
        deleteButton.textContent = 'Удалить';
        deleteButton.onclick = function () {
            deleteRow(row);
        };
        actionCell.appendChild(deleteButton);

        row.appendChild(actionCell);
        tbody.appendChild(row);
    });
}

function deleteRow(row) {
    const rowData = {}; // Создаем объект с данными строки
    row.querySelectorAll('td').forEach((td, index) => {
        const columnName = document.querySelector(`#data-table th:nth-child(${index + 1})`).textContent;
        rowData[columnName] = td.textContent; // Заполняем объект данными из ячеек строки
    });

    // Отправляем запрос DELETE на сервер с данными строки для удаления
    fetch('/ad/delete/review', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(rowData)
    })
        .then(response => response.json())
        .then(data => {
            console.log('Row deleted:', data);
            fetchData(); // Обновляем данные после удаления
        })
        .catch(error => {
            console.error('Error deleting row:', error);
        });
}