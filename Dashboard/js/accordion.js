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

const items3 = document.querySelectorAll(".accordion-header3");

items3.forEach(header => {
    header.addEventListener("click", () => {

        const item = header.parentElement;

        item.classList.toggle("active");
    });
});

const items4 = document.querySelectorAll(".accordion-header4");

items4.forEach(header => {
    header.addEventListener("click", () => {

        const item = header.parentElement;

        item.classList.toggle("active");
    });
});

const items5 = document.querySelectorAll(".accordion-header5");

items5.forEach(header => {
    header.addEventListener("click", () => {

        const item = header.parentElement;

        item.classList.toggle("active");
    });
});

document.querySelectorAll(".dropdown-button").forEach(button => {
    button.addEventListener("click", () => {
        const dropdown = button.parentElement;
        dropdown.classList.toggle("active");
    });
});

