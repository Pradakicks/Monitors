const delay = require('delay');
const cheerio = require('cheerio');
const axios = require('axios').default;
const fs = require('fs');
const { gzip } = require('zlib');
const { response } = require('express');
const { resolve } = require('path');
const { json } = require('body-parser');
const fetch = require('node-fetch');
const Discord = require('discord.js');
// require ('newrelic');
const rp = require('request-promise').defaults({
	followAllRedirects: true,
	resolveWithFullResponse: true,
	gzip : true,
    timeout: 50000
});
// https://discordapp.com/api/webhooks/816740348222767155/2APr1EdhzNO4hRWznexhMRlO0g7qOiCkI7HFtmuU7_r48PCWnGYmSTGJmRVX0LPCNN_t
 // const webhookClient1 = new Discord.WebhookClient('816740348222767155', '2APr1EdhzNO4hRWznexhMRlO0g7qOiCkI7HFtmuU7_r48PCWnGYmSTGJmRVX0LPCNN_t');
// const webhookClient = new Discord.WebhookClient('797249480410923018', 'NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW');
 const webhookClient1 = new Discord.WebhookClient('745279081247014942', '3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM');
// Test https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM

class bestBuyMonitor {
    constructor(sku) {
        this.trueSku = sku
        this.sku = sku.split(':')[0];
        this.skuName = sku.split(':')[1]
        this.delay = 850000; // this.delay = 390000
        this.startDelay = 1000; //  this.startDelay = 6000;
        this.availability = '';
        this.stockNumber = '';
        this.proxyList = [];
        this.isStock = false
        this.imageUrl = ''
        this.productName = ''

    }

    async task () {
        try {
            console.log('Start')
            console.log(this.sku)
            console.log(this.skuName)
            await this.getProxies()
         //   console.log(this.proxyList)
            await this.monitor()
        } catch (error) {
            fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                console.log(err)
            }))
        }
    }
    async getProxies() {
		try {
			// read contents of the file
			const data = fs.readFileSync('proxies.txt', 'utf-8');

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
            fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                console.log(err)
            }))
        }
    }
    async monitor () {
       try {
        console.log('Starting Monitoring')
        var testing = ''
        return new Promise( async ( resolve, reject) => {

                let i = 0
                var { skuBank } = require('../dms')
                let index = skuBank.findIndex(e => e.sku == this.trueSku)
                    while(!skuBank[index]?.stop){
                    skuBank[index].name = this.productName
                        if(i+1 == this.proxyList.length){
                            i = 0
                        }
                        let proxy = this.proxyList[i]
                        i++
                      //  console.log(`${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`);

                        (async () => {
                            try {
                                let fetchSite = await rp.get({
                                    headers : {
                                                "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
                                                "accept-language": "en-US,en;q=0.9",
                                                "cache-control": "max-age=0",
                                                "sec-fetch-dest": "document",
                                                "sec-fetch-mode": "navigate",
                                                "sec-fetch-site": "none",
                                                "sec-fetch-user": "?1",
                                                "upgrade-insecure-requests": "1",
                                                "user-agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML%2C like Gecko) Mobile/15E148",
                                    },
                                    url : `https://www.bestbuy.com/site/prada/${this.sku}.p?skuId=${this.sku}`,
                                    proxy : `http://${proxy.userAuth}:${proxy.userPass}@${proxy.ip}:${proxy.port}`
                                }, ((error, response, body) => {
                                    console.log(response?.statusCode)
                                }))   
                                let availability = fetchSite?.body?.split(',"availability":"')[1]?.split('.org/')[1]?.split('"')[0]
                                this.productName = fetchSite?.body?.split('"@context":"http://schema.org/","@type":"Product","name":"')[1]?.split('"')[0]
                                let price = fetchSite?.body?.split('"@context":"http://schema.org/","@type":"Product","name":"')[1]?.split('price":"')[1]?.split('"')[0]
                                let image = fetchSite?.body?.split('","image":"')[1]?.split('"')[0]
                                if(availability == "SoldOut"){
                                    this.availability = false
                                } else if (availability == "InStock"){
                                    this.availability = true
                                } else {
                                    console.log(availability)
                                }
                                
                                console.log(`Task ${i} : ${fetchSite.statusCode}`, this.availability, productName, price , this.isStock)
                                if(!this.isStock && this.availability) {
                                    // Send in stock webhook
                                    console.log(`Task ${i} : ${this.isStock} and ${this.availability}`)
                                    this.isStock = true
                                    let embed1 = new Discord.MessageEmbed()
                                    .setColor('#07bf6e')
                                    .setTitle('Best Buy Monitor')
                                    .setURL(`https://www.bestbuy.com/site/prada/${this.sku}.p?skuId=${this.sku}`)
                                    .addFields(
                                        { name : 'Product Name', value : `${this.productName}`},
                                        { name : 'Product Availability', value : `Product In Stock`, inline : true},
                                    //    { name : 'Stock Number', value : `${this.stockNumber}`, inline : true}, 
                                        { name : 'Current Price', value : price, inline : true} 
                                        )
                                    .addField("Links", `[ATC](https://api.bestbuy.com/click/tempo/${this.sku}/cart) | [Cart](https://www.bestbuy.com/cart)`)
                                    // .addField('Original Price', originalPrice, true)
                                    // .addField('Current Price', currentPrice, true)
                                    .setThumbnail(image)
                                    .setTimestamp()
                                    .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
                                    webhookClient1.send({
                                        username: 'Best Buy',
                                        avatarURL: 'https://cdn.discordapp.com/attachments/815507198394105867/816741454922776576/pfp.png',
                                        embeds: [embed1],
                                    })
                                } else if (!this.availability && this.isStock) {
                                    this.isStock = false
                                } else if (!this.availability && !this.isStock){
                                    this.isStock = false
                                    this.availability = false
                                    // Not Important
                                  //  console.log('False false')
                                } else if(this.availability && this.isStock){
                                    this.availability = true
                                    this.isStock = true
                                    // Not Important
                                    // console.log(true, true)
                                } else {
                                    fs.appendFileSync('what.txt', this.availability + this.isStock + '\n', (err =>{
                                        console.log(err)
                                    }))
                                }
            
                            
                                // console.log(originalPrice, currentPrice)
                               // console.log(parsedBod)
                
                
                            } catch (error) {
                                console.log(error)
                                // fs.appendFileSync('./errors.txt', error.toString() + '\n', (err =>{
                                //     console.log(err)
                                // }))
                                
                                if(error.message.includes('Unexpected token')){
                                      if(testing.includes('<!DOCTYPE html><html')){
                                        console.log("HTML")
                                      } else if (testing.includes('Are you a human?')) {
                                          console.log('Captcha')
                                      } else {
                                            console.log('No HTML')
                                            console.log(testing)
                                      }
                                   
                                } else if (error.message.includes('403')){
                                    // await delay(400000)
                                    console.log('403 Access Denied')
                                }
                            }
                        })()
                       
                        await delay(this.startDelay)
                    }
                    console.log('stopped!')
                    resolve('Stopped')
                    return
              
        })
       } catch (error) {
        fs.appendFileSync('bigError.txt', error.toString() + '\n', (err =>{
            console.log(err)
        }))
       }
  
    }

}

// const monitoring = new newEggMonitor(`19-113-569`);

// (async ()=>{
//     await monitoring.task()
// })()

module.exports = {
    bestBuyMonitor
}