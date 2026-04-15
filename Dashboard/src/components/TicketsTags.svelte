<script>
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



</script>

{@html `<div class="accordion1">
        <div class="accordion-item1">
            <div class="accordion-header1">
                <h3>Tag Settings</h3>
            </div>
                <div class="accordion-content1" id="content1">
                    <div class="ticket-overview" id="check-button1">
                        <h3 class="accordion1-texts">Enable Tags</h3>
                        <input type="checkbox" id="check1">
                        <label for="check1" class="button1"></label>
                    </div>
                    <div class="ticket-overview" id="check-button2">
                        <h3 class="accordion1-texts">Custom Tags</h3>
                        <input type="checkbox" id="check2">
                        <label for="check2" class="button2"></label>
                    </div>
                    <div class="ticket-overview" id="check-button3">
                        <h3 class="accordion1-texts">Tag Colors</h3>
                        <input type="checkbox" id="check3">
                        <label for="check3" class="button3"></label>
                    </div>
                    <div class="ticket-overview" id="check-button4">
                        <h3 class="accordion1-texts">Tag Permissions</h3>
                        <input type="checkbox" id="check4">
                        <label for="check4" class="button4"></label>
                    </div>
                </div>
            </div>
        </div>
    <div class="accordion2">
        <div class="accordion-item2">
            <div class="accordion-header2">
                <h3>Tag Management</h3>
            </div>
                <div class="accordion-content2" id="content2">
                    <div class="ticket-overview" id="check-button5">
                        <h3 class="accordion1-texts">Create Tags</h3>
                        <input type="checkbox" id="check5">
                        <label for="check5" class="button5"></label>
                    </div>
                    <div class="ticket-overview" id="check-button6">
                        <h3 class="accordion1-texts">Edit Tags</h3>
                        <input type="checkbox" id="check6">
                        <label for="check6" class="button6"></label>
                    </div>
                    <div class="ticket-overview" id="check-button7">
                        <h3 class="accordion1-texts">Delete Tags</h3>
                        <input type="checkbox" id="check7">
                        <label for="check7" class="button7"></label>
                    </div>
                    <div class="ticket-overview" id="check-button8">
                        <h3 class="accordion1-texts">Tag Categories</h3>
                        <input type="checkbox" id="check8">
                        <label for="check8" class="button8"></label>
                    </div>
                </div>
            </div>
        </div>
    </div>`}

<style>
body{
margin:0;
font-family:Arial, Helvetica, sans-serif;
background:#f4f6fb;
}

/* LAYOUT */

.dashboard{
display:flex;
height:100vh;
}


/* LEFT SIDEBAR */

.sidebar{
width:230px;
background:#111827;
color:white;

display:flex;
flex-direction:column;
padding:20px;
}

.logo{
display:flex;
align-items:center;
gap:10px;
margin-bottom:40px;
}

.logo-img{
width:40px;
height:40px;
border-radius:50%;
}

.sidebar nav{
display:flex;
flex-direction:column;
gap:15px;
}

.sidebar a{
color:#c7c7c7;
text-decoration:none;
padding:10px;
border-radius:8px;
}

.sidebar a:hover{
background:#1f2937;
color:white;
}

#Tickets{
background:#1f2937;
}

#Overview{
background:#1f2937;
}

#Servers{
background:#1f2937;
}

#Giveaways{
background:#1f2937;
}

#Moderation{
background:#1f2937;
}

#Leveling{
background:#1f2937;
}

#Dashboard{
background:#1f2937;
}

/* RIGHT SIDEBAR */

.right-sidebar{
position:fixed;
right: 10px;
top:200px;
width:230px;
height: 400px;
background:#111827;
color:white;
display:flex;
flex-direction:column;
padding:20px;
border-radius: 10px;
}

.right-sidebar nav{
display:flex;
flex-direction:column;
gap:15px;
}

.right-sidebar a{
color:#c7c7c7;
text-decoration:none;
padding:10px;
border-radius:8px;
}

.right-sidebar a:hover{
background:#1f2937;
color:white;
}

#Settings{
background:#1f2937;
}
/* MAIN */

.main{
flex:1;
display:flex;
flex-direction:column; 
}

.content{
    display: flex;
    flex-direction: column;
    gap: 10px; 
}

/* ACCORDION */

.accordion1,
.accordion2,
.accordion3,
.accordion4{
width:100%;
color:rgb(0, 0, 0);
}

.accordion-item1,
.accordion-item2,
.accordion-item3,
.accordion-item4{
border-bottom:1px solid rgba(255,255,255,0.1);
}

.accordion-header1,
.accordion-header2,
.accordion-header3,
.accordion-header4{
cursor:pointer;
padding:15px;
font-weight:bold;
display:flex;
flex-direction: row;
align-items:center;
justify-content:space-between;
}

.accordion-header1:hover,
.accordion-header2:hover,
.accordion-header3:hover,
.accordion-header4:hover{
background:rgba(255,255,255,0.05);
}

