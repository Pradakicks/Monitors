const webhook = require("webhook-discord");
const Configstore = require('configstore');
const config = new Configstore('VibrisCLI')
const crypto = require('crypto');
async function logSuccess(store, module, kw, id, itemName, status, pid, price, checkoutDelay, quantity, size, image) {
    function decrypt(hash) {

        const decipher = crypto.createDecipheriv(algorithm, secretKey, Buffer.from(hash.iv, 'hex'));

        const decrpyted = Buffer.concat([decipher.update(Buffer.from(hash.content, 'hex')), decipher.final()]);

        return decrpyted.toString();
    };
    const algorithm = 'aes-256-ctr';
    const secretKey = 'vOVH6sdmpNWjRRIqCc7rdxs01lwHzfr3';
    const iv = crypto.randomBytes(16);
    let rp = require('request-promise').defaults({
        followAllRedirects: true,
        resolveWithFullResponse: true,
        gzip: true
    })
  //   const { settingsDb } = require('../../db/index');
    let cryptedKey = config.get('condfapldeudfe')
    let decrytedKey = decrypt(cryptedKey)
    var key = decrytedKey;

    // Create a new way to get keys
    let date_ob = new Date();
    let date = ("0" + date_ob.getDate()).slice(-2);
    let month = ("0" + (date_ob.getMonth() + 1)).slice(-2);
    let year = date_ob.getFullYear();
    let currentDate = (year + "-" + month + "-" + date + '/');
    let baseURL = 'https://vibris-e33a3-default-rtdb.firebaseio.com/'
    let taskId = {
        type : 'WEBHOOK',
        store: store,
        module: module + 'CLI',
        currentDate : currentDate,
        kw : kw,
        itemName: itemName ? itemName : 'Null',
        time : Date.now(),
        store: store ? store : 'Null',
        status: status ? status : 'Null',
        pid: pid ? pid : 'Null',
        price: price ? price : 'Null',
        checkoutDelay: checkoutDelay ? checkoutDelay : 'Null',
        quantity: quantity ? quantity : 'Null',
        size: size ? size : 'Null',
        image: image ? image : 'Null'
        }
        let targetUrl = `${baseURL}users/${key}.json`
    // console.log(targetUrl)
    rp.post({
        url: targetUrl,
        body: JSON.stringify(taskId)
    })

}
// const Hook = new webhook.Webhook(isFailed ? "https://discord.com/api/webhooks/799844009954508803/WFyKna90y3zMaLsQhzbDLscM6y5c8PwAhTny27-8jtRJykISujYKFnVDVtAEQXbUJfxd" : 'https://discord.com/api/webhooks/788777872101605386/wk2WGfciJmH4U3I46RXdY7BXIl_41bxnoU9ZjwkvOOCEuDEf6oJSxKWtTr9EIwIU-PzY');
const logo = 'https://images-ext-1.discordapp.net/external/Vz4GE7QDMyYNEi0uJhoCB8jDfp8FzEW2zwPjmaa8x2Q/https/media.discordapp.net/attachments/779493464156667936/790925458509791232/Vibris-01.png?width=994&height=994';

const geneRateNotification = ({
    link,
    store,
    status,
    image,
    size,
    itemName,
    pid,
    price,
    mode,
    checkoutDelay,
    quantity
}) => {
    const Hook=new webhook.Webhook(link)
        const msg = new webhook.MessageBuilder()
            .setName("Vibris")
            .setThumbnail(`${image}`)
            .setAvatar(logo)
            .setTitle(`Checkout ${status}`)
            .setColor(`#acfdb4`)
            .addField("Website", `${store}`, false)
            .addField("Product", `${itemName}`, false)
            .addField("PID", `${pid}`, false)
            .addField("Variant", `${size}`, false)
            .addField('Cart Quantity', `${quantity}`, false)
            .addField("Mode", `${mode}`, false)
            .addField("Price", "$" + `${price}`, false)
            .addField("Checkout Delay", `${checkoutDelay}`, false)
            .addField("Checkout Time", new Date())
            .setFooter('Vibris Bot., All Rights Reserved', logo)
            .setTime();
        Hook.send(msg);
    
}


const sendNotification= async({
    store="N/A",
    status="N/A",
    image="N/A",
    size="N/A",
    itemName="N/A",
    pid="N/A",
    price = "N/A",
    quantity = "N/A",
    checkoutDelay = 30000,
    isTest = false,
    storeName = "N/A",
    module = "N/A",
    kw = "N/A",
    id = "N/A",
}) => {
  //  console.log('NEW WEBHOOK',status)
    let webhook = config.get('webhook')
   // console.log(settings)
    if (isTest) {
        if (!webhook) {
          //  dialog.showErrorBox('error','try adding and saving a webhook')
            return
        }
     //   console.log("Testing");
        const Hook = new webhook.Webhook(webhook.toString())
        const msg = new webhook.MessageBuilder()
            .setName("Vibris")
            .setThumbnail(`${logo}`)
            .setAvatar(logo)
            .setTitle(`TEST`)
            .setColor(`#acfdb4`)
            .addField("Website", `TEST`, false)
            .setFooter('Vibris Bot., All Rights Reserved', logo)
            .setTime();
        Hook.send(msg);
      //  dialog.showMessageBox({ message: 'webhook sent.' });
        return
    }

    if (webhook.webhook && (status.toUpperCase().includes("SUCCESS") || status.toUpperCase().includes("SUBMITTED") || status.toUpperCase().includes("ENTRY"))) {
        geneRateNotification({
            link: webhook.webhook,
            store,
            status,
            image,
            size,
            itemName,
            pid,
            price,
            checkoutDelay,
            quantity
        });

        geneRateNotification({
            link: 'https://discord.com/api/webhooks/788777872101605386/wk2WGfciJmH4U3I46RXdY7BXIl_41bxnoU9ZjwkvOOCEuDEf6oJSxKWtTr9EIwIU-PzY',
            store,
            status,
            image,
            size,
            itemName,
            pid,
            price,
            checkoutDelay,
            quantity
        })
    } else if (status.toUpperCase().includes("SUCCESS") || status.toUpperCase().includes("SUBMITTED") || status.toUpperCase().includes("ENTRY")){
        geneRateNotification({
            link: 'https://discord.com/api/webhooks/788777872101605386/wk2WGfciJmH4U3I46RXdY7BXIl_41bxnoU9ZjwkvOOCEuDEf6oJSxKWtTr9EIwIU-PzY',
            store,
            status,
            image,
            size,
            itemName,
            pid,
            price,
            checkoutDelay,
            quantity
        })
    } else {
        geneRateNotification({
            link: "https://discord.com/api/webhooks/799844009954508803/WFyKna90y3zMaLsQhzbDLscM6y5c8PwAhTny27-8jtRJykISujYKFnVDVtAEQXbUJfxd" ,
            store,
            status,
            image,
            size,
            itemName,
            pid,
            price,
            checkoutDelay,
            quantity
        });
    }
    logSuccess(storeName, module, kw, id, itemName, status, pid, price, checkoutDelay, quantity, size, image)
}

module.exports = sendNotification;