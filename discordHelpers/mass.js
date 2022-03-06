const { getValidatedIds, getSkuBank, updateSku } = require('./db');
const { startGoMonitor } = require('./go');
const config = require('../config.json');
const { rp,firstServer, port } = require('./helper')
const delay = require('delay');


async function mass(string, content, message, groupName) {
  let spaceLength = content.split(' ');
  console.log(`Space Length : ${spaceLength} : ${spaceLength.length}`);
  spaceLength.splice(0, 1);
  spaceLength.splice(0, 1);
  console.log(`Space : ${spaceLength} : ${spaceLength.length}`);
  const site = content?.split(' ')[1];
  console.log(site);

  let validatedIds = await getValidatedIds();
  let parsed = JSON.parse(validatedIds);
  let isValidated = false;
  let group;

  if (message) {
    parsed.forEach((e) => {
      let id = e?.split('-')[0];
      console.log(id, message.author.id);
      if (id == message.author.id) {
        isValidated = true;
        group = e?.split('-')[1];
      }
    });
  } else {
    isValidated = true;
    group = groupName;
  }
  if (isValidated) {
    //	const SKU = content.split(' ')[2];
    //	console.log(site)
    let g = spaceLength;
    //	console.log(g)
    for (let i = 0; i < g.length; i++) {
      if (!g[i].toUpperCase().includes('!MASSADD')) {
        let isContinue = true;
        let SKU;
        let pricerange = '';
        SKU = g[i].trim();
        let original = g[i];
        if (g[i].includes('[')) {
          pricerange = g[i].split('[')[1].split(']')[0];
          SKU = g[i].split(' ')[0];
        }
        if (site.toLowerCase().includes('shopifylink')) {
          SKU = SKU.split('https://')[1];
        }
        console.log(`Site : ${site.toUpperCase()} : SKU : ${SKU} : Pos : ${i}`);
        let skuBank = await getSkuBank();
        let caseSite = site.toUpperCase();
        let currentCompany = config.groups[group];

        // Check Site Limit
        let siteLimitObj = {};
        let allSites = Object.keys(skuBank);
        allSites?.forEach((e) => {
          if (e != '-M_bpveXSTSxZkahEQkQ') {
            siteLimitObj[e] = 0;
            let currentSkus = Object.keys(skuBank[e]);
            currentSkus.forEach((sku) => {
              skuBank[e][sku]?.companies?.forEach((company) => {
                if (company.company == group) {
                  siteLimitObj[e]++;
                }
              });
            });
          }
        });
        console.log(siteLimitObj, group);
        let shouldContinue = true;
        Object.keys(siteLimitObj).forEach((e) => {
          if (siteLimitObj[e] > config.groups[group][`${e}LIMIT`]) {
            message.channel.send(`${e} Limit Reached`);
            console.log(siteLimitObj[e], config.groups[group][`${e}LIMIT`]);
            shouldContinue = false;
          }
        });

        if (!shouldContinue) continue;
        // ------------------------- ---------------- -------- ------------------------ --------------------

        if (skuBank[caseSite]) {
          if (skuBank[caseSite][SKU]) {
            let skuWebhookArray = skuBank[caseSite][SKU]?.companies;
            let isPresent = false;
            skuWebhookArray.forEach((e) => {
              console.log(e.webhook, currentCompany[caseSite]);
              if (e.webhook == currentCompany[caseSite]) isPresent = true;
            });
            if (isPresent) {
              message.channel.send(`${SKU} is already present in monitor`);
            } else {
              message.channel.send(`${SKU} is being added to monitor`);
              let arr = [];
              skuWebhookArray.forEach((e) => {
                arr.push(e);
              });
              arr.push({
                company: group,
                webhook: currentCompany[caseSite],
                color: currentCompany?.companyColor,
                companyImage: currentCompany?.companyImage,
              });
              console.log(arr);
              skuBank[caseSite][SKU].companies = arr;
              await updateSku(site, SKU, skuBank[caseSite][SKU]);
            }
            console.log('Duplicate Found', isPresent);
            isContinue = false;
          }
        }
        if (isContinue) {
          let currentObj = {
            sku: SKU,
            site: site.toUpperCase(),
            stop: false,
            name: '',
            original: original,
            companies: [
              {
                company: group,
                webhook: currentCompany[caseSite],
                color: currentCompany?.companyColor,
                companyImage: currentCompany?.companyImage,
              },
            ],
          };

          let currentBody = {
            site: site.toUpperCase(),
            sku: SKU,
            priceRangeMin: parseInt(pricerange.split(',')[0]),
            priceRangeMax: parseInt(pricerange.split(',')[1]),
          };
          // console.log("SKU BANK", skuBank[caseSite], SKU)
          if (currentBody.priceRangeMax == NaN || !currentBody.priceRangeMax) {
            console.log('No Max Price Range Detected');
            currentBody.priceRangeMax = 100000;
          }
          if (currentBody.priceRangeMin == NaN || !currentBody.priceRangeMin) {
            console.log('No Min Price Range Detected');
            currentBody.priceRangeMin = 1;
          }

          switch (site.toUpperCase()) {
            case 'TARGET':
            case 'GAMESTOP':
            case 'BESTBUY':
            case 'BIGLOTS':
            case 'ACADEMY':
            case 'AMD':
            case 'SLICKDEALS':
            case 'SLICK':
            case 'SLICKDEAL':
            case 'BIGLOTS':
            case 'HOMEDEPOT':
            case 'SHOPIFY':
            case 'FANATICSNEWPRODUCTS':
              await pushSku(currentObj);
              startGoMonitor(currentBody, site.toUpperCase());
              break;
            case 'NEWEGG': {
              await pushSku(currentObj);
              currentBody['skuName'] = await getSku(SKU, proxyList);
              console.log(currentBody);
              startGoMonitor(currentBody, site.toUpperCase());
              break;
            }
            // case 'AMAZON': {
            //   await pushSku(currentObj);
            //   let oid = SKU.split(":")[1]
            //   let asin = SKU.split(":")[0]
            //   currentBody['skuName'] = oid
            //   currentBody['sku'] = asin
            //   console.log(currentBody, SKU.split(":"));
            //   startGoMonitor(currentBody, site.toUpperCase());
            //   break;
            // }
            case 'SHOPIFYLINK': {
              // currentObj.sku = currentObj.sku.split("https://")[1]
              // currentBody.sku = currentBody.sku.split("https://")[1]
              await pushSku(currentObj);
              startGoMonitor(currentBody, site.toUpperCase());
              break;
            }

            case 'WALMART': {
              await pushSku(currentObj);
              console.log(currentBody);
              startGoMonitor(currentBody, site.toUpperCase());
              await delay(10000);
              break;
            }
            case 'TARGETNEW': {
              await pushSku(currentObj);
              console.log(kw);
              let currentBody = {
                endpoint: SKU,
                keywords: kw,
              };
              console.log(currentBody);
              startGoMonitor(currentBody, site.toUpperCase());
              break;
            }
            default:
          }
        }
        await delay(7500);
      }
    }
    message.channel.send('Mass Add Sequence Completed!');
  } else {
    message.channel.send(`${message.author} is not a validated user`);
  }
}

async function pushSku(body) {
  try {
    console.log('PUSHING ', body);
    let pushSku = await rp.post({
      url: `${firstServer}:${port}/UPDATESKU`,
      body: JSON.stringify(body),
    });
    console.log(pushSku?.statusCode);
  } catch (error) {
    console.log(error.message);
  }
}

module.exports = mass;
