const targetMonitor = require('./targetMonitor');
const TargetMonitor = require('./targetMonitor1');
const rp = require('request-promise').defaults({
	followAllRedirects: true,
	resolveWithFullResponse: true,
	gzip : true,
});
function SKUADD(clients, triggerText, replyText) {
	try {
		clients.on('message', async (message) => {
			if(message.channel.type === 'dm' && message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				// message.author.send(replyText);
				const content = message.content;
				const SKU = content.split(' ')[1];
				//    fetch('')
				if(SKU.length == 8) {

					await rp.get({
						url : `https://montiors-default-rtdb.firebaseio.com/sites/target/${SKU}.json`,
					}, ((error, response, body) => {
						console.log(body);
                        console.log(response.statusCode);
                            let parsedBody = JSON.parse(body)
                            console.log(parsedBody)
                            if(parsedBody){
                                let status = parsedBody.Status
                                if(status == 'Active'){
                                    message.author.send(`${SKU} already active in DB`);
                                } else if (status == 'Not Active'){
                                    message.author.send(`${SKU} not active`);
                                    message.author.send(`Adding ${SKU} to monitor`);
                                    TargetMonitor(SKU)
                                } 
                            }
                            
                         else if(body === 'null' || body === null){
                            message.author.send(`New SKU detected - ${SKU}`);
                            message.author.send(`Adding ${SKU} to monitor`);
                            TargetMonitor(SKU)
                        } 
					}));


					
				} else {
					message.author.send('Invalid TCIN');
					message.author.send('Check TCIN Length or Contact Prada#4873');
				}


			}
		});
	} catch (error) {
		console.log(error);
	}

}

function findCommand(clients, triggerText, replyText) {
	clients.on('message', message => {
		if(message.channel.type === 'dm' && message.content.toLowerCase() === triggerText.toLowerCase()) {
			message.author.send(replyText);
		}
	});

}


module.exports = [SKUADD, findCommand];
