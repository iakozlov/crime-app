function initMap() {
    initApp();
    var uluru = {lat: 37.75740784154647, lng: -122.44584196135858};

    let message = document.getElementById("message");
    message.style.visibility = 'hidden';

// The map, centered at Uluru
    const map = new google.maps.Map(document.getElementById("map"), {
        zoom: 4,
        center: uluru,
    });
// The marker, positioned at Uluru
    var marker = new google.maps.Marker({
        position: uluru,
        map: map,
    });
    var btn = document.getElementById("submit_button");
    btn.innerHTML = "TRAVEL HERE";
    map.addListener("click", (e) => {
        var myAlert = document.getElementById('myAlert');
        myAlert.style.visibility = "hidden";
        if (e.latLng.lat() > 37.819975492297004|| e.latLng.lat < 37.7078790224135 ||
            e.latLng.lng() <  -122.513642064265 || e.latLng.lng() > -122.365240723693) {

            var myAlert = document.getElementById('myAlert');
            myAlert.innerText = "Sorry, we don't have information about this area.";
            myAlert.style.visibility = "visible";
        } else {
            marker.setMap(null);
            marker = new google.maps.Marker({
                position: e.latLng,
                map: map,
            });
            uluru.lat = e.latLng.lat();
            uluru.lng = e.latLng.lng();
            map.setZoom(12);
            map.panTo(e.latLng);
            btn.style.visibility = 'visible';
        }
    });
}

