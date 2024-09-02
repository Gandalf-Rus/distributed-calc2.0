// Функция для переключения между окнами входа и регистрации
function toggleView(show, hide, activeBtn, inactiveBtn) {
    document.getElementById(show).style.display = 'block';
    document.getElementById(hide).style.display = 'none';
    document.getElementById(activeBtn).classList.add('active');
    document.getElementById(inactiveBtn).classList.remove('active');
}

function showLogin() {
    toggleView('login', 'register', 'switchToLogin', 'switchToRegister');
}

function showRegister() {
    toggleView('register', 'login', 'switchToRegister', 'switchToLogin');
}

async function login() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const url = 'http://localhost:8080/api/login';
    const data = {
        name: username,
        password: password
    };

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            const result = await response.json();
            const token = result.jwt_token
            localStorage.setItem('authToken', token);
            console.log('Успешный вход:', result);
            showLoggedInState();
        } else {
            console.log('Ошибка:', response.statusText, response);
        }
    } catch (error) {
        console.error('Ошибка сети:', error);
    }
}

async function register() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const url = 'http://localhost:8080/api/register';
    const data = {
        name: username,
        password: password
    };

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            const result = await response.json();

            console.log('Успешная регистрация:', result);

        } else {
            console.log('Ошибка регистрации:', response.statusText);
        }
    } catch (error) {
        console.error('Ошибка сети:', error);
    }
}

function showLoggedInState() {
    document.getElementById('loginPage').style.display = 'none';
    document.getElementById('mainPage').style.display = 'block';
}

async function postExpression() {
    const exitId = document.getElementById('exitId').value;
    const expression = document.getElementById('expressionInput').value;
    const token = localStorage.getItem('authToken'); // Получение токена из localStorage

    if (!token) {
        console.log('Токен не найден. Пожалуйста, войдите или зарегистрируйтесь.');
        return;
    }

    const url = 'http://localhost:8080/api/expression';
    const data = {
        exit_id: exitId,
        expression: expression
    };

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`, // Добавление токена в заголовок Authorization
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Выражение успешно отправлено:', result);
        } else {
            console.log('Ошибка:', response.statusText);
        }
    } catch (error) {
        console.error('Ошибка сети:', error);
    }
}


async function getExpressions() {
    const token = localStorage.getItem('authToken'); // Получение токена из localStorage

    if (!token) {
        console.log('Токен не найден. Пожалуйста, войдите или зарегистрируйтесь.');
        return;
    }

    const url = 'http://localhost:8080/api/get_expressions';

    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`, // Добавление токена в заголовок Authorization
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            const result = await response.json();
            console.log('Полученные выражения:', result);

            // Пример вывода выражений в HTML
            const output = document.getElementById('expressionsOutput');
            output.textContent = JSON.stringify(result);
        } else {
            console.log('Ошибка:', response.statusText);
        }
    } catch (error) {
        console.error('Ошибка сети:', error);
    }
}


const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');

function switchTheme(e) {
    if (e.target.checked) {
        document.documentElement.setAttribute('class', 'dark-theme');
        localStorage.setItem('theme', 'dark');
    }
    else {
        document.documentElement.setAttribute('class', '');
        localStorage.setItem('theme', 'light');
    }    
}

toggleSwitch.addEventListener('change', switchTheme, false);

// Проверяем текущую тему при загрузке страницы
const currentTheme = localStorage.getItem('theme');
if (currentTheme) {
    document.documentElement.setAttribute('class', currentTheme === 'dark' ? 'dark-theme' : '');
    if (currentTheme === 'dark') {
        toggleSwitch.checked = true;
    }
}