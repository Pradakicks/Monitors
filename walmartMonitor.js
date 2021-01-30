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

class walmartMonitor {
    constructor(sku) {
        this.sku = sku;
        this.delay = 10000;
        this.availability = '';
        this.proxyList = [];

    }

    async task () {
        try {
            console.log('Start')
            await this.getProxies()
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
       
            return new Promise( (resolve , reject) => {
                var myInterval = setInterval( async () => {
                await rp.get({
                    headers : {
                        "sec-fetch-dest": "document",
                        'sec-fetch-mode': 'navigate',
                        'sec-fetch-site': 'none',
                        'sec-fetch-user': '?1',
                        'service-worker-navigation-preload': 'true',
                        'upgrade-insecure-requests': '1',
                        'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 '
                    },
                    url : `https://www.walmart.com/ip/${this.sku}`
                }, ((error, response, body) => {
                    if (response.statusCode == 200 || response.statusCode == 201) {
                        console.log(response.headers)
                        console.log(body)
                        let parsedHeaders =  (response.headers.stockstatus)
                        this.availability = parsedHeaders
                        console.log(this.availability)
                    } else {
                        console.log('Task Failed')
                    }
                   
    
                }))
            }, this.delay) 
            
        })
    }
}

const monitoring = new walmartMonitor(127446742);


(async ()=>{
    await monitoring.task()
})()