function buildPieChart(first_crime, second_crime, third_crime, svg, name1, name2, name3) {
    var width = 450
    height = 450
    margin = 40

// The radius of the pieplot is half the width or half the height (smallest one). I subtract a bit of margin.
    let radius = Math.min(width, height) / 2 - margin;

// append the svg object to the div called 'my_dataviz'
    d3.select("svg").remove();
    var svg = d3.select("#my_dataviz")
        .append("svg")
        .attr("width", width)
        .attr("height", height)
        .append("g")
        .attr("transform", "translate(" + width / 2 + "," + height / 2 + ")");

// Create dummy data
    var data = {};
    data[name1] = first_crime;
    data[name2] = second_crime;
    data[name3] = third_crime;

// set the color scale
    var color = d3.scaleOrdinal()
        .domain(data)
        .range(d3.schemeSet2);

// Compute the position of each group on the pie:
    var pie = d3.pie()
        .value(function (d) {
            return d.value;
        })
    var data_ready = pie(d3.entries(data))
// Now I know that group A goes from 0 degrees to x degrees and so on.

// shape helper to build arcs:
    var arcGenerator = d3.arc()
        .innerRadius(0)
        .outerRadius(radius)

// Build the pie chart: Basically, each part of the pie is a path that we build using the arc function.
    svg
        .selectAll('mySlices')
        .data(data_ready)
        .enter()
        .append('path')
        .attr('d', arcGenerator)
        .attr('fill', function (d) {
            return (color(d.data.key))
        })
        .attr("stroke", "black")
        .style("stroke-width", "2px")
        .style("opacity", 0.7)

// Now add the annotation. Use the centroid method to get the best coordinates
    svg
        .selectAll('mySlices')
        .data(data_ready)
        .enter()
        .append('text')
        .text(function (d) {
            return d.data.key+ " "+(Math.round(d.data.value*100))+"%";
        })
        .attr("transform", function (d) {
            return "translate(" + arcGenerator.centroid(d) + ")";
        })
        .style("text-anchor", "middle")
        .attr("font-family", "sans-serif")
        .style("font-size", 14)
}
function makeDataReadable(data){
    data = JSON.stringify(data);
    if(data == "[]"){
        return "You don't have any trips yet";
    }
    data = data.replace(/['"]+/g, '');
    data = data.replace(/\\/g, '');
    var begin = data.lastIndexOf("[");
    var end = data.indexOf("]");
    var res_str = data.substr(begin+1, (end-begin)-1);
    var splited_data = res_str.split("; ");
    document.getElementById("list_group").innerHTML = "";
    for(var i = 0; i < splited_data.length; i++){
        var node = document.createElement("LI");
        node.className = "list-group-item";// Create a <li> node
        var textnode = document.createTextNode((i+1).toString()+")"+splited_data[i]);         // Create a text node
        node.appendChild(textnode);
        document.getElementById("list_group").appendChild(node);
    }
}

function initApp(){
    //init modal windows settings
    var modal = document.getElementsByClassName('modal');
    var span = document.getElementsByClassName('close');
    var infoMenu = document.getElementById("info_menu");
    var crimeMenu = document.getElementById("crime_menu");
    infoMenu.onclick = function (){
        modal[0].style.display = "block";
        crimeMenu.disabled = true;
    }
    crimeMenu.onclick = function (){
        modal[1].style.display = "block";
        infoMenu.disabled = true;
    }
    span[0].onclick = function (){
        modal[0].style.display = "none";
        crimeMenu.disabled = false;
    }
    span[1].onclick = function (){
        modal[1].style.display = "none";
        infoMenu.disabled = false;
    }


    const geocoder = new google.maps.Geocoder();
    var uluru = {lat: 37.75740784154647, lng: -122.44584196135858};
    const url = "127.0.0.1:5555";
    let inputForm = document.getElementById("inputForm");
    var svg;
    var submitButton = document.getElementById("submit_button");
    //POST Request to get crime analysis
    inputForm.addEventListener("submit", (e) => {
        e.preventDefault()
        var time = document.getElementById("date").value;
        geocoder.geocode({ location: uluru }, (results, status) => {
            if (status === "OK") {
                if (results[0]) {
                    ;
                    var geocodedAddressArr = (results[0].formatted_address).split(", ");
                    var geocodedAddress = geocodedAddressArr.join(" ");
                    if (time != "") {
                        var time_h = document.getElementById("appt").value;
                        var a;
                        var myAlert = document.getElementById('myAlert')
                        myAlert.style.visibility = "hidden";
                        var a = new Object();
                        a.lat = uluru.lat.toString();
                        a.lng = uluru.lng.toString();
                        a.date = time;
                        if(time_h != ""){
                            a.time = time_h+":00";
                        }
                        else {
                            a.time = "12:00:00";
                        }
                        a.username = document.getElementById("usernameForTrips").value;
                        a.address = geocodedAddress;
                        submitButton.disabled = true;
                        fetch(url, {
                            method: "POST",
                            body: JSON.stringify(a),
                            headers:{
                                'Content-Type': 'application/json'
                            }
                        }).then(
                            response => response.text()
                        ).then(
                            (data) => {
                                data = JSON.stringify(data);
                                data = data.replace(/\\/g, '');
                                data = data.replace('r','');
                                data = data.replace(/['"]+/g, '');
                                buildPieChart(parseFloat(data.split(";")[0].split(":")[1]),
                                    parseFloat(data.split(";")[1].split(":")[1]),
                                    parseFloat(data.split(";")[2].split(":")[1]), svg,
                                    data.split(";")[0].split(":")[0],
                                    data.split(";")[1].split(":")[0],
                                    data.split(";")[2].split(":")[0]);
                                submitButton.disabled = false;
                            }

                        ).catch(
                            error => console.error(error)
                        )

                    } else {
                        var myAlert = document.getElementById('myAlert');
                        myAlert.innerText = "You should enter date of travelling.";
                        myAlert.style.visibility = "visible";
                    }
                } else {
                    alert("No results found");
                }
            } else {
                alert("Geocoder failed due to: " + status);
            }
        });

    });

    let historyForm = document.getElementById("getTripsForm");
    let tripsMessage = document.getElementById("tripsMessage");
    //GET Request to get history of requests
    historyForm.addEventListener("submit", (e) => {
        e.preventDefault()
        var name = document.getElementById("username").value;
        var myAlert = document.getElementById('myAlert');
        myAlert.style.visibility = "hidden";
        if(name.trim() != "") {
            tripsMessage.setAttribute('value', name);
            const formData = new FormData(historyForm)
            fetch("http://127.0.0.1:5555/reg" + "?username=" + name)
                .then(
                    response => response.text())
                .then(
                    (data) => {
                        makeDataReadable(data);
                    }
                ).catch(
                error => console.error(error)
            )
        } else{
            var myAlert = document.getElementById('myAlert');
            myAlert.innerText = "Name shouldn't be empty.";
            myAlert.style.visibility = "visible";
        }

    });
}

