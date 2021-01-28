const delay = require('delay');
const puppeteer = require('puppeteer-extra');
const cheerio = require('cheerio');
const axios = require('axios').default;
const fs = require('fs');
const { gzip } = require('zlib');
const { response } = require('express');
const { resolve } = require('path');
const { json } = require('body-parser');
const fetch = require('node-fetch');
const Discord = require('discord.js');
require ('newrelic');

const rp = require('request-promise').defaults({
	followAllRedirects: true,
	resolveWithFullResponse: true,
	gzip : true,
});
// https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM
// https://discordapp.com/api/webhooks/797249480410923018/NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW
const webhookClient1 = new Discord.WebhookClient('745279081247014942', '3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM');
const webhookClient = new Discord.WebhookClient('797249480410923018', 'NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW');


async function TargetMonitor (item) {
class targetMonitor {
	constructor(tcin) {
        this.tcin = tcin;
        this.proxyList = [];
        this.fatherUrl = `https://www.target.com/p/-/A-${this.tcin}`;
        this.monitorDelay =  22500;
        this.startTaskDelay = 500;
        this.redURL = `https://redsky.target.com/redsky_aggregations/v1/web/pdp_fulfillment_v1?key=ff457966e64d5e877fdbad070f276d18ecec4a01&tcin=${this.tcin}&store_id=2067&store_positions_store_id=2067&has_store_positions_store_id=true&scheduled_delivery_store_id=2067&pricing_store_id=2067&fulfillment_test_mode=grocery_opu_team_member_test`
        this.ProductTCIN = '1';
        this.availability = '1';
        this.stockNum = '1';
        this.testNum = 1;
        this.itemPicUrl = '1';
        this.productName = '1';
        this.getStatusByName = '';
        this.getIdByName ='';
        this.testStatusCode = '';
        this.g = 0;
        this.testTCINStatusCode = ''
    }
    async task() {
        try {
            await this.getProxies()
          //  console.log(this.proxyList)
            await this.testProduct()
            this.createItemsInDB()
            this.runMonitor()
        } catch (error) {
            console.log(JSON.stringify(error.message))
            let message = error.message.split('message')[1].split(',')[0]
            console.log(message)
            rp.patch({
                url : `https://montiors-default-rtdb.firebaseio.com/sites/target/${this.tcin}.json`,
                body : JSON.stringify({
                "Status" : "Not Active",
                "reason" : `${message}`
                })
            })
            throw new Error(error)
        }
    


    }
	async getProxies() {
		try {
			// read contents of the file
			const data = fs.readFileSync('proxies.txt', 'utf-8');
			// split the contents by new line
			const lines = data.split(/\r?\n/);

			// print all lines
			lines.forEach((line) => {
				const lineSplit = line.split(':');
				const item1 = {
					ip : lineSplit[0],
					port: lineSplit[1],
					userAuth: lineSplit[2],
					userPass: lineSplit[3],
				};
				this.proxyList.push(item1);

				// console.log(line);
				// console.log(item1);
				// console.log('\n\n\n\n\n');
			});
		} catch (err) {
			console.error(err);
		}
    }
    async getItemPic () {
        try {
            let getItemPic = await rp.get({url: this.fatherUrl}, ((error, response, body)=>{
                // console.log(body)
                 const $ = cheerio.load(response.body)
                 let itemPic = $('img').first().attr('src')
                  this.productName = $('h1').first().text()
                 // console.log(itemPic)
                  this.itemPicUrl = itemPic
             }))
        } catch (error) {
            throw new Error(error.message)
        }
       
    }
    async getItemsInDB () {
        try {
            await rp.get({url : 'https://quiet-dusk-97663.herokuapp.com/api/items'}, ((error, response, body) => {
            // console.log(JSON.parse(body))
            if(body){
                let parsed = JSON.parse(body)
                let string = (body)
                // console.log(parsed)
                console.log(`Searching for TCIN ${this.tcin}`)
                // console.log(string)
                
                if(string.split(`"name":"${this.tcin}","status":`)[1]){
                    this.getStatusByName = string.split(`"name":"${this.tcin}","status":`)[1].split(',')[0]
                }
                if (string.split(`","name":"${this.tcin}"`)[0]){
                    this.getIdByName = string.split(`","name":"${this.tcin}"`)[0].slice(-24)
                }
                
                // console.log(getIdByName)
                // console.log(getStatusByName)
            } else if(!body){
                console.log(response)
            }
           
        }))
            } catch (error) {
                throw new Error(error.message)
            }
    }
    async createItemsInDB() {
        try {
            await this.getItemsInDB()
            console.log(this.getIdByName)
            console.log(this.getStatusByName)
            
                if (this.getStatusByName == undefined) {
                console.log('Creating Item')
                let createItem = await fetch("https://quiet-dusk-97663.herokuapp.com/api/items", {
                    "headers": {
                      "accept": "*/*",
                      "accept-language": "en-US,en;q=0.9",
                      "cache-control": "no-cache",
                      "content-type": "application/json",
                    },
                    "referrerPolicy": "no-referrer-when-downgrade",
                    "body": `{\"name\":\"${this.tcin}\",\"status\":${this.testNum}}`,
                    "method": "POST",
                    "mode": "cors"
                  }).then((response)=>{
                      console.log(response)
                  });
            } else if (!this.getStatusByName == undefined) {
                console.log('TCIN already created in DB')
            }
        } catch (error) {
        throw new Error(error.message)   
        }
    }
    async updateItemsInDB2 () {
        try {
            let updateItem = await fetch(`https://quiet-dusk-97663.herokuapp.com/api/items/${this.getIdByName}`, {
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
        } catch (error) {
            throw new Error(error.message)
        }
    }
    async updateItemsinDB1 () {
        try {
            let updateItem = await fetch(`https://quiet-dusk-97663.herokuapp.com/api/items/${this.getIdByName}`, {
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
        } catch (error) {
            throw new Error(error.message)
        }
    }
    async taskMonitor(j) {
        try {
            console.log(`Task ${j}`)
                //   console.log(`Task ${g}`)
                  // console.log(`Task ${b}`)
                   console.log(this.g)
               
                       let item = {
                           url : this.fatherUrl,
                           monitorDelay : this.monitorDelay,
                           proxyIp: "",
                           proxyPort: "",
                           proxyFull: this.proxyList[this.g].ip + ':' + this.proxyList[this.g].port, // ENTER THE ENTIRE PROXY HERE IP ADDRESS // If user has userpass DO NOT ENTER THE USER AND PASS HERE ONLY THE ADDRESS
                           proxyUserAuth: this.proxyList[this.g].userAuth, // if user has userpass proxies enter the username here
                           proxyPassAuth: this.proxyList[this.g].userPass, // if user has userpass proxies enter the password here
                       }
                       this.g++
                       console.log(item)
                       for (let h = 0; h < 100; h++){
                           
                        await rp.get({proxy : `http://${item.proxyUserAuth}:${item.proxyPassAuth}@${item.proxyFull}`, url: this.redURL},
                               (error, response, body)  => {
                                if(response) console.log(response.statusCode)

                   if(body){
                       
                       if(!body.includes("No product found with tcin")){
                           this.ProductTCIN = (body.split('"tcin":"')[1].split('"')[0])
                           
                           if(body.split('"shipping_options":{"availability_status":"')[1]) {
                            this.availability = (body.split('"shipping_options":{"availability_status":"')[1].split('"')[0])
                           }
                           if(body.split('"available_to_promise_quantity":')[1]){
                            this.stockNum = (body.split('"available_to_promise_quantity":')[1].split(',')[0])
                           }
                           
                           console.log(this.ProductTCIN)  
                           console.log(this.availability)
                           console.log(this.stockNum)
                           if(this.availability){
                            if(this.availability.includes('OUT_OF_STOCK') || this.availability.includes('PRE_ORDER_UNSELLABLE')){
                                this.outOfStock()       
                             } else if (this.availability.includes('IN_STOCK') || this.availability.includes("PRE_ORDER_SELLABLE")){
                              this.inStock()
                             } else {
                                 console.log('Unknown Error')
                                 console.log(body)
                             }
                           }
                           
                           
                       } else {
                           console.log(`No product found with tcin - ${this.tcin}`)
                           console.log('Check TCIN')

                       }
                   }
                    
               })
                   await delay(this.monitorDelay)
       
                       
                       }
        } catch (error) {
            throw new Error(error.message)
        }
    }
    async inStock() {
        try {
            await this.getItemsInDB()
            console.log(this.getStatusByName)
            if (this.getStatusByName == 1){
            await this.updateItemsInDB2()
            let embed1 = new Discord.MessageEmbed()
             .setColor('#ff6666')
             .setTitle('Target Monitor')
             .setURL(`${this.fatherUrl}`)
             .addField('Product Name', `${this.productName}`)
             .addField('Product Tcin', `${this.ProductTCIN}`, true)
             .addField('Product Availability', 'Product In Stock',true)
             .addField('Stock Number', `${this.stockNum}`, true)
             .setImage(`${this.itemPicUrl}`)
             .setTimestamp()
             .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
             webhookClient.send('Restock!', {
                 username: 'Target',
                 avatarURL: 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg',
                 embeds: [embed1],
             })
             
             }
             console.log('In Stock')
        } catch (error) {
            throw new Error(error.message)
        }
       
    }
    async outOfStock() {
        try {
            await this.getItemsInDB()
            console.log('Out of Stock')
            if (this.getStatusByName == 2){
            await this.updateItemsinDB1 () 
            } else {
                return;
            }
                           
        } catch (error) {
            throw new Error(error.message)
        }
    }
    async testWebhook () {
        webhookClient.send('Webhook test', {
            username: 'Target',
            avatarURL: 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg',
            embeds: [embed],
        })
    }
    async testProduct () {
        try {
            let item = {
                url : this.fatherUrl,
                monitorDelay : this.monitorDelay,
                proxyIp: "",
                proxyPort: "",
                proxyFull: this.proxyList[5].ip + ':' + this.proxyList[5].port, // ENTER THE ENTIRE PROXY HERE IP ADDRESS // If user has userpass DO NOT ENTER THE USER AND PASS HERE ONLY THE ADDRESS
                proxyUserAuth: this.proxyList[5].userAuth, // if user has userpass proxies enter the username here
                proxyPassAuth: this.proxyList[5].userPass, // if user has userpass proxies enter the password here
            }
            let testTCIN = await rp.get({proxy : `http://${item.proxyUserAuth}:${item.proxyPassAuth}@${item.proxyFull}`, url: this.redURL}, ((error, response, body) => {
                console.log(response.statusCode)
                this.testTCINStatusCode = response.statusCode
                if (response.statusCode == 200) { 
                    rp.patch({
                    url : `https://montiors-default-rtdb.firebaseio.com/sites/target/${this.tcin}.json`,
                    body : `{
                        "Status" : "Active"
                        }`
                })
                } 
            }))
        } catch (error) {
            throw new Error(error.message)
        }
     
    }
    async runMonitor() {
    try {
    if (this.testTCINStatusCode == 200 || this.testTCINStatusCode == 201){
        this.getItemPic()
        for (let b = 0; b < this.proxyList.length; b++) {
            try {
                this.taskMonitor(b)
                await delay(this.startTaskDelay) 
            } catch (error) {
            rp.patch({
            url : `https://montiors-default-rtdb.firebaseio.com/sites/target/${this.tcin}.json`,
            body : `{
            "Status" : "Not Active"
            }`
        })
                
                console.log(error)   
            }
           
        }
    } else {
        console.log(`TCIN Not Found`)
        throw new Error(`TCIN ${this.tcin} is not FOUND`)
    }
    } catch (error) {
    throw new Error(error.message)
    }
    }


    

}
const monitoring = new targetMonitor(item);
(async () => {
    await monitoring.task()
})()
}


module.exports = TargetMonitor;