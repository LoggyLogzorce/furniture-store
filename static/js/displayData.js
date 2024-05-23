function displayData(data) {
    const tbody = document.getElementById('data-body');
    tbody.innerHTML = ''; // Очищаем содержимое tbody

    // Создаем строку для каждого элемента массива
    data.forEach(item => {
        const tr = document.createElement('tr');
        // Для каждого элемента массива создаем ячейку данных
        Object.keys(item).forEach(key => {
            const td = document.createElement('td');
            td.textContent = item[key]; // Устанавливаем текстовое содержимое ячейки
            tr.appendChild(td); // Добавляем ячейку в строку
        });

        // Добавляем кнопки "Редактировать" и "Удалить" в общий контейнер
        const tdActions = document.createElement('td');
        tdActions.className = 'action-buttons'; // Добавляем класс для стилизации кнопок
        tr.appendChild(tdActions); // Добавляем контейнер в строку

        // Добавляем кнопку "Редактировать" с обработчиком события
        const editButton = document.createElement('button');
        editButton.textContent = 'Редактировать';
        editButton.onclick = function () {
            editRow(tr); // Передаем объект строки и идентификатор
        };
        tdActions.appendChild(editButton); // Добавляем кнопку в контейнер

        // Добавляем кнопку "Удалить" с обработчиком события
        const deleteButton = document.createElement('button');
        deleteButton.textContent = 'Удалить';
        deleteButton.onclick = function () {
            deleteRow(tr); // Передаем объект строки для удаления
        };
        tdActions.appendChild(deleteButton); // Добавляем кнопку в контейнер

        tbody.appendChild(tr); // Добавляем строку в tbody
    });
}

function editRow(row) {
    const tds = row.querySelectorAll('td');
    tds.forEach(td => {
        const currentValue = td.textContent;
        const input = document.createElement('input');
        input.type = 'text';
        input.value = currentValue; // Устанавливаем текущее значение ячейки в поле ввода
        td.textContent = ''; // Очищаем содержимое ячейки
        td.appendChild(input); // Добавляем поле ввода в ячейку
    });

    // Заменяем кнопку "Редактировать" на кнопку "Сохранить"
    const saveButton = document.createElement('button');
    saveButton.textContent = 'Сохранить';
    saveButton.onclick = function () {
        saveRow(row); // Передаем объект строки и идентификатор для сохранения
    };
    const tdActions = row.querySelector('td:last-child');
    tdActions.textContent = ''; // Очищаем содержимое ячейки с кнопкой
    tdActions.appendChild(saveButton); // Добавляем кнопку "Сохранить"
}