const chalk = require('chalk');
const clear = require('clear');
const figlet = require('figlet');
const files = require('./lib/files');
const Configstore = require('configstore');
const inquirer = require('./lib/inquirer');

const config = new Configstore('VibrisCLI')
const { retrieveLicenseCLI , updateLicense} = require('./lib/keyauth')
const fs = require('fs');
const crypto = require('crypto');
const { taskWorkersCLI } = require('./lib/taskWorker.js');
const delay = require('delay');
const captcha = require('./modules/footsites/utils/captcha');
const algorithm = 'aes-256-ctr';
const secretKey = 'vOVH6sdmpNWjRRIqCc7rdxs01lwHzfr3';
const iv = crypto.randomBytes(16);


function encrypt (text) {

    const cipher = crypto.createCipheriv(algorithm, secretKey, iv);

    const encrypted = Buffer.concat([cipher.update(text), cipher.final()]);

    return {
        iv: iv.toString('hex'),
        content: encrypted.toString('hex')
    };
};

function decrypt (hash) {

    const decipher = crypto.createDecipheriv(algorithm, secretKey, Buffer.from(hash.iv, 'hex'));

    const decrpyted = Buffer.concat([decipher.update(Buffer.from(hash.content, 'hex')), decipher.final()]);

    return decrpyted.toString();
};

var taskVisibility = [{
    task : true
}]

