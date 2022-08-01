const login = document.getElementById("login")
const signUp = document.getElementById("signUp")

// """"""""get inputs""""""""
const emailLog = document.getElementById("emailLog")
const passwordLog = document.getElementById("passwordLog")
const email = document.getElementById("email")
const number = document.getElementById("number")
const password = document.getElementById("password")

login.onsubmit = evt =>{
    preventDefault()
    var data = JSON.stringify({
        "email":emailLog.value,
        "password":passwordLog.value
    });
    console.log(passwordLog.value);
    var xhr = new XMLHttpRequest();
    xhr.withCredentials = true;

    xhr.addEventListener("readystatechange", function() {
    if(this.readyState === 4) {
        console.log(this.responseText);
    }
    });

    xhr.open("POST", "https://127.0.0.1:4400/login");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.send(data);
}

signUp.onsubmit = evt =>{
    // preventDefault()
    var data = JSON.stringify({
        "email":email.value,
        "password":password.value,
        "phone":number.value
    });
    console.log(passwordLog.value);
    var xhr = new XMLHttpRequest();
    xhr.withCredentials = true;

    xhr.addEventListener("readystatechange", function() {
    if(this.readyState === 4) {
        console.log(this.responseText);
    }
    });

    xhr.open("POST", "https://127.0.0.1:4400/signup");
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.send(data);
}