const delay = require('delay');
const puppeteer = require('puppeteer-extra')
const cheerio = require('cheerio')
const axios = require('axios').default;
const fs = require('fs');
const { gzip } = require('zlib');
const { response } = require('express');
const { resolve } = require('path');
const { json } = require('body-parser');
const fetch = require("node-fetch");
const Discord = require('discord.js');
let rp = require('request-promise').defaults({
    followAllRedirects: true,
    resolveWithFullResponse: true,
    gzip : true,
})
// https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM
// https://discordapp.com/api/webhooks/797249480410923018/NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW
const webhookClient1 = new Discord.WebhookClient('745279081247014942', '3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM');
const webhookClient = new Discord.WebhookClient('797249480410923018', 'NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW')




function getChromiumExecPath() {
    let platform = process.platform;
    if (platform === "win32") {
        return "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe";
    } else if (platform === "darwin") {
        return "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome";
    } else {
        return puppeteer.executablePath().replace("app.asar", "app.asar.unpacked");
    }
}
var proxyList = []
var g = 0
var i = 0
try {
    // read contents of the file
    const data = fs.readFileSync('proxies.txt', 'utf-8');
    // split the contents by new line
    const lines = data.split(/\r?\n/);

    // print all lines
    lines.forEach((line) => {
        let lineSplit = line.split(':')
        let item1 = {
            ip : lineSplit[0],
            port: lineSplit[1],
            userAuth: lineSplit[2],
            userPass: lineSplit[3]
        }
        proxyList.push(item1)
        
        // console.log(line);
        // console.log(item1);
        // console.log('\n\n\n\n\n');
    });
} catch (err) {
    console.error(err);
}
 // console.log(proxyList)
let tcin = 81588857
let father = {
    url : `https://www.target.com/p/-/A-${tcin}`,
    monitorDelay : 20000,
    startTask : 1000
}

let redUrl = `https://redsky.target.com/redsky_aggregations/v1/web/pdp_fulfillment_v1?key=ff457966e64d5e877fdbad070f276d18ecec4a01&tcin=${tcin}&store_id=2067&store_positions_store_id=2067&has_store_positions_store_id=true&scheduled_delivery_store_id=2067&pricing_store_id=2067&fulfillment_test_mode=grocery_opu_team_member_test`
// axios.get(redUrl).then((response) =>{
//     let $ = cheerio.load(response)
//     console.log(response)
// })
var ProductTCIN = '1'
var availability = '1'
var stockNum = '1'
var testNum = 1
var itemPicUrl = '1'
var productName = '1'
var getStatusByName = ''
var getIdByName =''
async function getItemsinDB () {
await rp.get({url : 'http://localhost:5003/api/items'}, ((error, response, body)=> {
    // console.log(JSON.parse(body))
    let parsed = JSON.parse(body)
    let string = (body)
    console.log(`Searching for TCIN ${tcin}`)
    // console.log(string)
    getStatusByName = string.split(`"name":"${tcin}","status":`)[1].split(',')[0]
    getIdByName = string.split(`","name":"${tcin}"`)[0].slice(-24)
    // console.log(getIdByName)
    // console.log(getStatusByName)
}))
}  

async function createItemsinDB () {
await fetch("http://localhost:5003/api/items", {
  "headers": {
    "accept": "*/*",
    "accept-language": "en-US,en;q=0.9",
    "cache-control": "no-cache",
    "content-type": "application/json",
  },
  "referrerPolicy": "no-referrer-when-downgrade",
  "body": `{\"name\":\"${tcin}\",\"status\":${testNum}}`,
  "method": "POST",
  "mode": "cors"
}).then((response)=>{
    console.log(response)
});
} 
//  createItemsinDB()
async function updateItemsinDB2 () {
let updateItem = await fetch(`http://localhost:5003/api/items/${getIdByName}`, {
  "headers": {
    "accept": "*/*",
    "accept-language": "en-US,en;q=0.9",
    "cache-control": "no-cache",
    "content-type": "application/json",
  },
  "referrerPolicy": "no-referrer-when-downgrade",
  "body": "{\r\n    \"status\" : 2\r\n}",
  "method": "PATCH",
  "mode": "cors",
  "credentials": "omit"
});
}
async function updateItemsinDB1 () {
    let updateItem = await fetch(`http://localhost:5003/api/items/${getIdByName}`, {
      "headers": {
        "accept": "*/*",
        "accept-language": "en-US,en;q=0.9",
        "cache-control": "no-cache",
        "content-type": "application/json",
      },
      "referrerPolicy": "no-referrer-when-downgrade",
      "body": "{\r\n    \"status\" : 1\r\n}",
      "method": "PATCH",
      "mode": "cors",
      "credentials": "omit"
    });
    }


