let entrance=document.querySelector("#alreadyHaveAccount");
let registration=document.querySelector("#newUser");
let loginForm=document.querySelector("#login-form");
let registrationForm=document.querySelector("#registration");
entrance.onclick = function() {
    console.log("Вы нажали на вход");
    registrationForm.style.display = "none";
    loginForm.style.display = "block";
    event.preventDefault();
};
registration.onclick = function() {
    console.log("Вы нажали на регистрацию");
    registrationForm.style.display = "block";
    loginForm.style.display = "none";
    event.preventDefault();
};