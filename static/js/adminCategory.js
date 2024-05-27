// Добавляем обработчики событий для кнопок обновления и удаления
window.onload = function () {
    fetchData(); // Загружаем данные при загрузке страницы
}

function fetchData() {
    fetch('/data/categories') // Отправляем GET запрос на сервер для получения данных
        .then(response => response.json())
        .then(data => {
            displayData(data); // Отображаем полученные данные на странице
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
}

function deleteRow(row) {
    const rowData = {}; // Создаем объект с данными строки
    row.querySelectorAll('td').forEach((td, index) => {
        const columnName = document.querySelector(`#data-table th:nth-child(${index + 1})`).textContent;
        rowData[columnName] = td.textContent; // Заполняем объект данными из ячеек строки
    });

    // Отправляем запрос DELETE на сервер с данными строки для удаления
    fetch('/ad/delete/category', {
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

// Функция для сохранения изменений строки
function saveRow(row) {
    const inputs = row.querySelectorAll('input');
    const updatedData = {}; // Создаем объект с обновленными данными
    inputs.forEach(input => {
        const columnIndex = input.parentElement.cellIndex;
        const columnName = document.querySelector(`#data-table th:nth-child(${columnIndex + 1})`).textContent;
        updatedData[columnName] = input.value; // Заполняем объект значениями из полей ввода
    });

    // Отправляем запрос PUT на сервер с обновленными данными
    fetch('/ad/update/category', { // Используем переданный идентификатор
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(updatedData)
    })
        .then(response => response.json())
        .then(data => {
            console.log('Row updated:', data);
            fetchData(); // Обновляем данные

            // Восстанавливаем кнопку "Удалить"
            const tdActions = row.querySelector('.action-buttons');
            const deleteButton = document.createElement('button');
            deleteButton.textContent = 'Удалить';
            deleteButton.onclick = function () {
                deleteRow(row); // Передаем объект строки для удаления
            };
            tdActions.appendChild(deleteButton); // Добавляем кнопку "Удалить"
        })
        .catch(error => {
            console.error('Error updating row:', error);
        })
        .finally(() => {
            // Удаляем все поля ввода из ячеек
            inputs.forEach(input => {
                const td = input.parentElement;
                td.textContent = input.value;
            });

            // Заменяем кнопку "Сохранить" на кнопку "Редактировать"
            const editButton = document.createElement('button');
            editButton.textContent = 'Редактировать';
            editButton.onclick = function () {
                editRow(row); // Передаем объект строки и идентификатор
            };
            const tdActions = row.querySelector('.action-buttons');
            tdActions.textContent = ''; // Очищаем содержимое ячейки с кнопкой
            tdActions.appendChild(editButton); // Добавляем кнопку "Редактировать"
        });
}

function addRow() {
    const category = document.getElementById('category').value;

    // Отправляем запрос POST на сервер для добавления записи
    fetch('/ad/add/category', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({category})
    })
        .then(response => response.json())
        .then(data => {
            console.log('Row added:', data);
            fetchData(); // Обновляем данные после добавления
        })
        .catch(error => {
            console.error('Error adding row:', error);
        });
}