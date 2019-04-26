
let setup = new Setup()
let serverData = null


//Function to get params from URL
function getUrlVars() {
    var vars = {};
    var parts = window.location.href.replace(/[?&]([^=#]+)=([^&#]*)/gi, function(m,key,value) {
        vars[key] = value;
    });
    return vars;
}


ws.onopen = () => {

    console.log("Connected")
    var params = getUrlVars()
   
    //Send params from URL sent by server back to server to Init client player
    var X = parseFloat(params.PosX)
    var Y = parseFloat(params.PosY)
    ws.send(JSON.stringify({ID:params.ID,Position:{X:X,Y:Y}}))
   
    
}

ws.onmessage = (msg) => {
    serverData = JSON.parse(msg.data)

    // Control how many player position server sends and handle it
    setup.AddPlayer(serverData)
    setup.DeletePlayer(serverData)
    
    if (setup.player_list.length != 0) {

        for (var i = 0; i < setup.player_list.length; i++) {
            for (var j = 0; j < serverData.length; j++) {

                if (setup.player_list[i].ID == serverData[j].ID) {

                    // After recieving ID update position for particular ID
                    setup.player_list[i].UpdateData(serverData[j].Position.X, serverData[j].Position.Y)

                }
            }
        }


    }



}
//s

ws.onclose = () => {
    console.log("Disconnected")
}

// Control
document.onkeydown = function (event) {
    if (event.keyCode == 68)  //d
        ws.send(JSON.stringify({ Right: true }))
    if (event.keyCode === 83) //s
        ws.send(JSON.stringify({ Down: true }))
    if (event.keyCode === 65) //a
        ws.send(JSON.stringify({ Left: true }))
    if (event.keyCode === 87) //w
        ws.send(JSON.stringify({ Up: true }))
}
document.onkeyup = function (event) {
    if (event.keyCode == 68) { //d
        ws.send(JSON.stringify({ Right: false }))
    }
    if (event.keyCode === 83) //s
        ws.send(JSON.stringify({ Down: false }))
    if (event.keyCode === 65) //a
        ws.send(JSON.stringify({ Left: false }))
    if (event.keyCode === 87) //w
        ws.send(JSON.stringify({ Up: false }))
}

window.ws = ws