async function task (j) {
    console.log(`Task ${i}`)
    console.log(g)
        let item = {
            url : father.url,
            monitorDelay : father.monitorDelay,
            proxyIp: "",
            proxyPort: "",
            proxyFull: proxyList[g].ip + ':' + proxyList[g].port, // ENTER THE ENTIRE PROXY HERE IP ADDRESS // If user has userpass DO NOT ENTER THE USER AND PASS HERE ONLY THE ADDRESS
            proxyUserAuth: proxyList[g].userAuth, // if user has userpass proxies enter the username here
            proxyPassAuth: proxyList[g].userPass, // if user has userpass proxies enter the password here
        }
        g++
        console.log(item)
        for (let h = 0; h < 100; h++){
        

            try {
                await rp.get({proxy : `http://${item.proxyUserAuth}:${item.proxyPassAuth}@${item.proxyFull}`, url: redUrl},
        function (error, response, body)  {
    // console.log(body)
    // console.log(response)
    // console.log(body)
    // console.log(error)
    
        ProductTCIN = (body.split('"tcin":"')[1].split('"')[0])
     availability = (body.split('"shipping_options":{"availability_status":"')[1].split('"')[0])
     stockNum = (body.split('"available_to_promise_quantity":')[1].split(',')[0])
     console.log(ProductTCIN)  
     console.log(availability)
     console.log(stockNum)
    
    
     if(availability.includes('OUT_OF_STOCK')){
        async function outOfStock() {
        await getItemsinDB ()
        console.log('Out of Stock')
        if (getStatusByName == 2){
        await updateItemsinDB1 () 
        }
         
         } outOfStock ()
        
     } 
     else if (availability.includes('IN_STOCK')){
         async function inStock () {

         
        await getItemsinDB ()
        console.log(getStatusByName)
        if (getStatusByName == 1){
        await updateItemsinDB2 ()
        let embed1 = new Discord.MessageEmbed()
         .setColor('#ff6666')
         .setTitle('Target Monitor')
         .setURL(`${father.url}`)
         .addField('Product Name', `${productName}`)
         .addField('Product Tcin', `${ProductTCIN}`, true)
         .addField('Product Availability', 'Product In Stock',true)
         .addField('Stock Number', `${stockNum}`, true)
         .setImage(`${itemPicUrl}`)
         .setTimestamp()
         .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
         webhookClient.send('Restock!', {
             username: 'Target',
             avatarURL: 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg',
             embeds: [embed1],
         })
         
         }
         console.log('In Stock')
        } inStock ()
     } 
    if (error) {
        console.log(error)
        return;
    }
})
await delay(father.monitorDelay)
            } catch (error) {
                console.log(error)
            }
        
        }
    
}
async function testWebhook () {
        webhookClient.send('Webhook test', {
        username: 'Target',
        avatarURL: 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg',
        embeds: [embed],
    })
    }
    var embed = new Discord.MessageEmbed()
    .setColor('#ff6666')
	.setTitle('Target Monitor')
	.setURL(`${father.url}`)
	.addField('Product Name', `${productName}`)
	.addField('Product Tcin', `${ProductTCIN}`, true)
	.addField('Product Availability', 'Product In Stock',true)
	.addField('Stock Number', `${stockNum}`, true)
	.setImage(`${itemPicUrl}`)
	.setTimestamp()
	.setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
async function startTask ()  {
for (let b = 0; b < 20; b++) {
    task(b)
    await delay(father.startTask)
}
} 
// startTask ()
async function gettingItemPic () {
    let getItemPic = await rp.get({url: father.url}, ((error, response, body)=>{
       // console.log(body)
        const $ = cheerio.load(response.body)
        let itemPic = $('img').first().attr('src')
         productName = $('h1').first().text()
        // console.log(itemPic)
         itemPicUrl = itemPic
    }))
} 





async function mainTask () {   
    try {
      gettingItemPic ()
    startTask()  
    } catch (error) {
     console.log(error)   
    }
    




    //console.log(itemPicUrl)
    

    
} 
  mainTask ()
