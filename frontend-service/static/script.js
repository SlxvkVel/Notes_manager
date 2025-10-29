let currentUsername = '';
let currentEditingNoteId = null;

function toggleForm(formId) {
    const forms = document.querySelectorAll('.auth-form');
    forms.forEach(form => {
        form.style.display = 'none';
    });
    
    const form = document.getElementById(formId);
    if (form.style.display === 'block') {
        form.style.display = 'none';
    } else {
        form.style.display = 'block';
    }
    
    document.getElementById('userInfo').style.display = 'none';
}

function showUserInfo(username) {
    currentUsername = username;
    document.getElementById('currentUser').textContent = username;
    document.getElementById('userInfo').style.display = 'block';
    
    const forms = document.querySelectorAll('.auth-form');
    forms.forEach(form => {
        form.style.display = 'none';
    });
}

async function register() {
    const username = document.getElementById('regUsername').value;
    const email = document.getElementById('regEmail').value;
    const password = document.getElementById('regPassword').value;
    
    if (!username || !email || !password) {
        alert('Заполните все поля для регистрации');
        return;
    }
    
    try {
        const response = await fetch('/api/auth/register', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({
                username: username,
                email: email,
                password: password
            })
        });
        
        const responseText = await response.text();
        
        if (!response.ok) {
            let errorData;
            try {
                errorData = JSON.parse(responseText);
            } catch (e) {
                throw new Error(responseText || `HTTP error! status: ${response.status}`);
            }
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }
        
        const data = JSON.parse(responseText);
        
        document.getElementById('regUsername').value = '';
        document.getElementById('regEmail').value = '';
        document.getElementById('regPassword').value = '';
        
        toggleForm('registerForm');
        showUserInfo(username);
        alert(`✅ Регистрация успешна!\n\nДобро пожаловать, ${username}!`);
        
    } catch (error) {
        let errorMessage = 'Ошибка регистрации';
        
        if (error.message.includes('already exists') || error.message.includes('duplicate')) {
            errorMessage = '❌ Пользователь с таким email или именем уже существует';
        } else if (error.message.includes('400')) {
            errorMessage = '❌ Неверные данные для регистрации';
        } else if (error.message.includes('500')) {
            errorMessage = '❌ Ошибка сервера, попробуйте позже';
        } else {
            errorMessage = '❌ ' + error.message;
        }
        
        alert(errorMessage);
    }
}

async function login() {
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;
    
    if (!email || !password) {
        alert('Заполните email и пароль');
        return;
    }
    
    try {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({
                email: email,
                password: password
            })
        });
        
        const responseText = await response.text();
        
        if (!response.ok) {
            let errorData;
            try {
                errorData = JSON.parse(responseText);
            } catch (e) {
                throw new Error(responseText || `HTTP error! status: ${response.status}`);
            }
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }
        
        const data = JSON.parse(responseText);
        
        document.getElementById('loginEmail').value = '';
        document.getElementById('loginPassword').value = '';
        
        toggleForm('loginForm');
        showUserInfo(data.username || email.split('@')[0]);
        alert('✅ Вход выполнен успешно!');
        
        getNotes();
        
    } catch (error) {
        let errorMessage = 'Ошибка входа';
        
        if (error.message.includes('401') || error.message.includes('invalid credentials')) {
            errorMessage = '❌ Неверный email или пароль';
        } else if (error.message.includes('400')) {
            errorMessage = '❌ Заполните все поля';
        } else if (error.message.includes('500')) {
            errorMessage = '❌ Ошибка сервера, попробуйте позже';
        } else if (error.message.includes('user not found')) {
            errorMessage = '❌ Пользователь не найден';
        } else {
            errorMessage = '❌ ' + error.message;
        }
        
        alert(errorMessage);
    }
}

async function createNote() {
    const title = document.getElementById('title').value;
    const content = document.getElementById('content').value;
    
    if (!title || !content) {
        alert('Заполните заголовок и текст заметки');
        return;
    }
    
    try {
        const response = await fetch('/api/notes', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({
                title: title,
                content: content
            })
        });
        
        const responseText = await response.text();
        
        if (!response.ok) {
            let errorData;
            try {
                errorData = JSON.parse(responseText);
            } catch (e) {
                throw new Error(responseText || `HTTP error! status: ${response.status}`);
            }
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }
        
        const data = JSON.parse(responseText);
        alert('✅ Заметка создана!');
        
        document.getElementById('title').value = '';
        document.getElementById('content').value = '';
        
        getNotes();
    } catch (error) {
        let errorMessage = 'Ошибка создания заметки';
        
        if (error.message.includes('401') || error.message.includes('unauthorized')) {
            errorMessage = '❌ Для создания заметки необходимо войти в систему';
        } else if (error.message.includes('400')) {
            errorMessage = '❌ Неверные данные для заметки';
        } else if (error.message.includes('500')) {
            errorMessage = '❌ Ошибка сервера, попробуйте позже';
        } else {
            errorMessage = '❌ ' + error.message;
        }
        
        alert(errorMessage);
    }
}

