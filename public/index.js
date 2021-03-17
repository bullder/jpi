const cityTemplate = document.getElementById('city-template');
const cityCol = document.getElementById('city-col-template');

let citiesContainer = document.getElementById('cities');
let citiesAsyncContainer = document.getElementById('citiesAsync');
let citiesWs = document.getElementById('citiesWs');

let wsDuration = performance.now();

const searchUrl = './api/cities/';
const searchUrlAsync = './api/citiesAsync/';
const searchUrlWs = 'ws://' + location.host + '/ws';

const imagePlacehoder = "https://loremflickr.com/420/225/";
const ws = new WebSocket(searchUrlWs);
ws.binaryType = "blob";

ws.onmessage = function (event) {
    if (event.data instanceof Blob) {
        let reader = new FileReader();

        reader.onload = () => {
            console.log(reader.result);
            try {
                let city = JSON.parse(reader.result);
                citiesWs.appendChild(
                    renderItem(city.name, city.temp)
                );
                citiesWs.querySelector('.duration').innerHTML =  performance.now() - wsDuration + " ms"

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
    const col = document.importNode(cityCol.content, true);
    col.querySelector('.duration').innerHTML = response.duration;

    response.cities.forEach(city => {
        col.appendChild(
            renderItem(city.name, city.temp)
        );
    })

    element.appendChild(col)
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

function request(url, cities, targetEl) {
    fetch(url + cities)
        .then(
            function(response) {
                if (response.status !== 200) {
                    console.log('Status Code: ' +
                        response.status);
                    return;
                }

                response.json().then(function(data) {
                    console.log(data);
                    render(
                        targetEl,
                        data
                    );

                });
            }
        )
        .catch(function(err) {
            console.log('Fetch Error:', err);
        });
}

function send(cities, targetEl) {
    wsDuration = performance.now();

    const col = document.importNode(cityCol.content, true);
    col.querySelector('.duration').innerHTML = wsDuration.toString() + " ms";
    targetEl.appendChild(col)

    ws.send(cities);
}

document.getElementById("searchButton").addEventListener("click", function() {
    let cities = document.getElementById("search").value;

    citiesContainer.innerHTML = '';
    citiesAsyncContainer.innerHTML = '';
    citiesWs.innerHTML = '';

    request(searchUrl, cities, citiesContainer)
    request(searchUrlAsync, cities, citiesAsyncContainer)
    send(cities, citiesWs);
});
