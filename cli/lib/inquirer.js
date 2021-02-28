const inquirer = require('inquirer');



function keyAuth () {
    let askForKey = [{
        name: 'key',
        type: 'input',
        message: 'Enter your Key:',
        validate: function (value) {
            if (value.length) {
                return true;
            } else {
                return 'Please enter your Key.';
            }
        }
    }];
    return inquirer.prompt(askForKey);
}
function importTasks() {
    let importTasks = [{
        name: 'tasks',
        type: 'input',
        message: 'Enter your Tasks:',
        validate: function (value) {
            if (value.length) {
                return true;
            } else {
                return 'Please enter your Task.';
            }
        }
    }];
    return inquirer.prompt(importTasks)
}
function importProxies () {
    let importProxy = [{
        name: 'proxies',
        type: 'input',
        message: 'Enter your Proxies:',
        validate: function (value) {
            if (value.length) {
                return true;
            } else {
                return 'Please enter your Proxy Path.';
            }
        }
    }];
    return inquirer.prompt(importProxy)
}
function checkConfirmation () {
    let checkConfirmations = [{
        name: 'confirmation',
        type: 'confirm',
        message: 'Do you want to proceed with current task information?',
    }];
    return inquirer.prompt(checkConfirmations).then(res => {
       // console.log(res)
        if (res.confirmation) {
          //  console.log(res.confirmation)
        } else {
            throw new Error('No')
        }
    })
}
function mainMenu() {
    let main = [{
        name: 'mainMenu',
        type: 'list',
        message: 'Menu :',
        choices: ['Start Tasks', 'Import Settings'],
        defeat: function (value) {
       //     console.log(value)
            if (value.length) {
                return true;
            } else {
                return false;
            }
        }
    }];
    return inquirer.prompt(main)
}
function setWebhook() {
    let setWebhook = [{
        name: 'webhook',
        type: 'input',
        message: 'Enter your Webhook:',
        validate: function (value) {
            if (value.length) {
                return true;
            } else {
                return 'Please enter your Discord Webhook.';
            }
        }
    }];
    return inquirer.prompt(setWebhook)
}
function yesOrNo () {
    let checkConfirmations = [{
        name: 'confirmation',
        type: 'confirm',
        message: 'Do you want to add personal webhook?',
    }];
    return inquirer.prompt(checkConfirmations).then(res => {
        // console.log(res)
        if (res.confirmation) {
            return true
        } else {
           return false
        }
    })
}

function askForCaptcha (){
    let main = [{
        name: 'askForCaptcha',
        type: 'list',
        message: 'Captcha Method',
        choices: ['AYCD', '2Cap', 'Cap Monster', 'No Entry'],
        defeat: function (value) {
       //     console.log(value)
            if (value.length) {
                return true;
            } else {
                return false;
            }
        }
    }];
    return inquirer.prompt(main)
}

async function turnOffAYCD (){
    try {
    const AutoSolve = require('autosolve-client');
    const autoSolve = AutoSolve?.getInstance();
    console.log('Turning Off AYCD')
    await autoSolve.cancelAllRequests()
    } catch (error) {
        
    }
  
  
}

function aycd (){
    let aycd = [{
        name: 'key',
        type: 'input',
        message: 'Enter your AYCD Information ( ACCESSTOKEN:APIKEY )',
        validate: function (value) {
            if (value.length) {
                return true;
            } else {
                return 'Enter your AYCD Information';
            }
        }
    }];
    return inquirer.prompt(aycd)
}

function twoCap () {
    let twoCap = [{
        name: 'key',
        type: 'input',
        message: 'Enter your 2cap key : ',
        validate: function (value) {
            if (value.length) {
                return true;
            } else {
                return 'Enter your 2cap key : ';
            }
        }
    }];
    return inquirer.prompt(twoCap)

}

function capMonster () {
    let capMonster = [{
        name: 'key',
        type: 'input',
        message: 'Enter your Cap Monster Key : ',
        validate: function (value) {
            if (value.length) {
                return true;
            } else {
                return 'Enter your Cap Monster Key';
            }
        }
    }];
    return inquirer.prompt(capMonster)

}

function YesProxy() {
    let checkConfirmations = [{
        name: 'confirmation',
        type: 'confirm',
        message: 'Do you want to add proxies?',
    }];
    return inquirer.prompt(checkConfirmations).then(res => {
        // console.log(res)
        if (res.confirmation) {
            return true
        } else {
            return false
        }
    })
}

function yesTasks() {
    let checkConfirmations = [{
        name: 'confirmation',
        type: 'confirm',
        message: 'Do you want to add Tasks?',
    }];
    return inquirer.prompt(checkConfirmations).then(res => {
        // console.log(res)
        if (res.confirmation) {
            return true
        } else {
            return false
        }
    })
}




module.exports = {
    keyAuth,
    importTasks,
    importProxies,
    checkConfirmation,
    mainMenu,
    setWebhook,
    yesOrNo,
    YesProxy,
    yesTasks,
    askForCaptcha,
    aycd,
    twoCap,
    capMonster,
    turnOffAYCD
};