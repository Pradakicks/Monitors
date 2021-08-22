const { MessageAttachment, MessageEmbed } = require('discord.js');
const delay = require('delay');
const fs = require('fs').promises;
const config = require('./config.json');
const os = require('os');
//  var skuBank = []
let pushEndpoint = 'https://monitors-9ad2c-default-rtdb.firebaseio.com/monitor';
let discordIds =
  'https://monitors-9ad2c-default-rtdb.firebaseio.com/validatedUsers';
const rp = require('request-promise').defaults({
  followAllRedirects: true,
  resolveWithFullResponse: true,
  gzip: true,
});
const port = 7243;
// const secondServer = `http://ec2-3-236-148-149.compute-1.amazonaws.com`

let firstServer = `http://104.249.128.37`; // 12 Core z 24GB
let secondServer = `http://104.249.128.207`;
var proxyList = [];
// let thirdSfirstServererver = `http://64.227.28.51`;
getProxies();
function SKUADD(clients, triggerText, replyText) {
  try {
    clients.on('message', async (message) => {
      if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
        // message.author.send(replyText);
        let pricerange = '';
        let kw = [];
        const content = message.content;
        const { group, isValidated } = await checkIfUserValidated(message);
        console.log(group);
        if (isValidated) {
          const site = content.split(' ')[1];
          // let spaces = content.split(' ')
          // if(spaces.length)
          let original = content.split(`${site} `)[1];
          let SKU = content.split(' ')[2];
          // if (content.includes('[') && site.toUpperCase() !== 'TARGETNEW') {
          //   pricerange = content.split('[')[1].split(']')[0];
          // }

          // if (site.toUpperCase() == 'TARGETNEW') {
          //   let kwArray = content.split('[')[1].split(']')[0];
          //   console.log(kwArray);
          //   let eachItem = kwArray.split(',');
          //   eachItem.forEach((e) => {
          //     kw.push(e);
          //   });
          // }

          console.log(site);
          console.log(SKU);
          console.log(content);
          console.log(original);
          console.log(pricerange);
          if (SKU.length > 1 && site.length > 1) {
            let isContinue = true;
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
                return;
              }
            });

            if (!shouldContinue) return;
            // ------------------------- ---------------- -------- ------------------------ --------------------

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
            if (skuBank[caseSite]) {
              if (skuBank[caseSite][SKU]) {
                let skuWebhookArray = skuBank[caseSite][SKU]?.companies;
                let isPresent = false;
                skuWebhookArray?.forEach((e) => {
                  console.log(e.webhook, 'CURRENT', currentCompany[caseSite]);
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
                  skuBank[caseSite][SKU].companies = arr;
                  await updateSku(site, SKU, skuBank[caseSite][SKU]);
                }
                console.log('Duplicate Found', isPresent);
                isContinue = false;
              }
            }

            if (
              currentBody.priceRangeMax == NaN ||
              !currentBody.priceRangeMax
            ) {
              console.log('No Max Price Range Detected');
              currentBody.priceRangeMax = 100000;
            }
            if (
              currentBody.priceRangeMin == NaN ||
              !currentBody.priceRangeMin
            ) {
              console.log('No Min Price Range Detected');
              currentBody.priceRangeMin = 1;
            }

            if (isContinue) {
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
                case 'NEWEGG':
                  await pushSku(currentObj);
                  currentBody['skuName'] = await getSku(SKU, proxyList);
                  console.log(currentBody);
                  startGoMonitor(currentBody, site.toUpperCase());
                case 'WALMART':
                  await pushSku(currentObj);
                  console.log(currentBody);
                  startGoMonitor(currentBody, site.toUpperCase());
                  await delay(10000);
                  break;
                case 'SHOPIFYPRODUCT': {
                  currentBody['skuName'] = SKU.toUpperCase().split('_')[1];
                  currentBody.sku = SKU.toUpperCase().split('_')[0];
                  currentObj.sku = SKU.toUpperCase().split('_')[1];
                  await pushSku(currentObj);
                  console.log(currentBody);
                  startGoMonitor(currentBody, site.toUpperCase());
                  break;
                }
                case 'WALMARTNEW':
                case 'WALMART NEW':
                  currentBody['skuName'] =
                    'prg=desktop&facet=retailer:Walmart.com&sort=new';
                  await pushSku(currentObj);
                  await startGoMonitor(currentBody, site.toUpperCase());
                  currentBody['skuName'] =
                    'prg=desktop&facet=retailer%3AWalmart.com&soft_sort=false&sort=new';
                  await startGoMonitor(currentBody, site.toUpperCase());
                  currentBody['skuName'] =
                    'prg=desktop&cat_id=0&facet=brand%3APanini%7C%7Cbrand%3ATopps%7C%7Cretailer%3AWalmart.com&grid=false&query=panini&soft_sort=false&sort=new';
                  startGoMonitor(currentBody, site.toUpperCase());
                case 'TARGETNEW':
                  await pushSku(currentObj);
                  console.log(kw);
                  console.log({
                    endpoint: SKU,
                    keywords: kw,
                  });
                  startGoMonitor(
                    {
                      endpoint: SKU,
                      keywords: kw,
                    },
                    site.toUpperCase()
                  );
                default:
              }
              message.channel.send(`${SKU} Added to ${site}`);
            }
          }
        } else {
          message.channel.send(`${message.author} is not a validated user`);
        }
      }
    });
  } catch (error) {
    console.log(error);
  }
}
function findCommand(clients, triggerText, replyText) {
  clients.on('message', (message) => {
    if (message.content.toLowerCase() === triggerText.toLowerCase()) {
      message.author.send(replyText);
    }
  });
}

