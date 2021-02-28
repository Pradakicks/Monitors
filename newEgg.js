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

class newEggMonitor {
    constructor(sku) {
        this.sku = sku;
        this.delay = 10000;
        this.availability = '';
        this.stockNumber = '';
        this.proxyList = [];
        this.isStock = false

    }

    async task () {
        try {
            console.log('Start')
         //   await this.getProxies()
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
				// console.log(item1);
				// console.log('\n\n\n\n\n');
			});
		} catch (err) {
			console.error(err);
		}
    }

    async monitor () {
        console.log('Starting Monitoring')
        return new Promise( async ( resolve, reject) => {
            let montiorInterval = setInterval(async () => {
                
                try {
                    let fetchSite = await rp.get({
                        url : `https://www.newegg.com/product/api/ProductRealtime?ItemNumber=${this.sku}`
                    })   
                    console.log(fetchSite.statusCode)
                    let parsedBod = JSON.parse(fetchSite.body)
                    let productName = parsedBod.MainItem.Description.Title
                    let originalPrice = parsedBod.MainItem.OriginalUnitPrice
                    let currentPrice = parsedBod.MainItem.FinalPrice
                    this.availability = parsedBod.MainItem.Instock
                    this.stockNumber = parsedBod.MainItem.Stock 
                    if(!this.isStock && this.availability) {
                        // Send in stock webhook
                        this.isStock = true
                        let embed1 = new Discord.MessageEmbed()
                        .setColor('#00FF00')
                        .setTitle('New Egg Monitor')
                        .setURL(`https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg`)
                        .addField('Product Name', `${productName}`)
                        .addField('Product Availability', 'Product In Stock',true)
                        .addField('Stock Number', `${this.stockNumber}`, true)
                        .addField('Original Price', originalPrice)
                        .addField('Current Price', currentPrice)
                        .setImage(`https://c1.neweggimages.com/ProductImage/${this.sku}-V10.jpg`)
                        .setTimestamp()
                        .setFooter('Prada#4873', 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg');
                        webhookClient1.send('Restock!', {
                            username: 'New Egg',
                            avatarURL: 'https://cdn.discordapp.com/attachments/772173046235529256/795132477659152444/pradakicks.jpg',
                            embeds: [embed1],
                        })
                    } else if (!this.availability && this.isStock) {
                        this.isStock = false
                    }

                    console.log(this.availability, this.stockNumber, productName)
                    console.log(originalPrice, currentPrice)
                   // console.log(parsedBod)
    
    
                } catch (error) {
                    console.log(error)
                }
            }, 1000)

        })
    }

}

const monitoring = new newEggMonitor(`19-113-569`);

(async ()=>{
    await monitoring.task()
})()

module.exports = newEggMonitor