async function deleteNote(noteId) {
    if (!confirm('Вы уверены, что хотите удалить эту заметку?')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/notes/delete?id=${noteId}`, {
            method: 'DELETE',
            credentials: 'include'
        });
        
        const responseText = await response.text();
        
        if (!response.ok) {
            let errorData;
            try {
                errorData = JSON.parse(responseText);
            } catch (e) {
                throw new Error(responseText || `HTTP error! status: ${response.status}`);
            }
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }
        
        alert('✅ Заметка удалена!');
        getNotes();
    } catch (error) {
        let errorMessage = 'Ошибка удаления заметки';
        
        if (error.message.includes('401') || error.message.includes('unauthorized')) {
            errorMessage = '❌ Для удаления заметки необходимо войти в систему';
        } else if (error.message.includes('404')) {
            errorMessage = '❌ Заметка не найдена';
        } else if (error.message.includes('500')) {
            errorMessage = '❌ Ошибка сервера, попробуйте позже';
        } else {
            errorMessage = '❌ ' + error.message;
        }
        
        alert(errorMessage);
    }
}

function openEditModal(noteId, title, content) {
    currentEditingNoteId = noteId;
    document.getElementById('editTitle').value = title;
    document.getElementById('editContent').value = content;
    document.getElementById('editModal').style.display = 'block';
    document.getElementById('modalBackdrop').style.display = 'block';
}

function closeEditModal() {
    document.getElementById('editModal').style.display = 'none';
    document.getElementById('modalBackdrop').style.display = 'none';
    currentEditingNoteId = null;
}

async function updateNote() {
    const title = document.getElementById('editTitle').value;
    const content = document.getElementById('editContent').value;
    
    if (!title || !content) {
        alert('Заполните заголовок и текст заметки');
        return;
    }
    
    try {
        const response = await fetch(`/api/notes/update?id=${currentEditingNoteId}`, {
            method: 'PUT',
            headers: {'Content-Type': 'application/json'},
            credentials: 'include',
            body: JSON.stringify({
                title: title,
                content: content
            })
        });
        
        const responseText = await response.text();
        
        if (!response.ok) {
            let errorData;
            try {
                errorData = JSON.parse(responseText);
            } catch (e) {
                throw new Error(responseText || `HTTP error! status: ${response.status}`);
            }
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }
        
        alert('✅ Заметка обновлена!');
        closeEditModal();
        getNotes();
    } catch (error) {
        let errorMessage = 'Ошибка обновления заметки';
        
        if (error.message.includes('401') || error.message.includes('unauthorized')) {
            errorMessage = '❌ Для редактирования заметки необходимо войти в систему';
        } else if (error.message.includes('404')) {
            errorMessage = '❌ Заметка не найдена';
        } else if (error.message.includes('500')) {
            errorMessage = '❌ Ошибка сервера, попробуйте позже';
        } else {
            errorMessage = '❌ ' + error.message;
        }
        
        alert(errorMessage);
    }
}

async function getNotes() {
    try {
        const response = await fetch('/api/notes/list', {
            credentials: 'include'
        });
        
        const responseText = await response.text();
        
        if (!response.ok) {
            let errorData;
            try {
                errorData = JSON.parse(responseText);
            } catch (e) {
                throw new Error(responseText || `HTTP error! status: ${response.status}`);
            }
            throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
        }
        
        const data = JSON.parse(responseText);
        displayNotes(data.notes);
    } catch (error) {
        let errorMessage = 'Ошибка загрузки заметок';
        
        if (error.message.includes('401') || error.message.includes('unauthorized')) {
            errorMessage = '❌ Для просмотра заметок необходимо войти в систему';
        } else {
            errorMessage = '❌ ' + error.message;
        }
        
        document.getElementById('notes').innerHTML = `<div class="error-message">${errorMessage}</div>`;
    }
}

function displayNotes(notes) {
    const container = document.getElementById('notes');
    if (!notes || notes.length === 0) {
        container.innerHTML = '<p class="no-notes">Нет заметок</p>';
        return;
    }
    
    container.innerHTML = notes.map(note => `
        <div class="note">
            <div class="note-actions">
                <button class="edit-btn" onclick="openEditModal(${note.id}, '${note.title.replace(/'/g, "\\'")}', '${note.content.replace(/'/g, "\\'")}')">✏️</button>
                <button class="delete-btn" onclick="deleteNote(${note.id})">🗑️</button>
            </div>
            <h4>${note.title}</h4>
            <p>${note.content}</p>
            <small>Создано: ${new Date(note.created_at).toLocaleString()}</small>
        </div>
    `).join('');
}

async function logout() {
    try {
        await fetch('/api/auth/logout', {
            method: 'POST',
            credentials: 'include'
        });
        alert('✅ Выход выполнен');
        document.getElementById('notes').innerHTML = '';
        document.getElementById('userInfo').style.display = 'none';
        currentUsername = '';
    } catch (error) {
        alert('❌ Ошибка выхода: ' + error);
    }
}

window.onload = getNotes;