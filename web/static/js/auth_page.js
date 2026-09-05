document.addEventListener("DOMContentLoaded", () => {
    const tabs = document.querySelectorAll(".auth-tab");
    const panels = Array.from(document.querySelectorAll(".auth-panel"));

    function switchMode(mode) {
        tabs.forEach(t => t.classList.toggle("active", t.dataset.mode === mode));

        const next = mode === "register" ? 1 : 0;

        panels.forEach((p, i) => {
            p.classList.remove("leaving", "out");
            p.classList.toggle("active", i === next);
        });
    }

    tabs.forEach(tab => {
        tab.addEventListener("click", () => switchMode(tab.dataset.mode));
    });
});
