const footsites = require('../modules/footsites');
const Configstore = require('configstore');
const config = new Configstore('VibrisCLI')
const crypto = require('crypto');
async function logTask(store, module, kw, id) {
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
    let cryptedKey = config.get('condfapldeudfe')
    let decrytedKey = decrypt(cryptedKey)
    var key = decrytedKey;

    let date_ob = new Date();
    let date = ("0" + date_ob.getDate()).slice(-2);
    let month = ("0" + (date_ob.getMonth() + 1)).slice(-2);
    let year = date_ob.getFullYear();
    let currentDate = (year + "-" + month + "-" + date + '/');
    let baseURL = 'https://vibris-e33a3-default-rtdb.firebaseio.com/'
    let taskId = {
        type : 'TASK',
        store: store,
        module: module + 'CLI',
        currentDate: currentDate,
        kw : kw,
        id:Date.now()
    }
    let targetUrl = `${baseURL}users/${key}.json`
   // console.log(targetUrl)
    rp.post({
        url: targetUrl,
        body: JSON.stringify(taskId)
    })

}

const taskWorkersCLI = async (task) => {
    try {
        logTask(task.store.longCode, task.mode, task.keywords, task.id)
        switch (task.store.longCode.toUpperCase()) {
            case ('footlocker').toUpperCase():
            case ('champssports').toUpperCase():
            case ('footaction').toUpperCase():
            case ('eastbay').toUpperCase():
            case ('footlockerCa').toUpperCase():
            case ('kidsfootlocker').toUpperCase():
                await footsites(task)
                break;

            default:
                return;
        }
    } catch (error) {
        console.log(error)
    }
}

module.exports = {
    taskWorkersCLI
}