try {
    

    const run = async () => {
        try {
            clear();
            consoleIt = function(t) {
                process.stdout.write(t + '\n')
            }
    
            consoleIt(chalk.yellow(figlet.textSync('Vibris Cli', { horizontalLayout: 'full' })))

            async function checkKey() {
                var currentKey = ''
                let cryptedKey = config.get('condfapldeudfe')
                if (cryptedKey?.iv) {
                    let decrytedKey = decrypt(cryptedKey)
                    currentKey = decrytedKey
                    let response = await retrieveLicenseCLI(currentKey.toString())
                    if (response) {
                        consoleIt(chalk.green('Success'))
                    } else {
                        consoleIt(chalk.red('Key Rejected'))
                        await delay(2500)
                        await checkKey()
                       
                    }
                } else {
                    let keyAuth = await inquirer.keyAuth()
                    let realFilePath = keyAuth.key.split("\\").join("/");
                    if (realFilePath.includes('"') || realFilePath.includes("'")){
                        realFilePath = realFilePath.split('"').join('')
                        realFilePath = realFilePath.split("'").join('')
                    }
                    let fileContent = fs.readFileSync(`${realFilePath}`, 'utf-8');
                    currentKey = fileContent
                    config.set('condfapldeudfe', encrypt(fileContent))
                    let response = await retrieveLicenseCLI(currentKey)
                    if (response) {
                        consoleIt(chalk.green('Success'))
                    } else {
                        consoleIt(chalk.red('Key Rejected'))
                        await checkKey()
                    }
    
                }
    
                
               
         
      
              
            }  await checkKey()
            await inquirer.turnOffAYCD()
            async function importTask() {
                let yesOrNo = await inquirer.yesTasks()
                if (yesOrNo) {
                    let task = await inquirer.importTasks()
                    let realFilePath = task.tasks.split("\\").join("/");
                    if (realFilePath.includes('"') || realFilePath.includes("'")) {
                        realFilePath = realFilePath.split('"').join('')
                        realFilePath = realFilePath.split("'").join('')
                    }
                    let fileContent = fs.readFileSync(`${realFilePath}`, 'utf-8');
                    config.set('tasks', fileContent);
                    consoleIt(chalk.green('Tasks were added'))
                } else {
                    return
                }
              
            } //
            async function importProxy() {
                let yesOrNo = await inquirer.YesProxy()
                if (yesOrNo) {
                    let proxy = await inquirer.importProxies()
                    let realFilePath = proxy.proxies.split("\\").join("/");
                    if (realFilePath.includes('"') || realFilePath.includes("'")) {
                        realFilePath = realFilePath.split('"').join('')
                        realFilePath = realFilePath.split("'").join('')
                    }
                    let fileContent = fs.readFileSync(`${realFilePath}`, 'utf-8');
                    //  consoleIt(fileContent)
                    let arr = []
                    const data = fs.readFileSync(`${realFilePath}`, 'utf-8');
                    const lines = data.split(/\r?\n/);
                    await lines.forEach((line) => {
                        arr.push(line)
                    });
                 //   consoleIt(arr)
                    config.set('proxies', arr);
                    consoleIt(chalk.green('Proxies were added'))
                } else {
                    return
                }
                
            } // 
            async function importWebhook() {
                let yesOrNo = await inquirer.yesOrNo()
           //     consoleIt(yesOrNo)
                if (yesOrNo) {
                    let webhook = await inquirer.setWebhook()
                    let realFilePath = webhook.webhook.split("\\").join("/");
                    if (realFilePath.includes('"') || realFilePath.includes("'")) {
                        realFilePath = realFilePath.split('"').join('')
                        realFilePath = realFilePath.split("'").join('')
                    }
                    let fileContent = fs.readFileSync(`${realFilePath}`, 'utf-8');
                    config.set('webhook', fileContent);
                    consoleIt(chalk.green('Webhook Webhooks'))
                } else {
                    return
                }
                
            } // 
            async function askingForCaptcha () {
                let captcha = await inquirer.askForCaptcha()
           //     consoleIt(yesOrNo)
                if (captcha) {
                    if(captcha.askForCaptcha == 'AYCD'){
                        let aycd = await inquirer.aycd()
                        if(aycd){
                            config.set('captcha', "AYCD/" + aycd.key);
                            consoleIt(chalk.green('AYCD Added'))
                       }
                    } else if (captcha.askForCaptcha == '2Cap') {
                        let twoCap = await inquirer.twoCap()
                        if(twoCap){
                            config.set('captcha', "TwoCap/" + twoCap.key);
                            consoleIt(chalk.green('Two Cap Added'))
                       }
                    }  else if (captcha.askForCaptcha == 'Cap Monster') {
                        let capM = await inquirer.capMonster()
                        if(capM){
                            config.set('captcha', "capMonster/" + capM.key);
                            consoleIt(chalk.green('Cap Monster Added'))
                       }
                    } else {
                        return;
                    }
                } else {
                    return;
                }
            } //
            async function checkImportedTasks() {
                let tasks = config.get('tasks')
                let proxy = config.get('proxies')
                let cm = config.get('captcha')
                let captchaMethod = cm.split('/')[0]
                let parsedTasks = JSON.parse(tasks)
                let parsedProxies = proxy.length
              //  consoleIt(parsedTasks)
                consoleIt(chalk.blueBright(`Number of Tasks: ${parsedTasks.length}`))
                consoleIt(chalk.blueBright(`Number of Proxy: ${parsedProxies}`))
                consoleIt(chalk.blueBright(`Captcha Method: ${captchaMethod}`))
            } // 
            async function checkCon() {
            let confirmatinon = await inquirer.checkConfirmation()
             //  consoleIt(confirmatinon)
            } // await checkCon()
            async function mainMenu() {
                const readline = require('readline');
                readline.emitKeypressEvents(process.stdin);
                process.stdin.setRawMode(true);
                const keyMap = new Map();
                keyMap.set('C', 'clear');
                process.stdin.on('keypress', (str, key) => {
                   // console.log(key)
                   
                    if (key.sequence === '\x07') {
                        if(!taskVisibility[0].task){
                            taskVisibility[0].task = true
                            } else {
                                taskVisibility[0].task = false
                            }
                        console.log(taskVisibility[0].task)
                            return
                    } else {
                        return;
                      
                    }
                  });
                let menu = await inquirer.mainMenu()
               // consoleIt(menu)
                if (menu.mainMenu == 'Start Tasks') {
                    await checkImportedTasks()
                    await checkCon()
                    await startTasks()
                } else if (menu.mainMenu == 'Import Settings') {
                    await importTask()
                    await importProxy()
                    await importWebhook()
                    await askingForCaptcha()
                    await checkImportedTasks()
                    await mainMenu()
                }
            } await mainMenu()
            async function startTasks() {
                try {
                    let tasks = config.get('tasks')
                    let parsedTasks = JSON.parse(tasks)
                    let cm = config.get('captcha')
                    let captchaMethod = cm.split('/')[0]
                    let captchaKey = cm.split('/')[1]
                  //  console.log(captchaKey)
                    if(captchaMethod == 'AYCD'){
                        const { autoSolveInit } = require('./modules/footsites/utils/captcha/aycd')
                        let accessToken = captchaKey?.split(':')[0]
                        let apikey = captchaKey?.split(':')[1]
                        console.log(`AYCD ACCESS TOKEN : ${accessToken}`)
                        console.log(`AYCD API KEY: ${apikey}`)
                        autoSolveInit(accessToken, apikey)
                    }
                    for (let i = 0; i < parsedTasks.length; i++){
                        taskWorkersCLI(parsedTasks[i])
                        await delay(1)
                    }
                } catch (error) {
                    console.log(error)
                }
             
               
            } 

        } catch (error) {

            fs.appendFileSync('errorCli.txt', error.toString() + '\n', (err => {
                console.log(err)
            }))
            console.log(error)
        }
      
      //  await checkKey()
        
    };
    run();  

} catch (error) {

    fs.appendFileSync('errorCli.txt', error.toString() + '\n', (err =>{
        console.log(err)
    }))
}



module.exports = {
    encrypt,
    decrypt,
    taskVisibility
}