function deleteSku(clients, triggerText, replyText) {
  try {
    clients.on('message', async (message) => {
      if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
        const content = message.content;
        const site = content.split(' ')[1];

        let spaceLength = content.split(" ")
        console.log(`Space Length : ${spaceLength} : ${spaceLength.length}`)
        spaceLength.splice(0, 1)
        spaceLength.splice(0, 1)
        console.log(`Space : ${spaceLength} : ${spaceLength.length}`)

        spaceLength.forEach(async SKU => {
          console.log(site);
          console.log(`SKU - ${SKU}`);
          console.log(content);
          let skuBank = await getSkuBank();
          let caseSite = site.toUpperCase();
  
          let currentBody = skuBank[caseSite][SKU];
          const { group, isValidated } = await checkIfUserValidated(message);
  
          if (isValidated) {
            if (currentBody === undefined) {
              console.log('undefined body');
              message.channel.send(
                `${SKU} Not Present \nCannot Delete ${SKU} from ${replaceWithTheCapitalLetter(
                  site
                )}`
              );
            } else {
              console.log('Current length');
              console.log(currentBody.companies.length);
              let currentMessage = await message.channel.send(
                `Deleting ${SKU} from ${replaceWithTheCapitalLetter(site)}...`
              );
              if (currentBody.companies.length > 1) {
                for (let i = 0; i < currentBody.companies.length; i++) {
                  if ((currentBody.companies[i].company = group)) {
                    console.log(currentBody.companies);
                    currentBody.companies.splice(i, 1);
                    console.log(currentBody.companies);
                    skuBank[caseSite][SKU].companies = currentBody.companies;
                    await updateSku(site, SKU, skuBank[caseSite][SKU]);
                  }
                  currentMessage.edit(
                    `${SKU} Deleted From ${replaceWithTheCapitalLetter(site)}`
                  );
                }
              } else {
                currentBody.stop = true;
                console.log(currentBody);
                if (group == currentBody.companies[0].company) {
                  console.log('Perm Delete');
  
                  await updateSku(site, SKU, currentBody);
                  await delay(10000);
                  await deleteSkuEnd(site, SKU);
                  currentMessage.edit(
                    `${SKU} Deleted From ${replaceWithTheCapitalLetter(site)}`
                  );
                } else {
                  console.log(`${group} is not present for this sku`);
                  currentMessage.edit(
                    `${SKU} Not Present \nCannot Delete ${SKU} from ${replaceWithTheCapitalLetter(
                      site
                    )}`
                  );
                  currentMessage.edit(
                    `If this is an error please contact developer`
                  );
                }
              }
            }
          } else {
            message.channel.send(`${message.author} is not a validated user`);
          }
        })
        // const SKU = content.split(' ')[2];
        return;
      }
    });
  } catch (error) {
    console.log(error);
  }
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

