const cityTemplate = document.getElementById('city-template');

let citiesAjaxContainer = document.getElementById('citiesAjax');
let citiesContainer = document.getElementById('cities');
let citiesAsyncContainer = document.getElementById('citiesAsync');
let citiesWs = document.getElementById('citiesWs');

let wsDuration = performance.now();
let ajaxDuration = performance.now();

const searchCityUrl = './api/city/';
const searchUrl = './api/cities/';
const searchUrlAsync = './api/citiesAsync/';
const searchUrlWs = 'ws://' + location.host + '/ws';

const imagePlacehoder = "https://loremflickr.com/165/165/";

const ws = new WebSocket(searchUrlWs);
ws.binaryType = "blob";

ws.onmessage = function (event) {
    if (event.data instanceof Blob) {
        let reader = new FileReader();

        reader.onload = () => {
            console.log(reader.result);
            try {
                let city = JSON.parse(reader.result);
                getResultContainer(citiesWs).appendChild(
                    renderItem(city.name, city.temp)
                );
                getDurationElement(citiesWs).innerHTML = Math.floor(performance.now() - wsDuration) + " ms";

            } catch (e) {
                console.error("can't decode ws")
            }
        };

        reader.readAsText(event.data);
    } else {
        console.log("Result: " + event.data);
    }
}

function render(element, response) {
    if(Array.isArray(response.cities)) {
        response.cities.forEach(city => {
            getResultContainer(element).appendChild(
                renderItem(city.name, city.temp)
            );
        })
        getDurationElement(element).innerHTML = response.duration  + " ms";;
    } else {
        getResultContainer(element).appendChild(
            renderItem(response.name, response.temp)
        );
        // seems we're rendering ajax response we can update timer
        getDurationElement(element).innerHTML = Math.floor(performance.now() - ajaxDuration) + " ms";
    }
}

function renderItem(name, temp) {
    const item = document.importNode(cityTemplate.content, true);
    item.querySelector('.image').setAttribute(
        "src",
        imagePlacehoder + name
    );
    item.querySelector('.name').innerHTML = name;
    item.querySelector('.temp').innerHTML = temp + "°";

    return item
}
function requestAjax(url, cities, targetEl) {
    ajaxDuration = performance.now();

    let c = cities.split(",");
    c.forEach(city => {
        request(searchCityUrl, city, targetEl)
    })
}

function request(url, cities, targetEl) {
    fetch(url + cities)
        .then(
            function(response) {
                if (response.status !== 200) {
                    console.log('Status Code: ' + response.status);
                    return;
                }

                response.json().then(function(data) {
                    render(targetEl, data);
                });
            }
        )
        .catch(function(err) {
            console.log('Fetch Error:', err);
        });
}

function send(cities, targetEl) {
    wsDuration = performance.now();
    getResultContainer(targetEl).innerHTML = "";

    ws.send(cities);
}

document.getElementById("searchButton").addEventListener("click", function() {
    let cities = document.getElementById("search").value;

    getResultContainer(citiesAjaxContainer).innerHTML = '';
    getResultContainer(citiesContainer).innerHTML = '';
    getResultContainer(citiesAsyncContainer).innerHTML = '';
    getResultContainer(citiesWs).innerHTML = '';

    request(searchUrl, cities, citiesContainer)
    request(searchUrlAsync, cities, citiesAsyncContainer)
    send(cities, citiesWs);
    requestAjax(searchUrl, cities, citiesAjaxContainer)
});

function getResultContainer(el) {
    return el.querySelector('.results');
}

function getDurationElement(el) {
    return el.querySelector('.duration');
}