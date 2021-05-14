const Discord = require('discord.js');
const client = new Discord.Client();
const {
	prefix,
	token,
	bot_age,
	bot_info,
} = require('./config.json');
const {SKUADD , findCommand , deleteSku, checkBank, massAdd, validateUser} = require('./dms');
const fs = require('fs')

client.once('ready', () => {
	console.log(`Logged in as ${client.user.tag}!`);
	console.log(bot_age);
	console.log(prefix);
	console.log(token);
	console.log(bot_info.name);
	console.log(bot_info.version);

	findCommand(client, '!Add', 'Enter SKU like this\n[!]SKUAdd [SKU-NOBRACKETS]');
	SKUADD(client, '!SKUAdd', 'Testing');
	deleteSku(client, '!deleteSku', 'Deleted')
	deleteSku(client, '!skudelete', 'Deleted')
	deleteSku(client, '!skuremove', 'Deleted')
	deleteSku(client, '!removesku', 'Deleted')
	checkBank(client, '!checkBank', 'Returned')
	massAdd(client, '!massAdd', 'Returned')
	validateUser(client, '!validate', 'Returned')
	
	client.users.fetch('202862796965150720').then((user) => {
		user.send('Hello World');
	});
});


client.on('message', async  msg => {
	console.log(msg.content);
msg.repl
});



client.login(token);