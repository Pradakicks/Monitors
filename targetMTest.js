
// let task = async () =>{ puppeteer.launch({
//     headless: true,
//     executablePath: getChromiumExecPath()
// }).then(async browser => {

//     const page = (await browser.pages())[0] 
//     await page.goto(item.url)
//     async function checkItem () {
        
//         let information =  await page.evaluate(()=>{
//         let itemPrice = document.querySelectorAll('div.h-padding-b-default')[0].firstElementChild.innerText
//         console.log(itemPrice)
//         // Finding the Ship it Button
//         let lengthOfClass = document.getElementsByClassName('h-padding-l-tiny').length
//         let itemElement1 = document.getElementsByClassName('h-padding-l-tiny')
//         console.log(itemElement1)
//         let time = Date.now()
//         if(lengthOfClass >= 8) {
//             if (itemElement1[7]){
//                 if(itemElement1[7].innerHTML.includes('Ship it')){
//                  var itemAvailability = 'Item is in Stock'
//                 } else if(!itemElement1[7].innerHTML.includes('Ship it')){
//                  var itemAvailability = 'Item is out of Stock'
//                 }
                 
//              } 
//         } 
//         else {
//            if(itemElement1[itemElement1.length - 1].innerHTML.includes('Ship it')){
//             var itemAvailability = 'Item is in Stock'
//            } else if(!itemElement1[itemElement1.length - 1].innerHTML.includes('Ship it')){
//             var itemAvailability = 'Item is out of Stock'
//            }
            
//         } 

//         // var itemAvailability = itemElement1.innerHTML
//         // itemElement1.forEach((item)=>{
//         //     console.log(item)
//         //     if (item.innerHTML.includes('Ship it')){
//         //         var itemAvailability = 'Item is in Stock'
//         //     } else if (!item.innerHTML.includes('Ship it')){
//         //         var itemAvailability = 'Item is out of Stock'
//         //     }
//         // })
//         // if (itemElement1.innerHTML.includes('Ship it')){
//         //     var itemAvailability = 'Item is in Stock'
//         // } else if (!itemElement1.innerHTML.includes('Ship it')){
//         //     var itemAvailability = 'Item is out of Stock'
//         // }
//         console.log(itemAvailability)
       

//         return {
//             itemPrice,
//             itemAvailability,
//             time       
//         }
//     });
//     console.log(information)
//      await page.reload()
//     } 
//     for (let i = 0; i < 100; i--){
//     console.log(i)
//     await checkItem ()
//     await page.waitForTimeout(item.monitorDelay)
//     console.log('Checking...')
//     }


    

// })} 
// let taskProxy = async (h) =>{ 
//     console.log(h)
//     console.log(g)
//     g++
//     let item = {
//         url : father.url,
//         monitorDelay : father.monitorDelay,
//         proxyIp: "",
//         proxyPort: "",
//         proxyFull: proxyList[g].ip + ':' + proxyList[g].port, // ENTER THE ENTIRE PROXY HERE IP ADDRESS // If user has userpass DO NOT ENTER THE USER AND PASS HERE ONLY THE ADDRESS
//         proxyUserAuth: proxyList[g].userAuth, // if user has userpass proxies enter the username here
//         proxyPassAuth: proxyList[g].userPass, // if user has userpass proxies enter the password here
//     }
    
//     console.log(item) 
//     puppeteer.launch({
//     headless: true,
//     executablePath: getChromiumExecPath(),
//     args: [`--proxy-server=${item.proxyFull}`]
// }).then(async browser => {
//     const page = (await browser.pages())[0] 
//     await page.setUserAgent('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36');
//     await page.authenticate({
//         username: item.proxyUserAuth,
//         password: item.proxyPassAuth
//     })
    
//     await page.goto(item.url)
//     async function checkItem () {
        
//         let information =  await page.evaluate(()=>{
//         let itemPrice = document.querySelectorAll('div.h-padding-b-default')[0].firstElementChild.innerText
//         console.log(itemPrice)
//         // Finding the Ship it Button
//         let lengthOfClass = document.getElementsByClassName('h-padding-l-tiny').length
//         let itemElement1 = document.getElementsByClassName('h-padding-l-tiny')
//         console.log(itemElement1)
//         let time = Date.now()
//         if(lengthOfClass >= 8) {
//             if (itemElement1[7]){
//                 if(itemElement1[7].innerHTML.includes('Ship it')){
//                  var itemAvailability = 'Item is in Stock'
//                 } else if(!itemElement1[7].innerHTML.includes('Ship it')){
//                  var itemAvailability = 'Item is out of Stock'
//                 }
                 
//              } 
//         } 
//         else {
//            if(itemElement1[itemElement1.length - 1].innerHTML.includes('Ship it')){
//             var itemAvailability = 'Item is in Stock'
//            } else if(!itemElement1[itemElement1.length - 1].innerHTML.includes('Ship it')){
//             var itemAvailability = 'Item is out of Stock'
//            }
            
//         } 

//         // var itemAvailability = itemElement1.innerHTML
//         // itemElement1.forEach((item)=>{
//         //     console.log(item)
//         //     if (item.innerHTML.includes('Ship it')){
//         //         var itemAvailability = 'Item is in Stock'
//         //     } else if (!item.innerHTML.includes('Ship it')){
//         //         var itemAvailability = 'Item is out of Stock'
//         //     }
//         // })
//         // if (itemElement1.innerHTML.includes('Ship it')){
//         //     var itemAvailability = 'Item is in Stock'
//         // } else if (!itemElement1.innerHTML.includes('Ship it')){
//         //     var itemAvailability = 'Item is out of Stock'
//         // }
//         console.log(itemAvailability)
       

//         return {
//             itemPrice,
//             itemAvailability,
//             time       
//         }
//     });
//     console.log(information)
//      await page.reload()
//     } 
//     for (let i = 0; i < 100; i--){
//     console.log(i)
//     await checkItem ()
//     console.log(`Task ${h} Checking...`)
//     await page.waitForTimeout(item.monitorDelay)
//     }


    

// })}
// async function fullFunction () {


// await delay(father.startTask)
// taskProxy (1)
// await delay(father.startTask)
// taskProxy (2)
// await delay(father.startTask)
// taskProxy (3)
// await delay(father.startTask)
// taskProxy (4)
// await delay(father.startTask)
// taskProxy (5)
// await delay(father.startTask)
// taskProxy (6)
// await delay(father.startTask)
// taskProxy (7)
// await delay(father.startTask)
// taskProxy (8)
// await delay(father.startTask)
// taskProxy (9)
// await delay(father.startTask)
// taskProxy (10)
// } fullFunction ()