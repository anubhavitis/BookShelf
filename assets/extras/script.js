const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');

signUpButton.addEventListener('click', () => {
	container.classList.add("right-panel-active");
});

signInButton.addEventListener('click', () => {
	container.classList.remove("right-panel-active");
});

var password = document.getElementById("password")
  , c_password = document.getElementById("c_password");

function validatePassword(){
  if(password.value !== c_password.value) {
    c_password.setCustomValidity("Passwords Don't Match");
  } else {
    c_password.setCustomValidity('');
  }
}

password.onchange = validatePassword;
c_password.onkeyup = validatePassword;
