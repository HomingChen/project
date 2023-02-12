const loginModel = {
    async getAuthURL(authType){
        let result = await fetch("/api/getAuthURL?authType="+authType, {
            method: "GET"
        }).then((response)=>{
            return response.json();
        }).then((data)=>{
            console.log("model:", data);
            return data;
        }).catch((err)=>{
            console.log("fail sto excute function 'getAuthURL': ", err);
        });
        return result;
    },
};

const loginView = {
    renderElement(element){
        let tab = document.createElement(element.element);
        if(Object.hasOwn(element, "innerText")){
            tab.innerText = element.innerText;
        };
        if(Object.hasOwn(element, "attribute")){
            for(let i=0; i<element.attribute.length; i++){
                tab.setAttribute(element.attribute[i].name, element.attribute[i].value);
            };
        };
        document.querySelector(element.querrySelector).appendChild(tab);
    },
    popupWindow(){
        
    }
};

const loginControl = {
    redirectToAuthPage(btn){
        btn.addEventListener("click", async ()=>{
            let authURL = await loginModel.getAuthURL(btn.id);
            document.location.href = authURL;
            // console.log(authURL);
        });
    },
};

const googleBtn = document.querySelector("#google");
const notionBtn = document.querySelector("#notion");
const microsoftBtn = document.querySelector("#microsoft");
loginControl.redirectToAuthPage(googleBtn);
loginControl.redirectToAuthPage(notionBtn);
loginControl.redirectToAuthPage(microsoftBtn);

window.addEventListener("click", (event)=>{
    console.log("click: ", event.target.id);
});

// const loginWindow = {
//     "element": "div",
//     "querrySelector": "main",
//     "attribute": [
//         {
//             "name": "ID",
//             "value": "window"
//         },
//         {
//             "name:": "class",
//             "value": "class"
//         }
//     ]
// };
// const googleButton_div = {
//     "element": "div",
//     "querrySelector": "#google",
//     "attribute": [
//         {
//             "name": "id",
//             "value": "g_id_onload"
//         },
//         {
//             "name": "data-client_id",
//             "value": "1010807450371-q4l9h6dpse23g4crhb8kqdminrho7msp.apps.googleusercontent.com"
//         },
//         {
//             "name": "data-context",
//             "value": "signup"
//         },
//         {
//             "name": "data-ux_mode",
//             "value": "popup"
//         },
//         {
//             "name": "data-login_uri",
//             "value": "http://localhost:8080/api/getClientCode"
//         },
//         {
//             "name": "data-auto_prompt",
//             "value": "false"
//         }
//     ]
// };
// const googleButton_class = {
//     "element": "div",
//     "querrySelector": "#google",
//     "attribute": [
//         {
//             "name": "class",
//             "value": "g_id_signin"
//         },
//         {
//             "name": "data-type",
//             "value": "standard"
//         },
//         {
//             "name": "data-shape",
//             "value": "rectangular"
//         },
//         {
//             "name": "data-theme",
//             "value": "outline"
//         },
//         {
//             "name": "data-text",
//             "value": "signin_with"
//         },
//         {
//             "name": "data-size",
//             "value": "large"
//         },
//         {
//             "name": "data-logo_alignment",
//             "value": "left"
//         },
//         {
//             "name": "data-width",
//             "value": "300"
//         }
//     ]
// };
// loginView.renderElement(googleButton_div);
// loginView.renderElement(googleButton_class);