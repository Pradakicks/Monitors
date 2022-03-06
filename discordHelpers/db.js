const { rp,firstServer, port } = require('./helper')
const config = require('../config.json')
const delay = require('delay');

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
      // console.log(id, message.author.id);
      if (id == message.author.id) {
        returnObj.isValidated = true;
        returnObj.group = e?.split('-')[1];
      }
    });
    return returnObj;
  }
  async function getSkuBank() {
    let getbank = await rp.get({
      url: `${firstServer}:${port}/DB`,
    });
    return JSON.parse(getbank?.body);
  }
  async function getValidatedIds() {
    try {
      let getIds = await rp.get({
        url: `${firstServer}:${port}/DISCORDIDS`,
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
        url: `${firstServer}:${port}/DISCORDIDS`,
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
          let apikey = content?.split('apikey-')[1].trim();
          if (apikey == undefined) {
            message.reply('Please Submit Valid Api Key!');
          } else {
            console.log(`Api Key : ${apikey}`);
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
  async function updateSku(site, sku, newBody) {
    try {
      // No need for site and sku // Only reason I kept it here is for console logs
      console.log(`Updating Sku ${sku}/${site}`);
      let updateSku = await rp.post({
        url: `${firstServer}:${port}/UPDATESKU`,
        body: JSON.stringify(newBody),
      });
      console.log(updateSku?.statusCode);
    } catch (error) {
      console.log(error);
    }
  }
  function deleteSku(clients, triggerText, replyText) {
    try {
      clients.on('message', async (message) => {
        if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
          const content = message.content;
          const site = content.split(' ')[1];
  
          let spaceLength = content.split(' ');
          console.log(`Space Length : ${spaceLength} : ${spaceLength.length}`);
          spaceLength.splice(0, 1);
          spaceLength.splice(0, 1);
          console.log(`Space : ${spaceLength} : ${spaceLength.length}`);
  
          spaceLength.forEach(async (SKU) => {
            console.log(site);
            if (site.toLowerCase().includes('shopifylink')) {
              SKU = SKU.split('https://')[1];
            }
            console.log(`SKU - ${SKU}`);
            console.log(content);
            let skuBank = await getSkuBank();
            let caseSite = site.toUpperCase();
  
            let currentBody = skuBank[caseSite][SKU];
            // console.log(skuBank[caseSite]);
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
          });
          // const SKU = content.split(' ')[2];
          return;
        }
      });
    } catch (error) {
      console.log(error);
    }
  }
  function replaceWithTheCapitalLetter(values) {
    return values.charAt(0).toUpperCase() + values.slice(1);
  }
    async function deleteSkuEnd(site, sku, group) {
    try {
      console.log(`Deleting ${sku}/${site}`);
      let deleteSku = await rp.post({
        url: `${firstServer}:${port}/DELETESKU`,
        body: JSON.stringify({ site: site.toUpperCase(), sku: sku }),
      });
      console.log(deleteSku?.statusCode);
    } catch (error) {
      console.log(error);
    }
  }
module.exports = {
    checkIfUserValidated,
    getSkuBank,
    getValidatedIds,
    updateDiscordIdsDB,
    validateUser,
    updateSku,
    deleteSku
}