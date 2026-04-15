function getUsers(){
let sender = new XMLHttpRequest;
sender.open("GET","http://localhost:8084/getUsers",false)
sender.send()
build = JSON.parse(sender.responseText)
console.log(build)

addPanel()

}





function addPanel(){
    let findContainer = document.querySelector(".mainContainer")


    build.forEach((element)=> {
        let checkUserPanel = document.getElementById(element.Id)
       
        if (!checkUserPanel){
            let userPanel = document.createElement("div")
            userPanel.classList.add("userPanel")
            userPanel.setAttribute('id', element.Id);

            userPanel.innerHTML = `<b>${element.First_name}</b><br>ID: ${element.Id}`

            let message = document.createElement("div")
            message.classList.add("messages")
            message.innerHTML = `${element.Text}`
            userPanel.appendChild(message)

            let inputBox = document.createElement("div")
            inputBox.classList.add("inputBox")

            let sendBtn = document.createElement("button")
            sendBtn.classList.add("sendBtn")
            sendBtn.innerHTML =`Отправить`
            

            let input = document.createElement("input")
            input.classList.add("input")

            sendBtn.addEventListener("click",()=>{


                let sender = new XMLHttpRequest;
                sender.open("POST","http://localhost:8084/getMessages",false)
                sender.setRequestHeader("Content-Type", "application/json")
                
                let Object = {
                    Id:element.Id,
                    Text:input.value
                }
                input.value = ""
                let jsonObject = JSON.stringify(Object)
                console.log(jsonObject)
                sender.send(jsonObject)


                location.reload()
            })

            input.addEventListener("keydown", (e) => {
                if (e.key === "Enter") {
                    sendBtn.click() 
                }
            })
            

            findContainer.appendChild(userPanel)
            findContainer.appendChild(inputBox)
            inputBox.appendChild(input)
            inputBox.appendChild(sendBtn)

        }else{            
            let message = document.createElement("div")
            message.classList.add("messages")
            message.innerHTML = `${element.Text}`
            checkUserPanel.appendChild(message)

        }
    });

}

document.addEventListener('DOMContentLoaded', getUsers)
