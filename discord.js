const Discord = require('discord.js');
const client = new Discord.Client();
const os = require('os')
const {
	prefix,
	token,
	development,
	bot_age,
	bot_info,
} = require('./config.json');
const {SKUADD , findCommand , deleteSku, checkBank, massAdd, validateUser, walmartScraper} = require('./dms');
const fs = require('fs')

client.once('ready', () => {
	console.log(`Logged in as ${client.user.tag}!`);
	console.log(bot_age);
	console.log(prefix);
	console.log(token);
	console.log(bot_info.name);
	console.log(bot_info.version);

	// findCommand(client, '!Add', 'Enter SKU like this\n[!]SKUAdd [SKU-NOBRACKETS]');
	// Add Skus
	// SKUADD(client, '!skuadd', 'Testing');
	// SKUADD(client, '!add', 'Testing');
	// SKUADD(client, '!a', 'Testing');
	// massAdd(client, '!skuadd', 'Testing');
	massAdd(client, '!add', 'Testing');
	// massAdd(client, '!a', 'Testing');
	// Remove
	// deleteSku(client, '!deleteSku', 'Deleted')
	deleteSku(client, '!rm', 'Deleted')
	deleteSku(client, '!remove', 'Deleted')
	deleteSku(client, '!delete', 'Deleted')
	// deleteSku(client, '!skudelete', 'Deleted')
	// deleteSku(client, '!skuremove', 'Deleted')
	// deleteSku(client, '!removesku', 'Deleted')

	// Checking Bank
	// checkBank(client, '!checkBank', 'Returned')
	checkBank(client, '!check', 'Returned')
	checkBank(client, '!bank', 'Returned')
	checkBank(client, '!running', 'Returned')
	checkBank(client, '!list', 'Returned')

	// Mass Add
	// massAdd(client, '!massAdd', 'Returned')
	massAdd(client, '!madd', 'Returned')
	massAdd(client, '!mass', 'Returned')

	validateUser(client, '!validate', 'Returned')
	
	walmartScraper(client, '!walmart', 'Returned')
	
	client.users.fetch('202862796965150720').then((user) => {
		user.send('Hello World');
	});
});


client.on('message', async  msg => {
	console.log(msg.content);
msg.repl
});
client.on("message", (msg) => {
	if (msg.content === prefix + "pricing") {
	  msg.delete()
	  const pricingEmbed = new Discord.MessageEmbed()
		.setColor('#24223a')
		.setTitle('Pricing')
		.setDescription('*Our pricing system is based off how many members you have.*\n\n1-99 Members - **$100/m**\n\n100 - 199 Members - **$200/m**\n\n200 - 249 Members - **$250/m**\n\n250 - 299 Members - **$300/m**\n\n300 - 349 Members - **$350/m**\n\n350 - 399 Members - **$400/m**\n\n400+ Members - **$500/m**\n\n*Prices vary and are subject to change. All payments are made in USD.*')
		.setFooter('Fiber Monitors - fibermonitors.com', 'https://media.discordapp.net/attachments/865335845309644800/865335918609432616/PingMonitors-removebg-preview.png')
		msg.channel.send(pricingEmbed);
	}
  });
  client.on("message", (msg) => {
	if (msg.content.toUpperCase() === prefix + "COMMANDS") {
	  msg.delete()
	  const pricingEmbed = new Discord.MessageEmbed()
		.setColor('#24223a')
		.setTitle('Command List')
		.setDescription(`\n
		**!validate [apikey]**
		Validates User

		**!add [SITE] [SKU] **
		Add product to site

		**!remove [SITE] [SKU]** 
	 	Delete product from site

		**!list** 
		Returns running products

		**!add [site] [SKU] [SKU] [SKU] **
		Ability to add multiple skus at a time 

		**!walmart [SKU]** 
		Scrapes Specific Walmart Product

		*Without brackets and replace SITE and SKU.*
		*If you have any issues feel to ask questions in tickets.*
		`)
		.setFooter('Fiber Monitors - fibermonitors.com', 'https://media.discordapp.net/attachments/865335845309644800/865335918609432616/PingMonitors-removebg-preview.png')
		msg.channel.send(pricingEmbed);
	}
  });
  client.on("message", (msg) => {
	if (msg.content === prefix + "sitelist") {
	  msg.delete()
	  const sitelistEmbed = new Discord.MessageEmbed()
		.setColor('#24223a')
		.setTitle('Sitelist')
		.setDescription('*This sitelist will only grow as we develop. Sitelist varies and is subject to change. We are constantly adding new sites.*')
		.addFields(
		  {
			"name": "Shopify Filtered",
			"value": "🟢 Online",
			"inline": true
		  },
		  {
			"name": "Target",
			"value": "🟢 Online",
			"inline": true
		  },
		  {
			"name": "Walmart",
			"value": "🟢 Online",
			"inline": true
		  },
		  {
			"name": "BestBuy",
			"value": "🟢 Online",
			"inline": true
		  },
		  {
			"name": "Academy",
			"value": "🟢 Online",
			"inline": true
		  },
		  {
			"name": "Newegg",
			"value": "🟢 Online",
			"inline": true
		  },
		  {
			"name": "Gamestop",
			"value": "🟢 Online",
			"inline": true
		  },
		  {
			"name": "Slick Deals",
			"value": "🟢 Online",
			"inline": true
		  },
		  {
			"name": "Shopify Unfiltered",
			"value": "🟠 Under Construction",
			"inline": true
		  }
		)
		.setFooter('Fiber Monitors - fibermonitors.com', 'https://media.discordapp.net/attachments/865335845309644800/865335918609432616/PingMonitors-removebg-preview.png');
		msg.channel.send(sitelistEmbed);
	}
  });

console.log(os.platform())
if(os.platform() == "win32" || os.platform() == "darwin"){
	console.log("Development Environment")
	client.login(development)
} else {
	console.log("Production Environment")
	client.login(token);

}