.accordion-content1,
.accordion-content2,
.accordion-content3,
.accordion-content4{
max-height:0;
overflow:hidden;
transition:all 0.3s ease;
padding:0 15px;
display: flex;
flex-direction: row;
}

.accordion-content4 {
    display: flex;
    flex-direction: row;
    gap: 20px;
}

.accordion-item.active .accordion-content {
    max-height: 500px;
}

.accordion1 .accordion-content1 .accordion1-texts,
.accordion2 .accordion-content2 .accordion2-texts,
.accordion3 .accordion-content3 .accordion3-texts,
.accordion4 .accordion-content4 .accordion4-texts {
    display: flex;
    align-items: center;
    top: 100px;
}

/* open state */

.accordion-item1.active .accordion-content1,
.accordion-item2.active .accordion-content2,
.accordion-item3.active .accordion-content3,
.accordion-item4.active .accordion-content4{
max-height:500px; /* riittää */
padding:15px;
}

.content1,
.content2,
.content3,
.content4{
display:flex;
flex-direction: row;
gap: 10px;
}

/* CHECK BUTTONS */

#check-button1,
#check-button2,
#check-button3,
#check-button4,
#check-button5,
#check-button6,
#check-button7,
#check-button8,
#check-button9,
#check-button10,
#check-button11,
#check-button12,
#check-button13,
#check-button14,
#check-button15,
#check-button16{
    width: 100px;
    height: 100px;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: #f4f6fb;
}

/* Toggle buttons */

.button1,
.button2,
.button3,
.button4,
.button5,
.button6,
.button7,
.button8,
.button9,
.button10,
.button11,
.button12,
.button13,
.button14,
.button15,
.button16{
    width: 60px;
    height: 30px;
    background-color: #c7c7c7;
    border-radius: 15px;
    cursor: pointer;
    position: relative;
}

/* Toggle buttons */

.button1::before,
.button2::before,
.button3::before,
.button4::before,
.button5::before,
.button6::before,
.button7::before,
.button8::before,
.button9::before,
.button10::before,
.button11::before,
.button12::before,
.button13::before,
.button14::before,
.button15::before,
.button16::before{
    content: "";
    width: 30px;
    height: 30px;
    background-color: white;
    border-radius: 50%;
    position: absolute;
    top: 0;
    left: 0;
    transition: all 0.3s ease;
}

/* Toggle buttons checked state */

input:checked + .button1::before,
input:checked + .button2::before,
input:checked + .button3::before,
input:checked + .button4::before,
input:checked + .button5::before,
input:checked + .button6::before,
input:checked + .button7::before,
input:checked + .button8::before,
input:checked + .button9::before,
input:checked + .button10::before,
input:checked + .button11::before,
input:checked + .button12::before,
input:checked + .button13::before,
input:checked + .button14::before,
input:checked + .button15::before,
input:checked + .button16::before{
    transform: translateX(30px);
    background-color: #1f2937;
}

/* Accordion container */

.accordion1,
.accordion2,
.accordion3,
.accordion4{
    width: 100%;
    color:rgb(0, 0, 0);
    display: flex;
    flex-direction: row;
    gap: 10px;
    bottom: 100px;
}

input{
    display: none;
}

/* DROPDOWN MENUS */

.dropdown button{
    background-color: #111827;
    color: #c7c7c7;
    padding: 10px 5px;
    border: none;
    cursor: pointer;
    margin: 15px;
    margin-top: 45px;
    border-radius: 10px;
}

.dropdown a{
    display: block;
    color: #111827;
    text-decoration: none;

}

.dropdown-content {
    position: absolute;
    margin-top: 10px;
    background: white;
    border-radius: 10px;
    padding: 5px 0;
    min-width: 150px;
    box-shadow: 0 8px 20px rgba(0,0,0,0.15);
}

.dropdown-content a {
    padding: 10px 15px;
    display: block;
    color: #111827;
    font-size: 14px;
    transition: all 0.2s ease;
}

.dropdown-content a:hover {
    background: #f3f4f6;
    color: #111827;
}

.dropdown-content a:first-child {
    border-radius: 10px 10px 0 0;
}

.dropdown-content a:last-child {
    border-radius: 0 0 10px 10px;
}

.dropdown:hover button{
    background-color: #1f2937;
}

.dropdown-content {
    display: none;
}

.dropdown.active .dropdown-content {
    display: block;
}

.arrow {
    margin-left: 8px;
    transition: transform 0.3s ease;
}

.dropdown.active .arrow {
    transform: rotate(180deg);
}

/* TEXTS */

.ticket-overview {
    display: flex;
    flex-direction: column; /* tämä tekee tekstin napin päälle */
    align-items: center;
    justify-content: center;
    gap: 10px;
}

.accordion1-texts {
    font-size: 14px;
    text-align: center;
    margin: 0;
}

/* TOPBAR */

.topbar{
display:flex;
justify-content:space-between;
align-items:center;

padding:20px 40px;
background:white;
border-bottom:1px solid #eee;
}

.user img{
width:40px;
height:40px;
border-radius:50%;
}

</style>