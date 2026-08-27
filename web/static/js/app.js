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
            const email = document.getElementById("email").value;
            const password = document.getElementById("password").value;
            // const alertBox = document.getElementById("auth-alert");

            // alertBox.classList.add("hidden");

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
                alertBox.textContent = err.message;
                alertBox.classList.remove("hidden");
            });
        });
    }

    // Register Form Handler
    const registerForm = document.getElementById("register-form");
    if (registerForm) {
        registerForm.addEventListener("submit", (e) => {
            e.preventDefault();
            const email = document.getElementById("email").value;
            const firstname = document.getElementById("firstname").value;
            const lastname = document.getElementById("lastname").value;
            const password = document.getElementById("password").value;
            const confirm_password = document.getElementById("confirm_password").value;
            const alertBox = document.getElementById("auth-alert");
            const successBox = document.getElementById("auth-success");

            alertBox.classList.add("hidden");
            successBox.classList.add("hidden");

            if (password !== confirm_password) {
                alertBox.textContent = "Konfirmasi kata sandi tidak cocok";
                alertBox.classList.remove("hidden");
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
                successBox.textContent = "Registrasi sukses! Mengalihkan ke halaman masuk...";
                successBox.classList.remove("hidden");
                setTimeout(() => {
                    window.location.href = "/login";
                }, 1500);
            })
            .catch(err => {
                alertBox.textContent = err.message;
                alertBox.classList.remove("hidden");
            });
        });
    }
});