async function startGoMonitor(currentBody, site) {
  try {
    switch (site) {
      case 'BESTBUY': {
        rp.post(
          {
            url: `${secondServer}:${port}/${site}`,
            body: JSON.stringify(currentBody),
            headers: {
              'Content-Type': 'application/json',
            },
          },
          (response) => console.log(response?.statusCode)
        );
        break;
      }
      default: {
        rp.post(
          {
            url: `${firstServer}:${port}/${site}`,
            body: JSON.stringify(currentBody),
            headers: {
              'Content-Type': 'application/json',
            },
          },
          (response) => console.log(response?.statusCode)
        );
      }
    }
  } catch (error) {
    console.log(`Error Starting Go Monitor ${error}`);
  }
}

async function mass(string, content, message, groupName) {
  let spaceLength = content.split(" ")
  console.log(`Space Length : ${spaceLength} : ${spaceLength.length}`)
  spaceLength.splice(0, 1)
  spaceLength.splice(0, 1)
  console.log(`Space : ${spaceLength} : ${spaceLength.length}`)
  const site = content?.split(' ')[1]
  // console.log(content?.split(' ')[1].length);
  // console.log(site.toUpperCase().length);
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
    let g = spaceLength
    //	console.log(g)
    for (let i = 0; i < g.length; i++) {
      if (!g[i].toUpperCase().includes('!MASSADD')) {
        let isContinue = true;
        let SKU;
        let pricerange = '';
        SKU = g[i];
        let original = g[i];
        if (g[i].includes('[')) {
          pricerange = g[i].split('[')[1].split(']')[0];
          SKU = g[i].split(' ')[0];
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
            case 'BIGLOTS': {
              await pushSku(currentObj);
              console.log(currentBody);
              startGoMonitor(currentBody, site.toUpperCase());
              break;
            }
            case 'NEWEGG': {
              await pushSku(currentObj);
              currentBody['skuName'] = await getSku(SKU, proxyList);
              console.log(currentBody);
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
            case 'WALMARTNEW':
            case 'WALMART NEW': {
              console.log(pricerange);
              currentBody['skuName'] =
                'prg=desktop&cat_id=0&facet=brand%3APanini%7C%7Cbrand%3ATopps%7C%7Cretailer%3AWalmart.com&grid=false&query=panini&soft_sort=false&sort=new';
              await pushSku(currentObj);
              console.log(currentBody);
              startGoMonitor(currentBody, site.toUpperCase());
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
// Fire Base Sku Bank ----------------------------------------------
async function checkPresentSkus() {
  let skuBank = await rp.get({
    url: `${secondServer}:${port}/DB`,
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
  secondServer = 'http://localhost';
  firstServer = 'http://localhost';
} else {
  checkPresentSkus();
}

async function getSkuBank() {
  let getbank = await rp.get({
    url: `${secondServer}:${port}/DB`,
  });
  return JSON.parse(getbank?.body);
}
async function pushSku(body) {
  try {
    console.log('PUSHING ', body);
    let pushSku = await rp.post({
      url: `${secondServer}:${port}/UPDATESKU`,
      body: JSON.stringify(body),
    });
    console.log(pushSku?.statusCode);
  } catch (error) {
    console.log(error.message);
  }
}
async function deleteSkuEnd(site, sku, group) {
  try {
    console.log(`Deleting ${sku}/${site}`);
    let deleteSku = await rp.post({
      url: `${secondServer}:${port}/DELETESKU`,
      body: JSON.stringify({ site: site.toUpperCase(), sku: sku }),
    });
    console.log(deleteSku?.statusCode);
  } catch (error) {
    console.log(error);
  }
}
async function updateSku(site, sku, newBody) {
  try {
    // No need for site and sku // Only reason I kept it here is for console logs
    console.log(`Updating Sku ${sku}/${site}`);
    let updateSku = await rp.post({
      url: `${secondServer}:${port}/UPDATESKU`,
      body: JSON.stringify(newBody),
    });
    console.log(updateSku?.statusCode);
  } catch (error) {
    console.log(error);
  }
}
async function getValidatedIds() {
  try {
    let getIds = await rp.get({
      url: `${secondServer}:${port}/DISCORDIDS`,
    });
    console.log(getIds?.statusCode);
    return getIds?.body;
  } catch (error) {
    console.log(error);
  }
}
async function updateDiscordIdsDB(author, discordIdsArr, name) {
  try {
    let currentIds = [];
    let parsed = JSON.parse(discordIdsArr);
    var isPresent = false;
    // console.log(parsed, author, discordIdsArr, name);
    parsed?.forEach((e) => {
      console.log(e);
      currentIds.push(e);
      let id = e?.split('-')[0];
      if (id == author) isPresent = true;
    });
    if (!isPresent) currentIds.push(`${author}-${name}`);
    else return false;
    console.log({ ids: currentIds });
    let updateIds = await rp.post({
      url: `${secondServer}:${port}/DISCORDIDS`,
      body: JSON.stringify({ ids: currentIds }),
    });
    console.log(updateIds?.statusCode);
    return true;
  } catch (error) {
    console.log(error?.message);
    return false;
  }
}
async function validateUser(clients, triggerText, replyText) {
  try {
    clients.on('message', async (message) => {
      if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
        let content = message.content.split('!validate')[1];
        let apikey = content?.split('apikey-')[1]
        if (apikey == undefined) {
          message.reply('Please Submit Valid Api Key!');
        } else {
          console.log(apikey);
          let isPresent = false;
          const discordIdsDB = await getValidatedIds();
          Object.keys(config.groups).forEach(async (e) => {
            if (config.groups[e].apiKey == apikey) {
              console.log(config.groups[e]);
              message.reply('Valid Api Key');
              isPresent = true;
              var update = await updateDiscordIdsDB(
                message.author.id,
                discordIdsDB,
                e
              );
              if (update) message.reply(`${message.author} Validated`);
              else
                message.reply(
                  `${message.author} Could not be Validated \n Contact Dev if this continues to be an issue`
                );
            }
          });

          if (!isPresent) {
            message.reply('Unknown Api Key');
            console.log(apikey, config.groups);
          }
        }
        console.log(content);
      }
    });
  } catch (error) {
    console.log(error);
  }
}
async function checkIfUserValidated(message) {
  let returnObj = {
    isValidated: false,
    group: undefined,
  };
  let validatedIds = await getValidatedIds();
  let parsed = JSON.parse(validatedIds);
  // console.log(parsed)
  parsed.forEach((e) => {
    let id = e?.split('-')[0];
    console.log(id, message.author.id);
    if (id == message.author.id) {
      returnObj.isValidated = true;
      returnObj.group = e?.split('-')[1];
    }
  });
  return returnObj;
}
//-----------------------------------------------------------------

// Scraper
async function walmartScraper(clients, triggerText, replyText) {
  try {
    clients.on('message', async (message) => {
      if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
        const { group, isValidated } = await checkIfUserValidated(message);
        let currentCompany = config.groups[group];
        if (isValidated) {
          console.log('Walmart Scraper Initiated');
          const SKU = message.content.split('/')[5];
          let currentMessage = await message.channel.send(
            `${message.author} Walmart Scraper for ${SKU} Initiated`
          );
          let data = await getTerraSku(SKU);
          if (data == 'Error') {
            currentMessage.edit(
              `${message.author} Error Grabbing Data for ${SKU}`
            );
            return;
          }
          try {
            let ImageOffer = data.payload.selected.defaultImage;
            let selectedProduct = data.payload.selected.product;
            let productName =
              data.payload.products[selectedProduct].productAttributes
                .productName;
            let ImageUrl =
              data.payload.images[ImageOffer].assetSizeUrls.DEFAULT;
            let walmartSellerOffer;

            Object.keys(data.payload.sellers).map((e) => {
              let currentSeller = data?.payload?.sellers[e];
              if (currentSeller?.sellerName == 'Walmart.com')
                walmartSellerOffer = currentSeller?.sellerId;
            });

            console.log(ImageUrl, walmartSellerOffer);
            let offerObj;
            if (!walmartSellerOffer) {
              currentMessage.edit(
                `${message.author} ${SKU} is not sold by Walmart`
              );
              offerObj = {
                offerId: 'N/A',
                availabilityStatus: 'Not Sold By Walmart',
                price: 'N/A',
                image: ImageUrl,
                sku: SKU,
                productName: productName,
              };
            } else {
              Object.keys(data.payload.offers).forEach((e) => {
                let currentOffer = data.payload.offers[e];
                console.log(currentOffer?.pricesInfo);
                if (currentOffer.sellerId == walmartSellerOffer)
                  offerObj = {
                    offerId: currentOffer.id ? currentOffer.id : 'N/A',
                    availabilityStatus: currentOffer?.productAvailability
                      ?.availabilityStatus
                      ? currentOffer.productAvailability.availabilityStatus
                      : 'N/A',
                    price: currentOffer?.pricesInfo?.priceMap?.CURRENT?.price
                      ? currentOffer.pricesInfo.priceMap.CURRENT.price
                      : 'No Price Found',
                    image: ImageUrl,
                    sku: SKU,
                    productName: productName,
                  };
              });
            }

            var dt = new Date();
            // let body = {
            // 	"content": null,
            // 	"embeds": [
            // 	  {
            // 		"title": "Walmart Scraper",
            // 		"color": currentCompany.companyColor,
            // 		"fields": [
            // 		  {
            // 			"name": "Offer Id",
            // 			"value": offerObj.offerId
            // 		  },
            // 		  {
            // 			"name": "Availability",
            // 			"value": offerObj.availabilityStatus,
            // 			"inline": true
            // 		  },
            // 		  {
            // 			"name": "Walmart's Price",
            // 			"value": offerObj.price,
            // 			"inline": true
            // 		  },
            // 		  {
            // 			"name": "Links",
            // 			"value": `[Product](https://www.walmart.com/ip/prada/${offerObj.sku}) | [ATC](https://affil.walmart.com/cart/buynow?items=${offerObj.sku}) | [Checkout](https://www.walmart.com/checkout/) | [Cart](https://www.walmart.com/cart)`
            // 		  }
            // 		],
            // 		"footer": {
            // 		  "text": "Prada#4873"
            // 		},
            // 		"timestamp": dt.toISOString(),
            // 		"thumbnail": {
            // 		  "url": offerObj.image
            // 		}
            // 	  }
            // 	],
            // 	"avatar_url": currentCompany.companyImage
            //   }
            const currentEmbed = new MessageEmbed()
              .setColor(currentCompany.companyColorV2)
              .setTitle(offerObj.productName)
              .setURL(`https://www.walmart.com/ip/prada/${offerObj.sku}`)
              .setThumbnail(offerObj.image)
              .addFields(
                {
                  name: 'Offer Id',
                  value: offerObj.offerId,
                },
                {
                  name: 'PID',
                  value: offerObj.sku,
                },
                {
                  name: 'Availability',
                  value: offerObj.availabilityStatus,
                  inline: true,
                },
                {
                  name: "Walmart's Price",
                  value: offerObj.price,
                  inline: true,
                },
                {
                  name: 'Links',
                  value: `[Product](https://www.walmart.com/ip/prada/${offerObj.sku}) | [ATC](https://affil.walmart.com/cart/buynow?items=${offerObj.sku}) | [Checkout](https://www.walmart.com/checkout/) | [Cart](https://www.walmart.com/cart)`,
                }
              )
              .setFooter('Prada#4873')
              .setTimestamp();
            message.channel.send(currentEmbed);
            currentMessage.delete();
          } catch (error) {
            console.log(error);
            currentMessage.edit(
              `${message.author} Error Grabbing Data for ${SKU} BOJ`
            );
          }
        } else {
          message.channel.send(`${message.author} is not a validated user`);
        }
      }
    });
  } catch (error) {
    console.log(error);
  }
}
async function getSku(skuName, proxyList) {
  try {
    let proxy1 = proxyList[Math.floor(Math.random() * proxyList.length)];
    console.log(proxy1);
    let fetchProductPage = await rp.get({
      url: `https://www.newegg.com/prada/p/${skuName}`,
      headers: {
        accept:
          'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
        'accept-language': 'en-US,en;q=0.9',
        'cache-control': 'no-cache',
        pragma: 'no-cache',
        'sec-fetch-dest': 'document',
        'sec-fetch-mode': 'navigate',
        'sec-fetch-site': 'none',
        'sec-fetch-user': '?1',
        'upgrade-insecure-requests': '1',
      },
      proxy: `http://${proxy1.userAuth}:${proxy1.userPass}@${proxy1.ip}:${proxy1.port}`,
    });
    console.log(fetchProductPage?.statusCode);
    sku =
      fetchProductPage?.body?.split('/ProductImage/')[1].split('-')[0] + '-';
    console.log(sku);
    sku =
      sku +
      fetchProductPage?.body?.split('/ProductImage/')[1].split('-')[1] +
      '-';
    console.log(sku);
    sku =
      sku + fetchProductPage?.body?.split('/ProductImage/')[1].split('-')[2];
    console.log(sku);
    return sku;
  } catch (error) {
    console.log(error);
    await delay(10000);
    await getSku(skuName, proxyList);
  }
}
async function getTerraSku(SKU) {
  try {
    // let proxies = await getProxies()
    const proxy1 = proxyList[Math.floor(Math.random() * proxyList.length)];
    console.log(proxy1);
    let terra = await rp.get({
      url: `https://www.walmart.com/terra-firma/item/${SKU}`,
      headers: {
        accept:
          'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
        'accept-language': 'en-US,en;q=0.9',
        'cache-control': 'no-cache',
        pragma: 'no-cache',
        'sec-ch-ua':
          '" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"',
        'sec-ch-ua-mobile': '?0',
        'sec-fetch-dest': 'document',
        'sec-fetch-mode': 'navigate',
        'sec-fetch-site': 'none',
        'sec-fetch-user': '?1',
        'service-worker-navigation-preload': 'true',
        'upgrade-insecure-requests': '1',
        'user-agent':
          'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36',
      },
      proxy: `http://${proxy1.userAuth}:${proxy1.userPass}@${proxy1.ip}:${proxy1.port}`,
    });
    console.log(`Terra Status Code : ${terra.statusCode}`);
    if (terra.statusCode != 200 && terra.statusCode != 201) {
      return 'Error';
    } else {
      let parsed = JSON.parse(terra?.body);
      return parsed;
    }
  } catch (error) {
    console.log(error);
    return 'Error';
  }
}
// Helper
function replaceWithTheCapitalLetter(values) {
  return values.charAt(0).toUpperCase() + values.slice(1);
}
async function getProxies() {
  try {
    // read contents of the file
    // const data = await fs.readFile('./GoMonitor/cloud.txt', 'utf-8');
    let fetchProxies = await rp.get(`${secondServer}:${port}/PROXY`);
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
    fs.appendFileSync('./errors.txt', err.toString() + '\n', (err) => {
      console.log(err);
    });
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
  SKUADD,
  findCommand,
  deleteSku,
  checkBank,
  massAdd,
  validateUser,
  walmartScraper,
};
