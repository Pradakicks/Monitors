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
});
const webhookClient1 = new Discord.WebhookClient('745279081247014942', '3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM');
const webhookClient = new Discord.WebhookClient('797249480410923018', 'NPL3ktXS78z5EHo_cpYyrtFl_2iB0ARgz9IW5kwAZA-UkiseiinnBmUPJZlGgxw8TZiW');

const webhook = require("webhook-discord")
 // https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM
const Hook = new webhook.Webhook("https://discordapp.com/api/webhooks/745279081247014942/3TuT8vs6BUXr9HAK1uRKaB4t3Ap0LnoLfPJTgT1uhNzQvqR1GsUXW-d4_dxCrgOCdkBM")

class targetMonitor {
    constructor(sku) {
        this.sku = sku;
        this.delay = 10000;
        this.availability = '';
        this.stockNumber = '';
        this.proxyList = [];
        this.isStock = false;
        this.productName = '';
        this.itemPicUrl = ''

    }

    async task () {
        try {
            console.log('Start')
            await this.getProxies()
         //   console.log(this.proxyList)
            await this.monitor()
        } catch (error) {
            
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
			
				// console.log('\n\n\n\n\n');
			});
		} catch (err) {
			console.error(err);
		}
    }
    async monitor () {
        console.log('Starting Monitoring')
        var testing = ''
        return new Promise( async ( resolve, reject) => {
            let getItemPic = await rp.get({url: `https://www.target.com/p/-/A-${this.sku}`}, ((error, response, body)=>{
                    // console.log(body)
                     const $ = cheerio.load(response.body)
                     let itemPic = $('img').first().attr('src')
                      this.productName = $('h1').first().text()
                      this.itemPicUrl = itemPic
                 }))
            for(let i = 0 ; i < this.proxyList.length; i++){
                let monitorInterval = setInterval(async () => {
                    try {
                        let fetchSite = await rp.get({
                            url : `https://redsky.target.com/redsky_aggregations/v1/web/pdp_fulfillment_v1?key=ff457966e64d5e877fdbad070f276d18ecec4a01&tcin=${this.sku}&store_id=2067&store_positions_store_id=2067&has_store_positions_store_id=true&scheduled_delivery_store_id=2067&pricing_store_id=2067&fulfillment_test_mode=grocery_opu_team_member_test`
                        })   
                        console.log(fetchSite.statusCode)
                        testing = fetchSite.body
                        let parsedBod = JSON.parse(fetchSite.body)
                     //   console.log(fetchSite.body)
                   //     let productName = parsedBod.
                        let originalPrice = 'N/A'
                        let currentPrice = 'N/A'
                        this.availability = parsedBod.data.product.fulfillment.shipping_options.availability_status
                        this.stockNumber = parsedBod.data.product.fulfillment.shipping_options.available_to_promise_quantity 
                        let tcin = parsedBod.data.product.tcin
                        if(this.availability == "PRE_ORDER_SELLABLE" || this.availability == "IN_STOCK") {
                            this.availability = true
                        } else if(this.availability == "PRE_ORDER_UNSELLABLE" || this.availability == "UNAVAILABLE" || this.availability == undefined){
                            this.availability = false
                        }
    
    
                        if(!this.isStock && this.availability) {
                            // Send in stock webhook
                            this.isStock = true
                            let embed1 = new Discord.MessageEmbed()
                          //  .setColor('#00FF00')
                            .setTitle('Target Monitor')
                            .setThumbnail(`${this.itemPicUrl}`)
                            .setURL(`https://www.target.com/prada/-/A-${this.sku}`)
                            .addField('Product Name', `${this.productName}`)
                            .addField('Product Availability', 'Product In Stock',true)
                            .addField('Stock Number', `${this.stockNumber}`, true)
                            .addField('Original Price', originalPrice)
                            .addField('Current Price', currentPrice , true)
                          //  .setImage(`${this.itemPicUrl}`)
                            .setTimestamp()
                            .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
                            webhookClient1.send('Restock!', {
                                username: 'Target',
                                avatarURL: 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg',
                                embeds: [embed1],
                            })
                        } else if (!this.availability && this.isStock) {
                            this.isStock = false
                        }
    
                        console.log(`Task ${i} : ${this.availability}, ${this.productName}}`)
                       // console.log(originalPrice, currentPrice)
                       // console.log(parsedBod)
        
        
                    } catch (error) {
                        console.log(error)
                        if(error.message.includes('Unexpected token')){
                            console.log(testing)
                            clearInterval(monitorInterval)
                            resolve('g')
                        }
                    }
                }, this.delay)
                await delay(100)
            }
         

        })
    }

}

 // const monitoring = new targetMonitor(81409304);

// (async ()=>{
//     await monitoring.task()
// }) ()

module.exports = {
    targetMonitor
}