function logout(){
    fetch('/api/v1/auth/logout',{
        method: 'POST',
        headers: {'Content-Type' : 'application/json'}
    })
    .then(res => res.json())
    .then(data => {
        // Clear cookie client-side just in case
        document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
        window.location.href = '/login';
    })
    .catch(err => {
        console.error('Logout error:', err);
        window.location.href = '/';
    });
}

// Highlight active sidebar item
document.addEventListener("DOMContentLoaded", () => {
    const currentPath = window.location.pathname;
    const navItems = document.querySelectorAll(".nav-item");
    navItems.forEach(item => {
        if (item.getAttribute("href") === currentPath) {
            item.classList.add("active");
        }
    });

    // Login Form Handler
    const loginForm = document.getElementById("login-form");
    if (loginForm) {
        loginForm.addEventListener("submit", (e) => {
            e.preventDefault();
            const email = document.getElementById("login-email").value;
            const password = document.getElementById("login-password").value;

            fetch("/api/v1/auth/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, password })
            })
            .then(async res => {
                const data = await res.json();
                if (!res.ok) {
                    throw new Error(data.message || "Email atau password salah");
                }
                return data;
            })
            .then(data => {
                window.location.href = data.data.redirect || "/dashboard";
            })
            .catch(err => {
                alert(err.message);
            });
        });
    }

    // Register Form Handler
    const registerForm = document.getElementById("register-form");
    if (registerForm) {
        registerForm.addEventListener("submit", (e) => {
            e.preventDefault();
            const email = document.getElementById("reg-email").value;
            const firstname = document.getElementById("reg-firstname").value;
            const lastname = document.getElementById("reg-lastname").value;
            const password = document.getElementById("reg-password").value;
            const confirm_password = document.getElementById("reg-confirm").value;

            if (password !== confirm_password) {
                alert("Konfirmasi kata sandi tidak cocok");
                return;
            }

            fetch("/api/v1/auth/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, firstname, lastname, password, confirm_password })
            })
            .then(async res => {
                const data = await res.json();
                if (!res.ok) {
                    throw new Error(data.message || "Gagal mendaftarkan akun");
                }
                return data;
            })
            .then(data => {
                // auto login setelah register sukses
                return fetch("/api/v1/auth/login", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ email, password })
                }).then(async res => {
                    const loginData = await res.json();
                    if (!res.ok) {
                        throw new Error(loginData.message || "Login otomatis gagal");
                    }
                    return loginData;
                });
            })
            .then(loginData => {
                window.location.href = loginData.data.redirect || "/dashboard";
            })
            .catch(err => {
                alert(err.message);
            });
        });
    }
});
