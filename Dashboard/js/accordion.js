const items = document.querySelectorAll(".accordion-header1");

items.forEach(header => {
    header.addEventListener("click", () => {

        const item = header.parentElement;

        item.classList.toggle("active");

    });
});

const items2 = document.querySelectorAll(".accordion-header2");

items2.forEach(header => {
    header.addEventListener("click", () => {

        const item = header.parentElement;

        item.classList.toggle("active");
    });
});