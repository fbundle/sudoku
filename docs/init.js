import {game} from "./sketch.js";
import {httpPost} from "./http_transport.js";

function setup(response) {
    game.key = response.key;
    document.getElementById("key").value = response.key;
    document.getElementById("new_board").style.display = "none";
    document.getElementById("new_random").style.display = "none";
    document.getElementById("login_key").style.display = "none";
}

function login_random() {
    httpPost("api/new", {}, setup);
}

function login_key() {
    const textkey = document.getElementById("key");
    httpPost("api/login", {key: textkey.value}, setup);
}

function login_board() {
    const textboard = document.getElementById("board");
    if (textboard.value.length === 0) {
        textboard.value = new Array(81 + 1).join("0");
    }
    httpPost("api/new", {board: textboard.value}, setup);
}

document.querySelector("#new_random").addEventListener("click", login_random);
document.querySelector("#new_board").addEventListener("click", login_board);
document.querySelector("#login_key").addEventListener("click", login_key);
