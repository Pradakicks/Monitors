const { MessageAttachment, MessageEmbed } = require('discord.js');
const delay = require('delay');
const fs = require('fs').promises;
const os = require('os');


const config = require('./config.json');
const mass = require('./discordHelpers/mass')
const {checkIfUserValidated ,getSkuBank, getValidatedIds, updateDiscordIdsDB, validateUser, updateSku, deleteSku} = require('./discordHelpers/db')
const {walmartScraper} = require('./scrapers/walmartScraper')
const {startGoMonitor} = require('./discordHelpers/go')
const { rp,firstServer, port } = require('./discordHelpers/helper')
//  var skuBank = []
let pushEndpoint = 'https://monitors-9ad2c-default-rtdb.firebaseio.com/monitor';
let discordIds =
  'https://monitors-9ad2c-default-rtdb.firebaseio.com/validatedUsers';


// const firstServer = `http://ec2-3-236-148-149.compute-1.amazonaws.com`
// let firstServer = `http://104.249.128.37`; // 12 Core z 24GB
// let firstServer = `http://104.249.128.207`;
var proxyList = [];
// let thirdSfirstServererver = `http://64.227.28.51`;
function findCommand(clients, triggerText, replyText) {
  clients.on('message', (message) => {
    if (message.content.toLowerCase() === triggerText.toLowerCase()) {
      message.author.send(replyText);
    }
  });
}

function checkBank(clients, triggerText, replyText) {
  clients.on('message', async (message) => {
    try {
      if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
        const { group, isValidated } = await checkIfUserValidated(message);
        if (isValidated) {
          let skuBank = await getSkuBank();
          let siteLimitObj = {};
          if (skuBank.length != 0) {
            let bankArr = [];
            let sites = Object.keys(skuBank);
            console.log(group);
            sites?.forEach((e) => {
              if (e != '-M_bpveXSTSxZkahEQkQ') {
                siteLimitObj[e] = 0;
                let currentSkus = Object.keys(skuBank[e]);
                currentSkus.forEach((sku) => {
                  skuBank[e][sku]?.companies?.forEach((company) => {
                    if (company.company == group) {
                      siteLimitObj[e]++;
                      bankArr.push(`${e}-${sku}-${group}`);
                    }
                  });
                });
              }
            });

            console.log(bankArr);
            console.log(siteLimitObj);
            Object.keys(siteLimitObj).forEach((e) => {
              if (siteLimitObj[e] > config.groups[group][`${e}LIMIT`])
                console.log('Site Limit Reached');
              console.log(siteLimitObj[e], config.groups[group][`${e}LIMIT`]);
            });
            await fs.appendFile(
              `monitorBank-${message.author.username}.txt`,
              JSON.stringify(bankArr, null, 2),
              (err) => {
                if (err)
                  message.content.send('Error While Creating Text Document');
                else console.log('File Sent');
              }
            );
            let attachment = new MessageAttachment(
              `monitorBank-${message.author.username}.txt`
            );
            message.channel.send(attachment);
            message.author.send('Attachment Successfully Fetched and Sent');
            message.author.send(`You have ${bankArr.length} products running`);
            await delay(2500);
            await fs.unlink(
              `monitorBank-${message.author.username}.txt`,
              (err) => {
                if (err) console.log('Error doing the unthinkable');
              }
            );
          } else {
            message.channel.send('Monitor Bank is empty');
          }
        } else {
          message.channel.send(`${message.author} is not a validated user`);
        }
      }
    } catch (error) {
      console.log(error);
      message.channel.send('Error checking Bank');
    }
  });
}

function massAdd(clients, triggerText, replyText) {
  try {
    clients.on('message', async (message) => {
      if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
        let string = message.content;
        const content = message.content;
        mass(string, content, message);
      }
    });
  } catch (error) {
    console.log(error);
  }
}




// Fire Base Sku Bank ----------------------------------------------
async function checkPresentSkus() {
  let skuBank = await rp.get({
    url: `${firstServer}:${port}/DB`,
  });
  skuBank = JSON.parse(skuBank?.body);
  let sites = Object.keys(skuBank);
  let numberOfSkus = 0;
  sites.forEach((site) => {
    let skus = Object.keys(skuBank[site]);
    if (site != '-M_iJkLwZh3hW5Pjys5Z') {
      skus.forEach((sku) => {
        numberOfSkus++;
      });
    }
  });
  console.log(`Number Of Items In Monitor : ${numberOfSkus}`);
  await delay(3000);
  sites.forEach(async (e) => {
    if (e != '-M_iJkLwZh3hW5Pjys5Z') {
      let site = e;
      let skus = Object.keys(skuBank[e]);
      for (let i = 0; i < skus.length; i++) {
        let s = skus[i];
        let currentSku = skuBank[site][s].original;
        let pricerange = '';
        if (currentSku?.includes('[') && site?.toUpperCase() !== 'TARGETNEW') {
          pricerange = currentSku?.split('[')[1]?.split(']')[0];
          currentSku = currentSku?.split('[')[0];
        }
        let currentBody = {
          site: site,
          sku: currentSku?.trim(),
          priceRangeMin: parseInt(pricerange?.split(',')[0]),
          priceRangeMax: parseInt(pricerange?.split(',')[1]),
          skuName:
            site == 'NEWEGG'
              ? await getSku(currentSku?.trim(), proxyList)
              : site == 'WALMARTNEW'
              ? 'prg=desktop&cat_id=0&facet=brand%3APanini%7C%7Cbrand%3ATopps%7C%7Cretailer%3AWalmart.com&grid=false&query=panini&soft_sort=false&sort=new'
              : '',
        };
        if (currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax) {
          console.log('No Max Price Range Detected');
          currentBody.priceRangeMax = 100000;
        }
        if (currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin) {
          console.log('No Min Price Range Detected');
          currentBody.priceRangeMin = 1;
        }
        console.log(currentBody, i, skus.length);
        startGoMonitor(currentBody, site);
        if (site.toUpperCase() == 'WALMART') await delay(15000);
        await delay(2000);
      }
    }
  });
}
  
if (os.platform() == 'win32' || os.platform() == 'darwin') {
  console.log('Development Environment');
  // firstServer = 'http://localhost';
  // firstServer = 'http://localhost';
} else {
  checkPresentSkus();
}
getProxies();

//-----------------------------------------------------------------

// Helper

async function getProxies() {
  try {
    // read contents of the file
    // const data = await fs.readFile('./GoMonitor/cloud.txt', 'utf-8');
    let fetchProxies = await rp.get(`${firstServer}:${port}/PROXY`);
    let parsed = JSON.parse(fetchProxies.body);
    console.log(parsed);
    parsed.proxies.forEach((line) => {
      const lineSplit = line.split(':');
      const item1 = {
        ip: lineSplit[0],
        port: lineSplit[1],
        userAuth: lineSplit[2],
        userPass: lineSplit[3],
      };
      proxyList.push(item1);
    });
    console.log(`Proxy list Length : ${proxyList.length}`);
    return proxyList;
  } catch (err) {
    console.error(err);
    // fs.appendFileSync('./errors.txt', err.toString() + '\n', (err) => {
    //   console.log(err);
    // });
  }
}
async function sendWebhook(body, webhook) {
  try {
    let sent = await rp.post({
      url: webhook,
      body: body,
    });
    return sent.statusCode;
  } catch (error) {
    console.log(error);
    return error?.statusCode;
  }
}

module.exports = {
  findCommand,
  deleteSku,
  checkBank,
  massAdd,
  validateUser,
  walmartScraper,
  startGoMonitor
};
